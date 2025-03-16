package User

import (
	"encoding/binary"
	"fmt"
	"os"
	"proyecto1/DiskManagement"
	"proyecto1/Structs"
	"proyecto1/Utilities"
	"strings"
	"bytes"
)

// Estructura de partición de usuario
type PartitionUser struct {
	IDPartition string
	IDUsuario   string
}

//----------------geters y seters----------------
//Obtener el id de la particion
func (Data *PartitionUser) GetIDPartition() string {
	return Data.IDPartition
}

//Obtener el id del usuario
func (Data *PartitionUser) GetIDUsuario() string {
	return Data.IDUsuario
}

//Establecer el id de la particion
func (Data *PartitionUser) SetIDPartition(IDPartition string) {
	Data.IDPartition = IDPartition
}

//Establecer el id del usuario
func (Data *PartitionUser) SetIDUsuario(IDUsuario string) {
	Data.IDUsuario = IDUsuario
}

var Data PartitionUser

func Login(user string, pass string, id string, buffer *bytes.Buffer) {
	fmt.Println("======Start LOGIN======")
	fmt.Println("User:", user)
	fmt.Println("Pass:", pass)
	fmt.Println("Id:", id)

	// Verificar si el usuario ya está logueado buscando en las particiones montadas
	mountedPartitions := DiskManagement.GetMountedPartitions()
	var filepath string
	var partitionFound bool
	var login bool = false

	for _, partitions := range mountedPartitions {
		for _, Partition := range partitions {
			if Partition.ID == id && Partition.LoggedIn { // Verifica si ya está logueado
				fmt.Fprintf(buffer, "Error LOGIN: Ya existe un usuario logueado en la partición:%s\n", id)
				return
			}
			if Partition.ID == id { // Encuentra la partición correcta
				filepath = Partition.Path
				partitionFound = true
				break
			}
		}
		if partitionFound {
			break
		}
	}

	if !partitionFound {
		fmt.Fprintf(buffer, "Error en LOGIN: No se encontró ninguna partición montada con el ID proporcionado %s\n", id)
		return
	}

	// Abrir archivo binario
	file, err := Utilities.OpenFile(filepath, buffer)
	if err != nil {
		fmt.Println("Error: No se pudo abrir el archivo:", err)
		return
	}
	defer file.Close()

	var TempMBR Structs.MRB
	// Leer el MBR desde el archivo binario
	if err := Utilities.ReadObject(file, &TempMBR, 0, buffer); err != nil {
		fmt.Println("Error: No se pudo leer el MBR:", err)
		return
	}

	// Imprimir el MBR
	Structs.PrintMBR(TempMBR)
	fmt.Println("-------------")

	var index int = -1
	// Iterar sobre las particiones del MBR para encontrar la correcta
	for i := 0; i < 4; i++ {
		if TempMBR.MbrPartitions[i].Size != 0 {
			if strings.Contains(string(TempMBR.MbrPartitions[i].ID[:]), id) {
				fmt.Println("Partition found")
				if TempMBR.MbrPartitions[i].Status[0] == '1' {
					fmt.Println("Partition is mounted")
					index = i
				} else {
					fmt.Println("Partition is not mounted")
					return
				}
				break
			}
		}
	}

	// Iterar sobre las particiones del MBR para encontrar la correcta
	for i := 0; i < 4; i++ {
		if TempMBR.MbrPartitions[i].Size != 0 {
			if strings.Contains(string(TempMBR.MbrPartitions[i].ID[:]), id) {
				if TempMBR.MbrPartitions[i].Status[0] == '1' {
					index = i
				} else {
					return
				}
				break
			}
		}
	}

	if index == -1 {
		fmt.Fprintf(buffer, "Error en LOGIN: no se encontro nunguna particion con el ID %s\n", id)
		return
		}

	var tempSuperblock Structs.Superblock
	// Leer el Superblock desde el archivo binario
	if err := Utilities.ReadObject(file, &tempSuperblock, int64(TempMBR.MbrPartitions[index].Start), buffer); err != nil {
		fmt.Println("Error: No se pudo leer el Superblock:", err)
		return
	}

	// Buscar el archivo de usuarios /users.txt -> retorna índice del Inodo
	indexInode := InitSearch("/users.txt", file, tempSuperblock, buffer)

	var crrInode Structs.Inode
	// Leer el Inodo desde el archivo binario
	if err := Utilities.ReadObject(file, &crrInode, int64(tempSuperblock.SB_Inode_Start+indexInode*int32(binary.Size(Structs.Inode{}))), buffer); err != nil {
		fmt.Println("Error: No se pudo leer el Inodo:", err)
		return
	}

	// Leer datos del archivo
	data := GetInodeFileData(crrInode, file, tempSuperblock, buffer)

	// Dividir la cadena en líneas
	lines := strings.Split(data, "\n")

	// Iterar a través de las líneas para verificar las credenciales
	for _, line := range lines {
		words := strings.Split(line, ",")

		if len(words) == 5 {
			if (strings.Contains(words[3], user)) && (strings.Contains(words[4], pass)) {
				login = true
				break
			}
		}
	}

	// Imprimir información del Inodo
	fmt.Println("Inode", crrInode.IN_Block)

	// Si las credenciales son correctas y marcamos como logueado
	if login {
		fmt.Fprintf(buffer, "Usuario logueado con éxito en la partición:%s\n", id)
		fmt.Println("Usuario logueado con exito")
		DiskManagement.MarkPartitionAsLoggedIn(id) // Marcar la partición como logueada
	}

	// Establecer la partición y el usuario en la estructura de usuario
	Data.SetIDPartition(id)
	Data.SetIDUsuario(user)

	fmt.Println("======End LOGIN======")
}

