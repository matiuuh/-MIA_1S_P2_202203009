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
				if TempMBR.MbrPartitions[i].Status[0] == '1' {
					fmt.Println("particion montada\n")
					index = i
				} else {
					fmt.Println("particion no montada\n")
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
		fmt.Fprintf(buffer, "Error: No se pudo leer el Inodo:", err, " asegurese de haber ejecutado mkfs correctamente\n")
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
	var content string

	for _, block := range Inode.IN_Block {
		if block != -1 {
			if index < 13 {
				var crrFileBlock Structs.FileBlock
				if err := Utilities.ReadObject(file, &crrFileBlock, int64(tempSuperblock.SB_Block_Start+block*int32(binary.Size(Structs.FileBlock{}))), buffer); err != nil {
					return ""
				}

				// Importante: limpiar caracteres nulos claramente
				cleanData := strings.TrimRight(string(crrFileBlock.B_Content[:]), "\x00")
				content += cleanData

			} else {
				fmt.Print("indirectos")
			}
		}
		index++
	}

	fmt.Println("Contenido final obtenido de users.txt:", content)
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

//--------------------MKGRP--------------------
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
		fmt.Fprintf(buffer, "Error MKGRP: no se pudo abrir disco: %v\n", err)
		return
	}
	defer file.Close()

	// Leer el MBR de la partición
	var TempMBR Structs.MRB
	if err := Utilities.ReadObject(file, &TempMBR, 0, buffer); err != nil {
		fmt.Fprintf(buffer, "Error MKGRP: No se pudo leer el MBR: %v\n", err)
		return
	}

	var partitionIndex int = -1
	for i, part := range TempMBR.MbrPartitions {
		if strings.Trim(string(part.ID[:]), "\x00") == Data.GetIDPartition() && part.Status[0] == '1' {
			partitionIndex = i
			break
		}
	}
	if partitionIndex == -1 {
		fmt.Fprintf(buffer, "Error MKGRP: La partición no está montada o no existe.\n")
		return
	}

	var tempSuperblock Structs.Superblock
	sbStart := int64(TempMBR.MbrPartitions[partitionIndex].Start)
	if err := Utilities.ReadObject(file, &tempSuperblock, sbStart, buffer); err != nil {
		fmt.Fprintf(buffer, "Error MKGRP: No se pudo leer el Superblock: %v\n", err)
		return
	}

	// Buscar archivo /users.txt (obtienes índice de Inodo)
	indexInode := InitSearch("/users.txt", file, tempSuperblock, buffer)
	if indexInode == -1 {
		fmt.Fprintf(buffer, "Error MKGRP: No se encontró el archivo /users.txt\n")
		return
	}

	var usersInode Structs.Inode
	inodePos := int64(tempSuperblock.SB_Inode_Start + indexInode*int32(binary.Size(Structs.Inode{})))
	if err := Utilities.ReadObject(file, &usersInode, inodePos, buffer); err != nil {
		fmt.Fprintf(buffer, "Error MKGRP: No se pudo leer el Inodo: %v\n", err)
		return
	}

	data := GetInodeFileData(usersInode, file, tempSuperblock, buffer)
	lines := strings.Split(data, "\n")

	var groupID int = 1
	for _, line := range lines {
		if strings.TrimSpace(line) == "" {
			continue
		}
		fields := strings.Split(line, ",")
		if len(fields) == 3 && fields[1] == "G" {
			if fields[2] == name {
				fmt.Fprintf(buffer, "Error MKGRP: El grupo '%s' ya existe.\n", name)
				return
			}
			groupID++
		}
	}

	// Añadir el nuevo grupo al archivo claramente
	newGroup := fmt.Sprintf("%d,G,%s\n", groupID, name)

	err = AppendToFileBlock(&usersInode, newGroup, file, tempSuperblock, buffer)
	if err != nil {
		fmt.Fprintf(buffer, "Error MKGRP: Al escribir nuevo grupo: %v\n", err)
		return
	}

	fmt.Fprintf(buffer, "Grupo '%s' creado exitosamente.\n", name)
}

