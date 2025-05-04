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
	"math"
	"strconv"
)

// Estructura de partición de usuario
type PartitionUser struct {
	IDPartition string
	IDUsuario   string
	UID         int
	GID         int
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

// Obtener UID
func (Data *PartitionUser) GetUID() int {
	return Data.UID
}

// Establecer UID
func (Data *PartitionUser) SetUID(uid int) {
	Data.UID = uid
}

// Obtener GID
func (Data *PartitionUser) GetGID() int {
	return Data.GID
}

// Establecer GID
func (Data *PartitionUser) SetGID(gid int) {
	Data.GID = gid
}

var Data PartitionUser

func Login(user string, pass string, id string, buffer *bytes.Buffer) {
	fmt.Println("======Start LOGIN======")
	fmt.Println("User:", user)
	fmt.Println("Pass:", pass)
	fmt.Println("Id:", id)

	mountedPartitions := DiskManagement.GetMountedPartitions()
	var filepath string
	var partitionFound bool
	var login bool = false
	var namePart string

	for _, partitions := range mountedPartitions {
		for _, Partition := range partitions {
			if Partition.ID == id && Partition.LoggedIn {
				fmt.Fprintf(buffer, "Error LOGIN: Ya existe un usuario logueado en la partición:%s\n", id)
				return
			}
			if Partition.ID == id {
				filepath = Partition.Path
				namePart = Partition.Name
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

	file, err := Utilities.OpenFile(filepath, buffer)
	if err != nil {
		fmt.Println("Error: No se pudo abrir el archivo:", err)
		return
	}
	defer file.Close()

	var TempMBR Structs.MRB
	if err := Utilities.ReadObject(file, &TempMBR, 0, buffer); err != nil {
		fmt.Println("Error: No se pudo leer el MBR:", err)
		return
	}

	var index int = -1
	NameBytes := [16]byte{}
	copy(NameBytes[:], []byte(namePart))

	for i := 0; i < 4; i++ {
		if TempMBR.MbrPartitions[i].Size != 0 {
			if bytes.Equal(TempMBR.MbrPartitions[i].Name[:], NameBytes[:]) {
				if TempMBR.MbrPartitions[i].Status[0] == '1' {
					index = i
				} else {
					fmt.Fprintf(buffer, "Error LOGIN: La partición %s no está activa en disco.\n", namePart)
					return
				}
				break
			}
		}
	}

	if index == -1 {
		fmt.Fprintf(buffer, "Error en LOGIN: no se encontró ninguna partición con el nombre %s en el disco.\n", namePart)
		return
	}

	var tempSuperblock Structs.Superblock
	if err := Utilities.ReadObject(file, &tempSuperblock, int64(TempMBR.MbrPartitions[index].Start), buffer); err != nil {
		fmt.Println("Error: No se pudo leer el Superblock:", err)
		return
	}

	indexInode := InitSearch("/users.txt", file, tempSuperblock, buffer)
	if indexInode == -1 {
		fmt.Fprintln(buffer, "Error: No se encontró /users.txt, asegúrese de haber ejecutado mkfs.")
		return
	}

	var crrInode Structs.Inode
	if err := Utilities.ReadObject(file, &crrInode, int64(tempSuperblock.SB_Inode_Start+indexInode*int32(binary.Size(Structs.Inode{}))), buffer); err != nil {
		fmt.Fprintf(buffer, "Error: No se pudo leer el Inodo: %v\n", err)
		return
	}

	data := GetInodeFileData(crrInode, file, tempSuperblock, buffer)
	lines := strings.Split(data, "\n")

	var uidEncontrado, gidEncontrado int

	for _, line := range lines {
		words := strings.Split(line, ",")
		if len(words) == 5 {
			if words[3] == user && words[4] == pass {
				login = true
				uidEncontrado, _ = strconv.Atoi(words[0])       // UID del usuario
				gidEncontrado = obtenerIDGrupo(words[2], lines) // GID según nombre grupo
				break
			}
		}
	}

	if login {
		fmt.Fprintf(buffer, "Usuario logueado con éxito en la partición:%s\n", id)
		fmt.Println("Usuario logueado con éxito")
		DiskManagement.MarkPartitionAsLoggedIn(id)
		Data.SetIDPartition(id)
		Data.SetIDUsuario(user)
		Data.SetUID(uidEncontrado)
		Data.SetGID(gidEncontrado)
	} else {
		fmt.Fprintf(buffer, "Error LOGIN: Credenciales incorrectas o usuario no encontrado.\n")
	}

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
func AppendToFileBlock(inode *Structs.Inode, inodeIndex int32, newData string, file *os.File, superblock Structs.Superblock, buffer *bytes.Buffer) error {
	// Leer contenido existente
	existingData := GetInodeFileData(*inode, file, superblock, buffer)
	fullData := existingData + newData
	dataBytes := []byte(fullData)

	blockSize := binary.Size(Structs.FileBlock{})
	numBlocks := int(math.Ceil(float64(len(dataBytes)) / float64(blockSize)))

	//fmt.Fprintf(buffer, "[DEBUG] Datos existentes en users.txt:\n%s\n", existingData)
	//fmt.Fprintf(buffer, "[DEBUG] Tamaño total en bytes: %d, Bloques necesarios: %d\n", len(dataBytes), numBlocks)

	if numBlocks > 12 {
		return fmt.Errorf("el archivo users.txt excede el límite de bloques directos (12)")
	}

	for i := 0; i < numBlocks; i++ {
		start := i * blockSize
		end := start + blockSize
		if end > len(dataBytes) {
			end = len(dataBytes)
		}

		var block Structs.FileBlock
		copy(block.B_Content[:], dataBytes[start:end])

		// Si el bloque no está asignado, asignarlo del bitmap
		if inode.IN_Block[i] == -1 {
			//fmt.Fprintf(buffer, "[DEBUG] Bloque %d no asignado, buscando en bitmap...\n", i)
			var found bool
			for j := int32(0); j < superblock.SB_Blocks_Count; j++ {
				var bit byte
				pos := int64(superblock.SB_Bm_Block_Start + j)
				if err := Utilities.ReadObject(file, &bit, pos, buffer); err != nil {
					return err
				}
				if bit == 0 {
					inode.IN_Block[i] = j
					// Marcar como usado
					if err := Utilities.WriteObject(file, byte(1), pos, buffer); err != nil {
						return err
					}
					//fmt.Fprintf(buffer, "[DEBUG] Bloque libre encontrado y asignado: %d\n", j)
					found = true
					break
				}
			}
			if !found {
				return fmt.Errorf("error: no hay bloques libres disponibles")
			}
		}

		blockOffset := int64(superblock.SB_Block_Start + inode.IN_Block[i]*int32(blockSize))
		//fmt.Fprintf(buffer, "[DEBUG] Escribiendo bloque %d en posición %d\n", i, blockOffset)
		if err := Utilities.WriteObject(file, block, blockOffset, buffer); err != nil {
			return fmt.Errorf("error al escribir bloque %d: %v", i, err)
		}
	}

	inode.IN_Size = int32(len(dataBytes))

	inodeOffset := int64(superblock.SB_Inode_Start + inodeIndex*int32(binary.Size(Structs.Inode{})))
	//fmt.Fprintf(buffer, "[DEBUG] Escribiendo inodo actualizado en posición %d\n", inodeOffset)
	if err := Utilities.WriteObject(file, *inode, inodeOffset, buffer); err != nil {
		return fmt.Errorf("error al actualizar el inodo: %v", err)
	}

	//fmt.Fprintf(buffer, "[DEBUG] Archivo users.txt actualizado con éxito. Total bloques usados: %d\n", numBlocks)
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

	err = AppendToFileBlock(&usersInode, indexInode, newGroup, file, tempSuperblock, buffer)
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
				userFound = true
				return
			}
			userID++
		} else if len(fields) == 3 && fields[1] == "G" {
			// Verificar si el grupo existe
			if fields[2] == grp {
				if fields[0] == "0" {
					fmt.Fprintf(buffer, "Error MKUSR: El grupo '%s' ya fue eliminado y no puede usarse.\n", grp)
					return
				}
				grupoExiste = true
			}
		}
	}

	if userFound {
		// El usuario fue restaurado
		lines = append(lines, userLine) // Añadir la línea restaurada
		// Escribir el archivo actualizado
		newData := strings.Join(lines, "\n")
		err := AppendToFileBlock(&usersInode, indexInode, newData, file, tempSuperblock, buffer)
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
	err = AppendToFileBlock(&usersInode, indexInode, newUser, file, tempSuperblock, buffer)
	if err != nil {
		fmt.Fprintf(buffer, "Error MKUSR: Al escribir nuevo usuario: %v\n", err)
		return
	}

	fmt.Fprintf(buffer, "Usuario '%s' creado exitosamente en el grupo '%s'.\n", user, grp)
}

//--------------------RMUSR--------------------
func Rmusr(user string, buffer *bytes.Buffer) {
	fmt.Fprint(buffer, "=============RMUSR=============\n")

	if Data.GetIDUsuario() != "root" {
		fmt.Fprintf(buffer, "Error RMUSR: Solo 'root' puede eliminar usuarios.\n")
		return
	}

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
		fmt.Fprintf(buffer, "Error RMUSR: No se encontró ninguna partición montada con el ID: %s\n", Data.GetIDPartition())
		return
	}

	file, err := Utilities.OpenFile(filePath, buffer)
	if err != nil {
		fmt.Fprintf(buffer, "Error RMUSR: No se pudo abrir disco: %v\n", err)
		return
	}
	defer file.Close()

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

	indexInode := InitSearch("/users.txt", file, tempSuperblock, buffer)
	if indexInode == -1 {
		fmt.Fprintf(buffer, "Error RMUSR: No se encontró el archivo /users.txt\n")
		return
	}

	var usersInode Structs.Inode
	inodePos := int64(tempSuperblock.SB_Inode_Start + indexInode*int32(binary.Size(Structs.Inode{})))
	if err := Utilities.ReadObject(file, &usersInode, inodePos, buffer); err != nil {
		fmt.Fprintf(buffer, "Error RMUSR: No se pudo leer el Inodo: %v\n", err)
		return
	}

	data := GetInodeFileData(usersInode, file, tempSuperblock, buffer)
	lines := strings.Split(data, "\n")

	var updatedLines []string
	userFound := false

	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}

		fields := strings.Split(line, ",")
		if len(fields) == 5 && fields[1] == "U" && fields[3] == user && fields[0] != "0" {
			fields[0] = "0" // Marcar como eliminado
			userFound = true
			updatedLine := strings.Join(fields, ",")
			updatedLines = append(updatedLines, updatedLine)
		} else {
			updatedLines = append(updatedLines, line)
		}
	}

	if !userFound {
		fmt.Fprintf(buffer, "Error RMUSR: El usuario '%s' no existe o ya fue eliminado.\n", user)
		return
	}

	newData := strings.Join(updatedLines, "\n") + "\n"

	err = OverwriteFileBlock(&usersInode, indexInode, newData, file, tempSuperblock, buffer)
	if err != nil {
		fmt.Fprintf(buffer, "Error RMUSR: Al escribir el archivo actualizado: %v\n", err)
		return
	}

	fmt.Fprintf(buffer, "Usuario '%s' eliminado exitosamente.\n", user)
}