func InitSearch(path string, file *os.File, tempSuperblock Structs.Superblock, buffer *bytes.Buffer) int32 {
	fmt.Println("======Start BUSQUEDA INICIAL ======")
	fmt.Println("path:", path)
	// path = "/ruta/nueva"

	// split the path by /
	TempStepsPath := strings.Split(path, "/")
	StepsPath := TempStepsPath[1:]

	fmt.Println("StepsPath:", StepsPath, "len(StepsPath):", len(StepsPath))
	for _, step := range StepsPath {
		fmt.Println("step:", step)
	}

	var Inode0 Structs.Inode
	// Read object from bin file
	if err := Utilities.ReadObject(file, &Inode0, int64(tempSuperblock.SB_Inode_Start), buffer); err != nil {
		return -1
	}

	fmt.Println("======End BUSQUEDA INICIAL======")

	return SarchInodeByPath(StepsPath, Inode0, file, tempSuperblock, buffer)
}

// stack
func pop(s *[]string) string {
	lastIndex := len(*s) - 1
	last := (*s)[lastIndex]
	*s = (*s)[:lastIndex]
	return last
}

func SarchInodeByPath(StepsPath []string, Inode Structs.Inode, file *os.File, tempSuperblock Structs.Superblock, buffer *bytes.Buffer) int32 {
	fmt.Println("======Start BUSQUEDA INODO POR PATH======")
	index := int32(0)
	SearchedName := strings.Replace(pop(&StepsPath), " ", "", -1)

	fmt.Println("========== SearchedName:", SearchedName)

	// Iterate over i_blocks from Inode
	for _, block := range Inode.IN_Block {
		if block != -1 {
			if index < 13 {
				//CASO DIRECTO

				var crrFolderBlock Structs.FolderBlock
				// Read object from bin file
				if err := Utilities.ReadObject(file, &crrFolderBlock, int64(tempSuperblock.SB_Block_Start+block*int32(binary.Size(Structs.FolderBlock{}))), buffer); err != nil {
					return -1
				}

				for _, folder := range crrFolderBlock.B_Content {
					// fmt.Println("Folder found======")
					fmt.Println("Folder === Name:", string(folder.B_Name[:]), "B_inodo", folder.B_Inode)

					if strings.Contains(string(folder.B_Name[:]), SearchedName) {

						fmt.Println("len(StepsPath)", len(StepsPath), "StepsPath", StepsPath)
						if len(StepsPath) == 0 {
							fmt.Println("Folder found======")
							return folder.B_Inode
						} else {
							fmt.Println("NextInode======")
							var NextInode Structs.Inode
							// Read object from bin file
							if err := Utilities.ReadObject(file, &NextInode, int64(tempSuperblock.SB_Inode_Start+folder.B_Inode*int32(binary.Size(Structs.Inode{}))), buffer); err != nil {
								return -1
							}
							return SarchInodeByPath(StepsPath, NextInode, file, tempSuperblock, buffer)
						}
					}
				}

			} else {
				fmt.Print("indirectos")
			}
		}
		index++
	}

	fmt.Println("======End BUSQUEDA INODO POR PATH======")
	return 0
}

func GetInodeFileData(Inode Structs.Inode, file *os.File, tempSuperblock Structs.Superblock, buffer *bytes.Buffer) string {
	fmt.Println("======Start CONTENIDO DEL BLOQUE======")
	index := int32(0)
	// define content as a string
	var content string

	// Iterate over i_blocks from Inode
	for _, block := range Inode.IN_Block {
		if block != -1 {
			//Dentro de los directos
			if index < 13 {
				var crrFileBlock Structs.FileBlock
				// Read object from bin file
				if err := Utilities.ReadObject(file, &crrFileBlock, int64(tempSuperblock.SB_Block_Start+block*int32(binary.Size(Structs.FileBlock{}))), buffer); err != nil {
					return ""
				}

				content += string(crrFileBlock.B_Content[:])

			} else {
				fmt.Print("indirectos")
			}
		}
		index++
	}

	fmt.Println("======End CONTENIDO DEL BLOQUE======")
	return content
}