//--------------------MKUSR--------------------
func Mkusr(user string, pass string, grp string, buffer *bytes.Buffer) {
	fmt.Fprint(buffer, "=============MKUSR=============\n")

	// Verificar que el usuario logueado sea 'root'
	if Data.GetIDUsuario() != "root" {
		fmt.Fprintf(buffer, "Error MKUSR: Solo 'root' puede crear usuarios.\n")
		return
	}

	// Validación longitud máxima
	if len(user) > 10 {
		fmt.Fprintf(buffer, "Error MKUSR: Nombre del usuario máximo 10 caracteres.\n")
		return
	}
	if len(pass) > 10 {
		fmt.Fprintf(buffer, "Error MKUSR: Contraseña máximo 10 caracteres.\n")
		return
	}
	if len(grp) > 10 {
		fmt.Fprintf(buffer, "Error MKUSR: Nombre del grupo máximo 10 caracteres.\n")
		return
	}

	// Buscar la partición montada
	mountedPartitions := DiskManagement.GetMountedPartitions()
	var filePath string
	var partitionFound bool
	for _, Particiones := range mountedPartitions {
		for _, Particion := range Particiones {
			if Particion.ID == Data.GetIDPartition() {
				filePath = Particion.Path
				partitionFound = true
				break
			}
		}
		if partitionFound {
			break
		}
	}

	if !partitionFound {
		fmt.Fprintf(buffer, "Error MKUSR: No se encontró partición montada.\n")
		return
	}

	// Abrir el archivo de la partición
	file, err := Utilities.OpenFile(filePath, buffer)
	if err != nil {
		fmt.Fprintf(buffer, "Error MKUSR: No se pudo abrir disco: %v\n", err)
		return
	}
	defer file.Close()

	// Leer el MBR
	var TempMBR Structs.MRB
	if err := Utilities.ReadObject(file, &TempMBR, 0, buffer); err != nil {
		fmt.Fprintf(buffer, "Error MKUSR: No se pudo leer el MBR: %v\n", err)
		return
	}

	// Buscar la partición correcta
	var partitionIndex int = -1
	for i, part := range TempMBR.MbrPartitions {
		if strings.Trim(string(part.ID[:]), "\x00") == Data.GetIDPartition() && part.Status[0] == '1' {
			partitionIndex = i
			break
		}
	}
	if partitionIndex == -1 {
		fmt.Fprintf(buffer, "Error MKUSR: Partición no montada o inexistente.\n")
		return
	}

	// Leer el Superblock
	var tempSuperblock Structs.Superblock
	sbStart := int64(TempMBR.MbrPartitions[partitionIndex].Start)
	if err := Utilities.ReadObject(file, &tempSuperblock, sbStart, buffer); err != nil {
		fmt.Fprintf(buffer, "Error MKUSR: No se pudo leer Superblock: %v\n", err)
		return
	}

	// Buscar archivo /users.txt
	indexInode := InitSearch("/users.txt", file, tempSuperblock, buffer)
	if indexInode == -1 {
		fmt.Fprintf(buffer, "Error MKUSR: No se encontró archivo /users.txt\n")
		return
	}

	var usersInode Structs.Inode
	inodePos := int64(tempSuperblock.SB_Inode_Start + indexInode*int32(binary.Size(Structs.Inode{})))
	if err := Utilities.ReadObject(file, &usersInode, inodePos, buffer); err != nil {
		fmt.Fprintf(buffer, "Error MKUSR: No se pudo leer el Inodo: %v\n", err)
		return
	}

	// Leer los datos del archivo /users.txt
	data := GetInodeFileData(usersInode, file, tempSuperblock, buffer)
	lines := strings.Split(data, "\n")

	var userID int = 1
	var grupoExiste bool = false
	var userFound bool = false
	var userLine string

	// Buscar si el usuario existe en cualquier grupo
	for _, line := range lines {
		if strings.TrimSpace(line) == "" {
			continue
		}
		fields := strings.Split(line, ",")
		if len(fields) == 5 && fields[1] == "U" {
			// Si el usuario ya existe y está eliminado (ID = 0), restaurarlo
			if fields[3] == user && fields[0] == "0" {
				// Restaurar el usuario, cambiar el ID a 1
				fields[0] = "1"
				fields[4] = pass // Actualizar la contraseña
				line = strings.Join(fields, ", ") + "\n"
				userFound = true
				userLine = line
				break
			} else if fields[3] == user {
				// Si el usuario ya existe y no está eliminado, no permitir su creación
				fmt.Fprintf(buffer, "Error MKUSR: El usuario '%s' ya existe.\n", user)
				return
			}
			userID++
		} else if len(fields) == 3 && fields[1] == "G" {
			// Verificar si el grupo existe
			if fields[2] == grp {
				grupoExiste = true
			}
		}
	}

	if userFound {
		// El usuario fue restaurado
		lines = append(lines, userLine) // Añadir la línea restaurada
		// Escribir el archivo actualizado
		newData := strings.Join(lines, "\n")
		err := AppendToFileBlock(&usersInode, newData, file, tempSuperblock, buffer)
		if err != nil {
			fmt.Fprintf(buffer, "Error MKUSR: Al escribir el archivo actualizado: %v\n", err)
			return
		}
		fmt.Fprintf(buffer, "Usuario '%s' restaurado exitosamente.\n", user)
		return
	}

	// Verificar si el grupo existe
	if !grupoExiste {
		// El grupo no existe
		fmt.Fprintf(buffer, "Error MKUSR: El grupo '%s' no existe.\n", grp)
		return
	}

	// Crear un nuevo usuario si no existe
	newUser := fmt.Sprintf("%d,U,%s,%s,%s\n", userID, grp, user, pass)

	// Añadir el nuevo usuario al archivo
	err = AppendToFileBlock(&usersInode, newUser, file, tempSuperblock, buffer)
	if err != nil {
		fmt.Fprintf(buffer, "Error MKUSR: Al escribir nuevo usuario: %v\n", err)
		return
	}

	fmt.Fprintf(buffer, "Usuario '%s' creado exitosamente en el grupo '%s'.\n", user, grp)
}