func OverwriteFileBlock(inode *Structs.Inode, inodeIndex int32, newData string, file *os.File, superblock Structs.Superblock, buffer *bytes.Buffer) error {
	// Limpiar bloques anteriores
	//blockSize := binary.Size(Structs.FileBlock{})
	for i := 0; i < 12; i++ {
		if inode.IN_Block[i] != -1 {
			// Limpiar el bitmap
			pos := int64(superblock.SB_Bm_Block_Start + inode.IN_Block[i])
			if err := Utilities.WriteObject(file, byte(0), pos, buffer); err != nil {
				return err
			}
			inode.IN_Block[i] = -1
		}
	}

	// Usar AppendToFileBlock para volver a escribir desde cero
	return AppendToFileBlock(inode, inodeIndex, newData, file, superblock, buffer)
}

//--------------------RMGRP--------------------
func Rmgrp(name string, buffer *bytes.Buffer) {
	fmt.Fprint(buffer, "=============RMGRP=============\n")

	// Validar que el usuario sea root
	if Data.GetIDUsuario() != "root" {
		fmt.Fprintf(buffer, "Error RMGRP: Solo el usuario 'root' puede eliminar grupos.\n")
		return
	}

	// Buscar partición montada
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
		fmt.Fprintf(buffer, "Error RMGRP: No se encontró ninguna partición montada con el ID: %s\n", Data.GetIDPartition())
		return
	}

	// Abrir archivo del disco
	file, err := Utilities.OpenFile(filePath, buffer)
	if err != nil {
		fmt.Fprintf(buffer, "Error RMGRP: No se pudo abrir el disco: %v\n", err)
		return
	}
	defer file.Close()

	// Leer MBR y Superblock
	var TempMBR Structs.MRB
	if err := Utilities.ReadObject(file, &TempMBR, 0, buffer); err != nil {
		fmt.Fprintf(buffer, "Error RMGRP: No se pudo leer el MBR: %v\n", err)
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
		fmt.Fprintf(buffer, "Error RMGRP: La partición no está montada o no existe.\n")
		return
	}

	var tempSuperblock Structs.Superblock
	sbStart := int64(TempMBR.MbrPartitions[partitionIndex].Start)
	if err := Utilities.ReadObject(file, &tempSuperblock, sbStart, buffer); err != nil {
		fmt.Fprintf(buffer, "Error RMGRP: No se pudo leer el Superblock: %v\n", err)
		return
	}

	// Buscar /users.txt
	indexInode := InitSearch("/users.txt", file, tempSuperblock, buffer)
	if indexInode == -1 {
		fmt.Fprintf(buffer, "Error RMGRP: No se encontró el archivo /users.txt\n")
		return
	}

	var usersInode Structs.Inode
	inodePos := int64(tempSuperblock.SB_Inode_Start + indexInode*int32(binary.Size(Structs.Inode{})))
	if err := Utilities.ReadObject(file, &usersInode, inodePos, buffer); err != nil {
		fmt.Fprintf(buffer, "Error RMGRP: No se pudo leer el Inodo: %v\n", err)
		return
	}

	// Leer contenido actual
	data := GetInodeFileData(usersInode, file, tempSuperblock, buffer)
	lines := strings.Split(data, "\n")

	var updatedLines []string
	groupFound := false

	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		fields := strings.Split(line, ",")
		if len(fields) == 3 && fields[1] == "G" && fields[2] == name && fields[0] != "0" {
			fields[0] = "0" // marcar como eliminado
			groupFound = true
			updatedLines = append(updatedLines, strings.Join(fields, ","))
		} else {
			updatedLines = append(updatedLines, line)
		}
	}

	if !groupFound {
		fmt.Fprintf(buffer, "Error RMGRP: El grupo '%s' no existe o ya fue eliminado.\n", name)
		return
	}

	// Guardar el contenido nuevo
	newData := strings.Join(updatedLines, "\n") + "\n"
	err = OverwriteFileBlock(&usersInode, indexInode, newData, file, tempSuperblock, buffer)
	if err != nil {
		fmt.Fprintf(buffer, "Error RMGRP: Al escribir el archivo actualizado: %v\n", err)
		return
	}

	fmt.Fprintf(buffer, "Grupo '%s' eliminado exitosamente.\n", name)
}