func LogOut(buffer *bytes.Buffer) {
	fmt.Fprint(buffer, "==========LOGOUT==========\n")
	mountedPartitions := DiskManagement.GetMountedPartitions()
	var SesionActiva bool

	if len(mountedPartitions) == 0 {
		fmt.Fprintf(buffer, "Error LOGOUT: No hay ninguna partición montada.\n")
		return
	}

	for _, Particiones := range mountedPartitions {
		for _, Particion := range Particiones {
			if Particion.LoggedIn {
				SesionActiva = true
				break
			}
		}
		if SesionActiva {
			break
		}
	}
	if !SesionActiva {
		fmt.Fprintf(buffer, "Error LOGOUT: No hay ninguna sesión activa.\n")
		return
	} else {
		DiskManagement.MarkPartitionAsLoggedOut(Data.GetIDPartition())
		fmt.Fprintf(buffer, "Sesión cerrada con éxito de la partición:%s\n", Data.GetIDPartition())
	}
	Data.SetIDPartition("")
	Data.SetIDUsuario("")
}

// MKUSER
func AppendToFileBlock(inode *Structs.Inode, newData string, file *os.File, superblock Structs.Superblock, buffer *bytes.Buffer) error {
	// Leer el contenido existente del archivo utilizando la función GetInodeFileData
	existingData := GetInodeFileData(*inode, file, superblock, buffer)

	// Concatenar el nuevo contenido
	fullData := existingData + newData

	// Asegurarse de que el contenido no exceda el tamaño del bloque
	if len(fullData) > len(inode.IN_Block)*binary.Size(Structs.FileBlock{}) {
		// Si el contenido excede, necesitas manejar bloques adicionales
		return fmt.Errorf("el tamaño del archivo excede la capacidad del bloque actual y no se ha implementado la creación de bloques adicionales")
	}

	// Escribir el contenido actualizado en el bloque existente
	var updatedFileBlock Structs.FileBlock
	copy(updatedFileBlock.B_Content[:], fullData)
	if err := Utilities.WriteObject(file, updatedFileBlock, int64(superblock.SB_Block_Start+inode.IN_Block[0]*int32(binary.Size(Structs.FileBlock{}))), buffer); err != nil {
		return fmt.Errorf("error al escribir el bloque actualizado: %v", err)
	}

	// Actualizar el tamaño del inodo
	inode.IN_Size = int32(len(fullData))
	if err := Utilities.WriteObject(file, *inode, int64(superblock.SB_Inode_Start+inode.IN_Block[0]*int32(binary.Size(Structs.Inode{}))), buffer); err != nil {
		return fmt.Errorf("error al actualizar el inodo: %v", err)
	}

	return nil
}

func Mkgrp(name string, buffer *bytes.Buffer) {
	fmt.Fprint(buffer, "=============MKGRP=============\n")

	mountedPartitions := DiskManagement.GetMountedPartitions()
	var filePath string
	var PartitionFound bool

	// Buscar la partición donde se ha iniciado sesión
	for _, Particiones := range mountedPartitions {
		for _, Particion := range Particiones {
			if Particion.ID == Data.GetIDPartition() {
				filePath = Particion.Path
				PartitionFound = true
				break
			}
		}
		if PartitionFound {
			break
		}
	}

	if !PartitionFound {
		fmt.Fprintf(buffer, "Error MKGRP: No se encontró ninguna partición montada con el ID: %s\n", Data.GetIDPartition())
		return
	}

	// Validación del usuario "root"
	if Data.GetIDUsuario() != "root" {
		fmt.Fprintf(buffer, "Error MKGRP: Solo el usuario 'root' puede crear grupos.\n")
		return
	}

	// Abrir el archivo de la partición
	file, err := Utilities.OpenFile(filePath, buffer)
	if err != nil {
		return
	}
	defer file.Close()

	// Leer el MBR de la partición
	var TempMBR Structs.MRB
	if err := Utilities.ReadObject(file, &TempMBR, 0, buffer); err != nil {
		return
	}

	// Buscar la partición correcta en el MBR
	var index int = -1
	for i := 0; i < 4; i++ {
		if TempMBR.MbrPartitions[i].Size != 0 {
			if strings.Contains(string(TempMBR.MbrPartitions[i].ID[:]), Data.GetIDPartition()) {
				if TempMBR.MbrPartitions[i].Status[0] == '1' {
					index = i
				} else {
					fmt.Fprintf(buffer, "Error MKGRP: La partición con ID %s no está montada.\n", Data.GetIDPartition())
					return
				}
				break
			}
		}
	}

	if index == -1 {
		fmt.Fprintf(buffer, "Error MKGRP: No se encontró ninguna partición con el ID: %s\n", Data.GetIDPartition())
		return
	}

	fmt.Fprintf(buffer, "Grupo '%s' creado exitosamente.\n", name)
}