//--------------------RMUSR--------------------
func Rmusr(user string, buffer *bytes.Buffer) {
	fmt.Fprint(buffer, "=============RMUSR=============\n")

	// Verificar que el usuario logueado sea 'root'
	if Data.GetIDUsuario() != "root" {
		fmt.Fprintf(buffer, "Error RMUSR: Solo 'root' puede eliminar usuarios.\n")
		return
	}

	mountedPartitions := DiskManagement.GetMountedPartitions()
	var filePath string
	var partitionFound bool

	// Buscar la partición donde se ha iniciado sesión
	for _, Particiones := range mountedPartitions {
		for _, Particion := range Particiones {
			if Particion.ID == Data.GetIDPartition() {
				filePath = Particion.Path
				partitionFound = true
				break
			}
		}
		if partitionFound {
			break
		}
	}

	if !partitionFound {
		fmt.Fprintf(buffer, "Error RMUSR: No se encontró ninguna partición montada con el ID: %s\n", Data.GetIDPartition())
		return
	}

	// Abrir el archivo de la partición
	file, err := Utilities.OpenFile(filePath, buffer)
	if err != nil {
		fmt.Fprintf(buffer, "Error RMUSR: No se pudo abrir disco: %v\n", err)
		return
	}
	defer file.Close()

	// Leer el MBR de la partición
	var TempMBR Structs.MRB
	if err := Utilities.ReadObject(file, &TempMBR, 0, buffer); err != nil {
		fmt.Fprintf(buffer, "Error RMUSR: No se pudo leer el MBR: %v\n", err)
		return
	}

	var partitionIndex int = -1
	for i, part := range TempMBR.MbrPartitions {
		if strings.Trim(string(part.ID[:]), "\x00") == Data.GetIDPartition() && part.Status[0] == '1' {
			partitionIndex = i
			break
		}
	}
	if partitionIndex == -1 {
		fmt.Fprintf(buffer, "Error RMUSR: La partición no está montada o no existe.\n")
		return
	}

	var tempSuperblock Structs.Superblock
	sbStart := int64(TempMBR.MbrPartitions[partitionIndex].Start)
	if err := Utilities.ReadObject(file, &tempSuperblock, sbStart, buffer); err != nil {
		fmt.Fprintf(buffer, "Error RMUSR: No se pudo leer Superblock: %v\n", err)
		return
	}

	// Buscar archivo /users.txt
	indexInode := InitSearch("/users.txt", file, tempSuperblock, buffer)
	if indexInode == -1 {
		fmt.Fprintf(buffer, "Error RMUSR: No se encontró el archivo /users.txt\n", err)
		return
	}

	var usersInode Structs.Inode
	inodePos := int64(tempSuperblock.SB_Inode_Start + indexInode*int32(binary.Size(Structs.Inode{})))
	if err := Utilities.ReadObject(file, &usersInode, inodePos, buffer); err != nil {
		fmt.Fprintf(buffer, "Error RMUSR: No se pudo leer el Inodo: %v\n", err)
		return
	}

	// Leer los datos del archivo /users.txt
	data := GetInodeFileData(usersInode, file, tempSuperblock, buffer)
	lines := strings.Split(data, "\n")

	// Buscar y eliminar al usuario
	var newData string
	userFound := false
	for _, line := range lines {
		if strings.TrimSpace(line) == "" {
			continue
		}
		fields := strings.Split(line, ",")
		if len(fields) == 5 && fields[1] == "U" && fields[3] == user {
			// Eliminar la línea correspondiente al usuario
			userFound = true
			fields[0] = "0"  // Cambiar el ID a 0 (usuario eliminado)
			line = strings.Join(fields, ", ") + "\n" // Modificar la línea con ID = 0
		}
		// Conservar el resto de los datos
		newData += line + "\n"
	}

	if !userFound {
		fmt.Fprintf(buffer, "Error RMUSR: El usuario '%s' no existe.\n", user)
		return
	}

	// Escribir el archivo /users.txt actualizado
	err = AppendToFileBlock(&usersInode, newData, file, tempSuperblock, buffer)
	if err != nil {
		fmt.Fprintf(buffer, "Error RMUSR: Al escribir el archivo actualizado: %v\n", err)
		return
	}

	// Confirmación de eliminación
	fmt.Fprintf(buffer, "Usuario '%s' eliminado exitosamente.\n", user)
}