//--------------------CHGRP--------------------
func Chgrp(user string, newGroup string, buffer *bytes.Buffer) {
	fmt.Fprint(buffer, "=============CHGRP=============\n")

	// Validar que solo 'root' puede ejecutar el comando
	if Data.GetIDUsuario() != "root" {
		fmt.Fprintf(buffer, "Error CHGRP: Solo el usuario 'root' puede cambiar grupos.\n")
		return
	}

	// Buscar la partición activa
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
		fmt.Fprintf(buffer, "Error CHGRP: No se encontró partición montada con ID: %s\n", Data.GetIDPartition())
		return
	}

	// Abrir el archivo
	file, err := Utilities.OpenFile(filePath, buffer)
	if err != nil {
		fmt.Fprintf(buffer, "Error CHGRP: No se pudo abrir el archivo: %v\n", err)
		return
	}
	defer file.Close()

	// Leer MBR y superblock
	var TempMBR Structs.MRB
	if err := Utilities.ReadObject(file, &TempMBR, 0, buffer); err != nil {
		fmt.Fprintf(buffer, "Error CHGRP: No se pudo leer el MBR: %v\n", err)
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
		fmt.Fprintf(buffer, "Error CHGRP: La partición no está montada o no existe.\n")
		return
	}

	var tempSuperblock Structs.Superblock
	sbStart := int64(TempMBR.MbrPartitions[partitionIndex].Start)
	if err := Utilities.ReadObject(file, &tempSuperblock, sbStart, buffer); err != nil {
		fmt.Fprintf(buffer, "Error CHGRP: No se pudo leer el Superblock: %v\n", err)
		return
	}

	// Buscar users.txt
	indexInode := InitSearch("/users.txt", file, tempSuperblock, buffer)
	if indexInode == -1 {
		fmt.Fprintf(buffer, "Error CHGRP: No se encontró el archivo /users.txt\n")
		return
	}

	var usersInode Structs.Inode
	inodePos := int64(tempSuperblock.SB_Inode_Start + indexInode*int32(binary.Size(Structs.Inode{})))
	if err := Utilities.ReadObject(file, &usersInode, inodePos, buffer); err != nil {
		fmt.Fprintf(buffer, "Error CHGRP: No se pudo leer el inodo de users.txt: %v\n", err)
		return
	}

	data := GetInodeFileData(usersInode, file, tempSuperblock, buffer)
	lines := strings.Split(data, "\n")

	var updatedLines []string
	var userFound, groupExists bool

	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}

		fields := strings.Split(line, ",")

		// Verificar si el grupo existe (ID != 0)
		if len(fields) == 3 && fields[1] == "G" && fields[2] == newGroup && fields[0] != "0" {
			groupExists = true
		}

		// Buscar al usuario y cambiar su grupo
		if len(fields) == 5 && fields[1] == "U" && fields[3] == user && fields[0] != "0" {
			userFound = true
			fields[2] = newGroup // Cambiar grupo
			line = strings.Join(fields, ",")
		}

		updatedLines = append(updatedLines, line)
	}

	if !groupExists {
		fmt.Fprintf(buffer, "Error CHGRP: El grupo '%s' no existe o está eliminado.\n", newGroup)
		return
	}

	if !userFound {
		fmt.Fprintf(buffer, "Error CHGRP: El usuario '%s' no existe o está eliminado.\n", user)
		return
	}

	newData := strings.Join(updatedLines, "\n") + "\n"
	err = OverwriteFileBlock(&usersInode, indexInode, newData, file, tempSuperblock, buffer)
	if err != nil {
		fmt.Fprintf(buffer, "Error CHGRP: No se pudo actualizar el archivo: %v\n", err)
		return
	}

	fmt.Fprintf(buffer, "Grupo del usuario '%s' actualizado exitosamente a '%s'.\n", user, newGroup)
}

func obtenerIDGrupo(nombreGrupo string, lines []string) int {
	for _, line := range lines {
		parts := strings.Split(line, ",")
		if len(parts) == 3 && parts[1] == "G" && parts[2] == nombreGrupo {
			id, _ := strconv.Atoi(parts[0])
			return id
		}
	}
	return -1
}
