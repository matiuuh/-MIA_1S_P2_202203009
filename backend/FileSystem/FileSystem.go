
package FileSystem

import (
	"encoding/binary"
	"fmt"
	"bytes"
	"time"
	"os"
	"proyecto1/DiskManagement"
	"proyecto1/Structs"
	"proyecto1/Utilities"
	"strings"
	"proyecto1/User"
	"strconv"
)

func Mkfs(id string, type_ string, fs_ string, buffer *bytes.Buffer) {
	fmt.Fprintf(buffer, "========MKFS========\n")

	var MountedPartitions DiskManagement.MountedPartition
	var ParticionEncontrada bool

	for _, Particiones := range DiskManagement.GetMountedPartitions() {
		for _, Particion := range Particiones {
			if Particion.ID == id {
				MountedPartitions = Particion
				ParticionEncontrada = true
				break
			}
		}
		if ParticionEncontrada {
			break
		}
	}

	if !ParticionEncontrada {
		fmt.Fprintf(buffer, "Error MFKS: La partición: %s no existe.\n", id)
		return
	}

	if MountedPartitions.Status != '1' {
		fmt.Fprintf(buffer, "Error MFKS: La partición %s aún no está montada.\n", id)
		return
	}

	archivo, err := Utilities.OpenFile(MountedPartitions.Path, buffer)
	if err != nil {
		return
	}

	var TempMBR Structs.MRB
	if err := Utilities.ReadObject(archivo, &TempMBR, 0, buffer); err != nil {
		return
	}

	var IndiceParticion int = -1
	for i := 0; i < 4; i++ {
		if TempMBR.MbrPartitions[i].Size != 0 {
			if strings.Contains(string(TempMBR.MbrPartitions[i].ID[:]), id) {
				IndiceParticion = i
				break
			}
		}
	}

	if IndiceParticion == -1 {
		fmt.Fprintf(buffer, "Error MFKS: La partición: %s no existe.\n", id)
		return
	}

	numerador := int32(TempMBR.MbrPartitions[IndiceParticion].Size - int32(binary.Size(Structs.Superblock{})))
	denrominador_base := int32(4 + int32(binary.Size(Structs.Inode{})) + 3*int32(binary.Size(Structs.FileBlock{})))
	denrominador := denrominador_base
	n := int32(numerador / denrominador)

	// Crear el Superbloque
	var NuevoSuperBloque Structs.Superblock
	NuevoSuperBloque.SB_FileSystem_Type = 2
	NuevoSuperBloque.SB_Inodes_Count = n
	NuevoSuperBloque.SB_Blocks_Count = 3 * n
	NuevoSuperBloque.SB_Free_Blocks_Count = 3*n - 2
	NuevoSuperBloque.SB_Free_Inodes_Count = n - 2
	FechaActual := time.Now()
	FechaString := FechaActual.Format("2006-01-02 15:04:05")
	FechaBytes := []byte(FechaString)
	copy(NuevoSuperBloque.SB_Mtime[:], FechaBytes)
	copy(NuevoSuperBloque.SB_Umtime[:], FechaBytes)
	NuevoSuperBloque.SB_Mnt_Count = 1
	NuevoSuperBloque.SB_Magic = 0xEF53
	NuevoSuperBloque.SB_Inode_Size = int32(binary.Size(Structs.Inode{}))
	NuevoSuperBloque.SB_Block_Size = int32(binary.Size(Structs.FileBlock{}))
	// Calcular las posiciones de los bloques
	NuevoSuperBloque.SB_Bm_Inode_Start = TempMBR.MbrPartitions[IndiceParticion].Start + int32(binary.Size(Structs.Superblock{}))
	NuevoSuperBloque.SB_Bm_Block_Start = NuevoSuperBloque.SB_Bm_Inode_Start + n
	NuevoSuperBloque.SB_Inode_Start = NuevoSuperBloque.SB_Bm_Block_Start + 3*n
	NuevoSuperBloque.SB_Block_Start = NuevoSuperBloque.SB_Inode_Start + n*int32(binary.Size(Structs.Inode{}))
	// Escribir el superbloque en el archivo
	SistemaEXT2(n, TempMBR.MbrPartitions[IndiceParticion], NuevoSuperBloque, FechaString, archivo, buffer)
	defer archivo.Close()
}

func SistemaEXT2(n int32, Particion Structs.Partition, NuevoSuperBloque Structs.Superblock, Fecha string, archivo *os.File, buffer *bytes.Buffer) {
	for i := int32(0); i < n; i++ {
		err := Utilities.WriteObject(archivo, byte(0), int64(NuevoSuperBloque.SB_Bm_Inode_Start+i), buffer)
		if err != nil {
			return
		}
	}
	for i := int32(0); i < 3*n; i++ {
		err := Utilities.WriteObject(archivo, byte(0), int64(NuevoSuperBloque.SB_Bm_Block_Start+i), buffer)
		if err != nil {
			return
		}
	}
	// Inicializa inodos y bloques con valores predeterminados
	if err := initInodesAndBlocks(n, NuevoSuperBloque, archivo, buffer); err != nil {
		fmt.Println("Error: ", err)
		return
	}
	// Crea la carpeta raíz y el archivo users.txt
	if err := createRootAndUsersFile(NuevoSuperBloque, Fecha, archivo, buffer); err != nil {
		fmt.Println("Error: ", err)
		return
	}
	// Escribe el superbloque actualizado al archivo
	if err := Utilities.WriteObject(archivo, NuevoSuperBloque, int64(Particion.Start), buffer); err != nil {
		fmt.Println("Error: ", err)
		return
	}
	// Marca los primeros inodos y bloques como usados
	if err := markUsedInodesAndBlocks(NuevoSuperBloque, archivo, buffer); err != nil {
		fmt.Println("Error: ", err)
		return
	}
	// Imprimir el Superblock final
	Structs.PrintSuperblock(NuevoSuperBloque)
	fmt.Fprintf(buffer, "Partición: %s formateada exitosamente.\n", string(Particion.Name[:]))

	// Depuración: Leer e imprimir el Inodo 0 y el FolderBlock 0
	var inode0 Structs.Inode
	if err := Utilities.ReadObject(archivo, &inode0, int64(NuevoSuperBloque.SB_Inode_Start), buffer); err == nil {
		fmt.Fprintln(buffer, "\n--- Inodo Raíz (0) ---")
		Structs.PrintInode(inode0)
	} else {
		fmt.Fprintln(buffer, "Error leyendo Inodo raíz para depuración:", err)
	}

	var folder0 Structs.FolderBlock
	if err := Utilities.ReadObject(archivo, &folder0, int64(NuevoSuperBloque.SB_Block_Start), buffer); err == nil {
		fmt.Fprintln(buffer, "\n--- Folder Block Raíz (bloque 0) ---")
		Structs.PrintFolderblock(folder0)
	} else {
		fmt.Fprintln(buffer, "Error leyendo FolderBlock raíz para depuración:", err)
	}

	// Depuración: Leer e imprimir el Inodo 1 y FileBlock 1
	var inode1 Structs.Inode
	if err := Utilities.ReadObject(archivo, &inode1, int64(NuevoSuperBloque.SB_Inode_Start+int32(binary.Size(Structs.Inode{}))), buffer); err == nil {
		fmt.Fprintln(buffer, "\n--- Inodo users.txt (1) ---")
		Structs.PrintInode(inode1)
	} else {
		fmt.Fprintln(buffer, "Error leyendo Inodo users.txt para depuración:", err)
	}

	var fileblock1 Structs.FileBlock
	if err := Utilities.ReadObject(archivo, &fileblock1, int64(NuevoSuperBloque.SB_Block_Start+int32(binary.Size(Structs.FolderBlock{}))), buffer); err == nil {
		fmt.Fprintln(buffer, "\n--- FileBlock contenido de users.txt ---")
		fmt.Fprintf(buffer, "Contenido: %s\n", string(fileblock1.B_Content[:]))
	} else {
		fmt.Fprintln(buffer, "Error leyendo FileBlock users.txt para depuración:", err)
	}

}

// Función auxiliar para inicializar inodos y bloques
func initInodesAndBlocks(n int32, newSuperblock Structs.Superblock, file *os.File, buffer *bytes.Buffer) error {
	var newInode Structs.Inode
	for i := int32(0); i < 15; i++ {
		newInode.IN_Block[i] = -1
	}

	for i := int32(0); i < n; i++ {
		if err := Utilities.WriteObject(file, newInode, int64(newSuperblock.SB_Inode_Start+i*int32(binary.Size(Structs.Inode{}))), buffer); err != nil {
			return err
		}
	}

	var newFileblock Structs.FileBlock
	for i := int32(0); i < 3*n; i++ {
		if err := Utilities.WriteObject(file, newFileblock, int64(newSuperblock.SB_Block_Start+i*int32(binary.Size(Structs.FileBlock{}))), buffer); err != nil {
			return err
		}
	}

	return nil
}

// Función auxiliar para crear la carpeta raíz y el archivo users.txt
func createRootAndUsersFile(newSuperblock Structs.Superblock, date string, file *os.File, buffer *bytes.Buffer) error {
	var Inode0, Inode1 Structs.Inode
	initInode(&Inode0, date)
	initInode(&Inode1, date)

	Inode0.IN_Block[0] = 0
	Inode1.IN_Block[0] = 1

	// Asignar el tamaño real del contenido
	data := "1,G,root\n1,U,root,root,123\n"
	actualSize := int32(len(data))
	Inode1.IN_Size = actualSize // Esto ahora refleja el tamaño real del contenido

	var Fileblock1 Structs.FileBlock
	copy(Fileblock1.B_Content[:], data) // Copia segura de datos a FileBlock

	var Folderblock0 Structs.FolderBlock
	Folderblock0.B_Content[0].B_Inode = 0
	copy(Folderblock0.B_Content[0].B_Name[:], ".")
	Folderblock0.B_Content[1].B_Inode = 0
	copy(Folderblock0.B_Content[1].B_Name[:], "..")
	Folderblock0.B_Content[2].B_Inode = 1
	copy(Folderblock0.B_Content[2].B_Name[:], "users.txt")

	// Escribir los inodos y bloques en las posiciones correctas
	if err := Utilities.WriteObject(file, Inode0, int64(newSuperblock.SB_Inode_Start), buffer); err != nil {
		return err
	}
	if err := Utilities.WriteObject(file, Inode1, int64(newSuperblock.SB_Inode_Start+int32(binary.Size(Structs.Inode{}))), buffer); err != nil {
		return err
	}
	if err := Utilities.WriteObject(file, Folderblock0, int64(newSuperblock.SB_Block_Start), buffer); err != nil {
		return err
	}
	if err := Utilities.WriteObject(file, Fileblock1, int64(newSuperblock.SB_Block_Start+int32(binary.Size(Structs.FolderBlock{}))), buffer); err != nil {
		return err
	}

	return nil
}

// Función auxiliar para inicializar un inodo
func initInode(inode *Structs.Inode, date string) {
	inode.IN_Uid = 1
	inode.IN_Gid = 1
	inode.IN_Size = 0
	copy(inode.IN_Atime[:], date)
	copy(inode.IN_Ctime[:], date)
	copy(inode.IN_Mtime[:], date)
	copy(inode.IN_Perm[:], "664")

	for i := int32(0); i < 15; i++ {
		inode.IN_Block[i] = -1
	}
}

// Función auxiliar para marcar los inodos y bloques usados
func markUsedInodesAndBlocks(newSuperblock Structs.Superblock, file *os.File, buffer *bytes.Buffer) error {
	if err := Utilities.WriteObject(file, byte(1), int64(newSuperblock.SB_Bm_Inode_Start), buffer); err != nil {
		return err
	}
	if err := Utilities.WriteObject(file, byte(1), int64(newSuperblock.SB_Bm_Inode_Start+1), buffer); err != nil {
		return err
	}
	if err := Utilities.WriteObject(file, byte(1), int64(newSuperblock.SB_Bm_Block_Start), buffer); err != nil {
		return err
	}
	if err := Utilities.WriteObject(file, byte(1), int64(newSuperblock.SB_Bm_Block_Start+1), buffer); err != nil {
		return err
	}
	return nil
}

func CAT(files []string, buffer *bytes.Buffer) {
	fmt.Println("Usuario activo:", User.Data.GetIDUsuario())
	fmt.Println("Partición activa:", User.Data.GetIDPartition())
	//fmt.Fprintf(buffer, "Error: No hay un usuario logueado")
	// Check if a user is logged in
	if !isUserLoggedIn() {
		fmt.Fprintf(buffer, "Error: No hay un usuario logueado")
		return
	}

	// Check if the user has permission
	if !tienePermiso(buffer) {
		fmt.Fprintf(buffer, "Error: El usuario no tiene permiso de lectura (permiso 777 requerido)")
		return
	}

	// Get the mounted partition information
	ParticionesMount := DiskManagement.GetMountedPartitions()
	var filepath string
	var id string

	// Find the logged-in partition
	for _, partitions := range ParticionesMount {
		for _, partition := range partitions {
			if partition.LoggedIn {
				filepath = partition.Path
				id = partition.ID
				break
			}
		}
		if filepath != "" {
			break
		}
	}

	// Open the file
	file, err := Utilities.OpenFile(filepath, buffer)
	if err != nil {
		fmt.Println("Error: No se pudo abrir el archivo:", err)
		return
	}
	defer file.Close()

	// Read the MBR
	var TempMBR Structs.MRB
	if err := Utilities.ReadObject(file, &TempMBR, 0, buffer); err != nil {
		fmt.Println("Error: No se pudo leer el MBR:", err)
		return
	}

	// Find the correct partition
	var index int = -1
	for i := 0; i < 4; i++ {
		if TempMBR.MbrPartitions[i].Size != 0 && strings.Contains(string(TempMBR.MbrPartitions[i].ID[:]), id) {
			if TempMBR.MbrPartitions[i].Status[0] == '1' {
				index = i
				break
			}
		}
	}

	if index == -1 {
		fmt.Println("Error: No se encontró la partición")
		return
	}

	// Read the Superblock
	var tempSuperblock Structs.Superblock
	if err := Utilities.ReadObject(file, &tempSuperblock, int64(TempMBR.MbrPartitions[index].Start), buffer); err != nil {
		fmt.Println("Error: No se pudo leer el Superblock:", err)
		return
	}

	// Process each file in the input
	for _, filePath := range files {
		fmt.Printf("Imprimiendo el contenido de %s\n", filePath)

		indexInode := buscarInodoPorRuta(filePath, file, tempSuperblock, buffer)
		if indexInode == -1 {
			fmt.Printf("Error: No se pudo encontrar el archivo %s\n", filePath)
			//Structs.PrintFileblock(Structs.FileBlock{}, buffer)
			//fmt.Fprintf(buffer, "------esto despues de que no se pudo encontrar------")
			continue
		}

		var crrInode Structs.Inode
		if err := Utilities.ReadObject(file, &crrInode, int64(tempSuperblock.SB_Inode_Start+indexInode*int32(binary.Size(Structs.Inode{}))), buffer); err != nil {
			fmt.Printf("Error: No se pudo leer el Inode para %s\n", filePath)
			continue
		}

		// Read and print the content of each block in the file
		for _, block := range crrInode.IN_Block {
			if block != -1 {
				var fileblock Structs.FileBlock
				if err := Utilities.ReadObject(file, &fileblock, int64(tempSuperblock.SB_Block_Start+block*int32(binary.Size(Structs.FileBlock{}))), buffer); err != nil {
					fmt.Printf("Error: No se pudo leer el FileBlock para %s\n", filePath)
					continue
				}
				Structs.PrintFileblock(fileblock, buffer)
			}
		}

		fmt.Println("------FIN CAT------")
	}
}

// Función para verificar si un usuario está logueado
func isUserLoggedIn() bool {
	fmt.Println("DEBUG: Usuario actual:", User.Data.GetIDUsuario())
	fmt.Println("DEBUG: Partición actual:", User.Data.GetIDPartition())
	return User.Data.GetIDUsuario() != "" && User.Data.GetIDPartition() != ""
}

// Función para verificar si el usuario tiene permisos
func tienePermiso(buffer *bytes.Buffer) bool {
	// Verificamos si hay sesión activa
	if !isUserLoggedIn() {
		fmt.Fprintln(buffer, "DEBUG: No hay sesión activa.")
		return false
	}

	// Si el usuario es root, siempre tiene permisos
	if User.Data.GetIDUsuario() == "root" {
		fmt.Fprintln(buffer, "DEBUG: Usuario root tiene todos los permisos.")
		return true
	}
	
	// Obtener la partición montada activa
	ParticionesMount := DiskManagement.GetMountedPartitions()
	var filepath, id string
	for _, partitions := range ParticionesMount {
		for _, partition := range partitions {
			if partition.ID == User.Data.GetIDPartition() && partition.LoggedIn {
				filepath = partition.Path
				id = partition.ID
				break
			}
		}
	}

	// Si no encontró la partición activa
	if filepath == "" {
		fmt.Fprintln(buffer, "DEBUG: No se encontró la partición activa.")
		return false
	}

	// Abrimos el archivo del disco
	file, err := Utilities.OpenFile(filepath, buffer)
	if err != nil {
		fmt.Fprintln(buffer, "DEBUG: No se pudo abrir el archivo:", err)
		return false
	}
	defer file.Close()

	// Leemos el MBR
	var TempMBR Structs.MRB
	if err := Utilities.ReadObject(file, &TempMBR, 0, buffer); err != nil {
		fmt.Fprintln(buffer, "DEBUG: No se pudo leer el MBR:", err)
		return false
	}

	// Buscamos el índice de la partición
	var index int = -1
	for i := 0; i < 4; i++ {
		if TempMBR.MbrPartitions[i].Size != 0 &&
			strings.Contains(string(TempMBR.MbrPartitions[i].ID[:]), id) &&
			TempMBR.MbrPartitions[i].Status[0] == '1' {
			index = i
			break
		}
	}
	if index == -1 {
		fmt.Fprintln(buffer, "DEBUG: No se encontró la partición válida.")
		return false
	}

	// Leemos el Superblock
	var sb Structs.Superblock
	if err := Utilities.ReadObject(file, &sb, int64(TempMBR.MbrPartitions[index].Start), buffer); err != nil {
		fmt.Fprintln(buffer, "DEBUG: No se pudo leer el Superblock.")
		return false
	}

	// Buscamos el inodo de /users.txt
	indexInode := buscarStart("/users.txt", file, sb, buffer)
	if indexInode == -1 {
		fmt.Fprintln(buffer, "DEBUG: No se encontró el inodo de /users.txt")
		return false
	}

	// Leemos el Inodo
	var inode Structs.Inode
	if err := Utilities.ReadObject(file, &inode, int64(sb.SB_Inode_Start+indexInode*int32(binary.Size(Structs.Inode{}))), buffer); err != nil {
		fmt.Fprintln(buffer, "DEBUG: No se pudo leer el inodo.")
		return false
	}

	perm := strings.Trim(string(inode.IN_Perm[:]), "\x00")
	fmt.Fprintf(buffer, "DEBUG: Permiso leído: %s\n", perm)

	// Verificamos si tiene permisos de lectura (simplificado)
	return strings.HasPrefix(perm, "6")
}

// Función modificada para buscar y leer Fileblocks en lugar de Folderblocks
func buscarStart(path string, file *os.File, tempSuperblock Structs.Superblock, buffer *bytes.Buffer) int32 {
	TempStepsPath := strings.Split(path, "/")
	RutaPasada := TempStepsPath[1:]

	var Inode0 Structs.Inode
	if err := Utilities.ReadObject(file, &Inode0, int64(tempSuperblock.SB_Inode_Start), buffer); err != nil {
		return -1
	}

	return BuscarInodoRuta(RutaPasada, Inode0, file, tempSuperblock, buffer)
}

// Cambiado para manejar FileBlock en lugar de Folderblock
func BuscarInodoRuta(RutaPasada []string, Inode Structs.Inode, file *os.File, tempSuperblock Structs.Superblock, buffer *bytes.Buffer) int32 {
	if len(RutaPasada) == 0 {
		return -1
	}

	SearchedName := strings.TrimSpace(RutaPasada[0])
	RutaPasada = RutaPasada[1:]

	for _, block := range Inode.IN_Block {
		if block == -1 {
			continue
		}

		blockOffset := int64(tempSuperblock.SB_Block_Start + block*int32(binary.Size(Structs.FolderBlock{})))
		var folder Structs.FolderBlock
		if err := Utilities.ReadObject(file, &folder, blockOffset, buffer); err != nil {
			continue
		}

		for _, entry := range folder.B_Content {
			name := strings.Trim(string(entry.B_Name[:]), "\x00")
			fmt.Fprintf(buffer, "DEBUG: comparando con entrada '%s'\n", name)
			if name == SearchedName {
				if len(RutaPasada) == 0 {
					// Último elemento, retornamos el inodo encontrado
					return entry.B_Inode
				}
				// Si hay más pasos en la ruta, seguimos recursivamente
				var nextInode Structs.Inode
				inodeOffset := int64(tempSuperblock.SB_Inode_Start + entry.B_Inode*int32(binary.Size(Structs.Inode{})))
				if err := Utilities.ReadObject(file, &nextInode, inodeOffset, buffer); err != nil {
					return -1
				}
				return BuscarInodoRuta(RutaPasada, nextInode, file, tempSuperblock, buffer)
			}
		}
	}

	return -1
}


// Función auxiliar para extraer el último elemento de un slice
func pop(s *[]string) string {
	lastIndex := len(*s) - 1
	last := (*s)[lastIndex]
	*s = (*s)[:lastIndex]
	return last
}

func Mkdir(path string, p bool, buffer *bytes.Buffer) {
	fmt.Fprintln(buffer, "=========== MKDIR ===========")

	// 1. Validar sesión activa
	if !isUserLoggedIn() {
		fmt.Fprintln(buffer, "Error MKDIR: No hay una sesión activa.")
		return
	}

	// 2. Obtener partición montada con sesión activa
	partitions := DiskManagement.GetMountedPartitions()
	var filepath, id string
	found := false
	for _, parts := range partitions {
		for _, part := range parts {
			if part.LoggedIn {
				filepath = part.Path
				id = part.ID
				found = true
				break
			}
		}
		if found {
			break
		}
	}
	if !found {
		fmt.Fprintln(buffer, "Error MKDIR: No se encontró partición activa.")
		return
	}

	// 3. Abrir disco
	file, err := Utilities.OpenFile(filepath, buffer)
	if err != nil {
		fmt.Fprintln(buffer, "Error MKDIR: No se pudo abrir el archivo.")
		return
	}
	defer file.Close()

	// 4. Leer MBR y Superblock
	var mbr Structs.MRB
	if err := Utilities.ReadObject(file, &mbr, 0, buffer); err != nil {
		fmt.Fprintln(buffer, "Error MKDIR: No se pudo leer el MBR.")
		return
	}

	// Buscar la partición correspondiente
	var index int = -1
	for i := 0; i < 4; i++ {
		if strings.Contains(string(mbr.MbrPartitions[i].ID[:]), id) && mbr.MbrPartitions[i].Status[0] == '1' {
			index = i
			break
		}
	}
	if index == -1 {
		fmt.Fprintln(buffer, "Error MKDIR: Partición no montada.")
		return
	}

	var sb Structs.Superblock
	if err := Utilities.ReadObject(file, &sb, int64(mbr.MbrPartitions[index].Start), buffer); err != nil {
		fmt.Fprintln(buffer, "Error MKDIR: No se pudo leer el Superblock.")
		return
	}

	// 5. Procesar ruta
	pathParts := strings.Split(strings.Trim(path, "/"), "/")
	if len(pathParts) == 0 {
		fmt.Fprintln(buffer, "Error MKDIR: Ruta inválida.")
		return
	}

	// 6. Buscar carpeta padre y crear carpetas si no existen
	rootInode := Structs.Inode{}
	if err := Utilities.ReadObject(file, &rootInode, int64(sb.SB_Inode_Start), buffer); err != nil {
		fmt.Fprintln(buffer, "Error MKDIR: No se pudo leer el Inodo raíz.")
		return
	}

	err = crearCarpetas(pathParts, rootInode, file, sb, p, buffer)
	if err != nil {
		fmt.Fprintln(buffer, "Error MKDIR:", err)
		return
	}

	fmt.Fprintln(buffer, "Carpeta creada exitosamente:", path)
}

func crearCarpetas(pathParts []string, currentInode Structs.Inode, file *os.File, sb Structs.Superblock, p bool, buffer *bytes.Buffer) error {
	fmt.Fprintf(buffer, "DEBUG crearCarpetas: p=%v, pathParts=%v\n", p, pathParts)
	for i, part := range pathParts {
		found := false
		var nextInode Structs.Inode
		var nextInodeIndex int32 = -1

		// Buscar si la carpeta ya existe
		for _, block := range currentInode.IN_Block {
			if block == -1 {
				continue
			}
			var folder Structs.FolderBlock
			blockOffset := int64(sb.SB_Block_Start) + int64(block)*int64(binary.Size(Structs.FolderBlock{}))
			if err := Utilities.ReadObject(file, &folder, blockOffset, buffer); err != nil {
				return err
			}
			for _, entry := range folder.B_Content {
				name := strings.Trim(string(entry.B_Name[:]), "\x00")
				if name == part && entry.B_Inode != -1 {
					// Carpeta encontrada
					found = true
					nextInodeIndex = entry.B_Inode
					inodeOffset := int64(sb.SB_Inode_Start) + int64(nextInodeIndex)*int64(binary.Size(Structs.Inode{}))
					if err := Utilities.ReadObject(file, &nextInode, inodeOffset, buffer); err != nil {
						return err
					}
					break
				}
			}
			if found {
				break
			}
		}

		if found {
			// Continuamos con la siguiente carpeta
			currentInode = nextInode
			continue
		} else {
			// No se encontró: ¿puede crear?
			if !p && i != len(pathParts)-1 {
				fmt.Fprintf(buffer, "DEBUG: no se encontró carpeta '%s' y -p no está activado\n", part)
				return fmt.Errorf("no existe la carpeta '%s' y no se especificó -p", part)
			}

			// Crear carpeta
			newInodeIndex, newBlockIndex, err := reservarInodoYBloque(file, sb, buffer)
			if err != nil {
				return fmt.Errorf("error al reservar espacio para carpeta '%s': %v", part, err)
			}

			// Crear el nuevo FolderBlock
			var newFolder Structs.FolderBlock
			copy(newFolder.B_Content[0].B_Name[:], ".")
			newFolder.B_Content[0].B_Inode = newInodeIndex
			copy(newFolder.B_Content[1].B_Name[:], "..")
			newFolder.B_Content[1].B_Inode = buscarIndiceInodo(currentInode, file, sb, buffer)

			// Guardar el bloque
			blockOffset := int64(sb.SB_Block_Start + newBlockIndex*int32(binary.Size(Structs.FolderBlock{})))
			if err := Utilities.WriteObject(file, newFolder, blockOffset, buffer); err != nil {
				return fmt.Errorf("error al escribir folder block: %v", err)
			}

			// Crear el nuevo Inode
			var newInode Structs.Inode
			newInode.IN_Uid = 1
			newInode.IN_Gid = 1
			newInode.IN_Size = 0
			copy(newInode.IN_Perm[:], "664")
			for i := 0; i < 15; i++ {
				newInode.IN_Block[i] = -1
			}
			newInode.IN_Block[0] = newBlockIndex

			inodeOffset := int64(sb.SB_Inode_Start + newInodeIndex*int32(binary.Size(Structs.Inode{})))
			if err := Utilities.WriteObject(file, newInode, inodeOffset, buffer); err != nil {
				return fmt.Errorf("error al escribir inode: %v", err)
			}

			// Insertar la nueva carpeta en el FolderBlock del inodo padre
			if err := agregarEntradaAFolderBlock(&currentInode, part, newInodeIndex, file, sb, buffer); err != nil {
				return fmt.Errorf("error al agregar entrada de carpeta '%s': %v", part, err)
			}

			// Actualizar inodo padre en disco
			currentInodeIndex := buscarIndiceInodo(currentInode, file, sb, buffer)
			if currentInodeIndex == -1 {
				return fmt.Errorf("no se pudo encontrar índice del inodo padre")
			}
			inodeOffset = int64(sb.SB_Inode_Start + currentInodeIndex*int32(binary.Size(Structs.Inode{})))
			if err := Utilities.WriteObject(file, currentInode, inodeOffset, buffer); err != nil {
				return fmt.Errorf("error al actualizar inodo padre: %v", err)
			}

			// Continuar con el nuevo inodo creado
			currentInode = newInode
		}
	}
	return nil
}

func reservarInodoYBloque(file *os.File, sb Structs.Superblock, buffer *bytes.Buffer) (int32, int32, error) {
	// Buscar primer bit libre en el bitmap de inodos
	var inodoLibre int32 = -1
	for i := int32(0); i < sb.SB_Inodes_Count; i++ {
		var bit byte
		if err := Utilities.ReadObject(file, &bit, int64(sb.SB_Bm_Inode_Start+i), buffer); err != nil {
			return -1, -1, err
		}
		if bit == 0 {
			inodoLibre = i
			break
		}
	}
	if inodoLibre == -1 {
		return -1, -1, fmt.Errorf("no hay inodos disponibles")
	}
	if err := Utilities.WriteObject(file, byte(1), int64(sb.SB_Bm_Inode_Start+inodoLibre), buffer); err != nil {
		return -1, -1, err
	}

	// Buscar primer bit libre en el bitmap de bloques
	var bloqueLibre int32 = -1
	for i := int32(0); i < sb.SB_Blocks_Count; i++ {
		var bit byte
		if err := Utilities.ReadObject(file, &bit, int64(sb.SB_Bm_Block_Start+i), buffer); err != nil {
			return -1, -1, err
		}
		if bit == 0 {
			bloqueLibre = i
			break
		}
	}
	if bloqueLibre == -1 {
		return -1, -1, fmt.Errorf("no hay bloques disponibles")
	}
	if err := Utilities.WriteObject(file, byte(1), int64(sb.SB_Bm_Block_Start+bloqueLibre), buffer); err != nil {
		return -1, -1, err
	}

	return inodoLibre, bloqueLibre, nil
}

func buscarIndiceInodo(target Structs.Inode, file *os.File, sb Structs.Superblock, buffer *bytes.Buffer) int32 {
	for i := int32(0); i < sb.SB_Inodes_Count; i++ {
		var temp Structs.Inode
		offset := int64(sb.SB_Inode_Start + i*int32(binary.Size(Structs.Inode{})))
		if err := Utilities.ReadObject(file, &temp, offset, buffer); err != nil {
			continue
		}

		if temp.IN_Uid == target.IN_Uid &&
			temp.IN_Gid == target.IN_Gid &&
			temp.IN_Size == target.IN_Size &&
			string(temp.IN_Perm[:]) == string(target.IN_Perm[:]) &&
			temp.IN_Block[0] == target.IN_Block[0] {
			return i
		}
	}
	return -1
}

func agregarEntradaAFolderBlock(parentInode *Structs.Inode, name string, inodeIndex int32, file *os.File, sb Structs.Superblock, buffer *bytes.Buffer) error {
	for i := 0; i < 15; i++ {
		block := parentInode.IN_Block[i]
		if block == -1 {
			// Crear un nuevo FolderBlock si no existe
			newBlockIndex := int32(-1)
			for j := int32(0); j < sb.SB_Blocks_Count; j++ {
				var bit byte
				if err := Utilities.ReadObject(file, &bit, int64(sb.SB_Bm_Block_Start+j), buffer); err != nil {
					return err
				}
				if bit == 0 {
					newBlockIndex = j
					break
				}
			}
			if newBlockIndex == -1 {
				return fmt.Errorf("no hay bloques disponibles para nueva carpeta")
			}

			// Marcar el bloque como usado
			if err := Utilities.WriteObject(file, byte(1), int64(sb.SB_Bm_Block_Start+newBlockIndex), buffer); err != nil {
				return err
			}

			var newFolder Structs.FolderBlock
			copy(newFolder.B_Content[0].B_Name[:], name)
			newFolder.B_Content[0].B_Inode = inodeIndex

			blockOffset := int64(sb.SB_Block_Start + newBlockIndex*int32(binary.Size(Structs.FolderBlock{})))
			if err := Utilities.WriteObject(file, newFolder, blockOffset, buffer); err != nil {
				return err
			}

			// Actualizar el inodo con el nuevo bloque
			parentInode.IN_Block[i] = newBlockIndex

			// DEPURACIÓN
			fmt.Fprintf(buffer, "DEBUG: Nuevo FolderBlock creado en bloque %d\n", newBlockIndex)
			fmt.Fprintf(buffer, "DEBUG: Insertando nombre '%s' con inodo %d en posición 0\n", name, inodeIndex)
			fmt.Fprintf(buffer, "DEBUG: Bytes del nombre: %v\n", newFolder.B_Content[0].B_Name)

			return nil
		} else {
			// Revisar si hay espacio en el FolderBlock existente
			var folder Structs.FolderBlock
			blockOffset := int64(sb.SB_Block_Start + block*int32(binary.Size(Structs.FolderBlock{})))
			if err := Utilities.ReadObject(file, &folder, blockOffset, buffer); err != nil {
				return err
			}
			for j := 0; j < 4; j++ {
				if folder.B_Content[j].B_Inode == -1 {
					copy(folder.B_Content[j].B_Name[:], name)
					folder.B_Content[j].B_Inode = inodeIndex
					if err := Utilities.WriteObject(file, folder, blockOffset, buffer); err != nil {
						return err
					}

					// DEPURACIÓN
					fmt.Fprintf(buffer, "DEBUG: Insertando nombre '%s' en FolderBlock existente en bloque %d, posición %d\n", name, block, j)
					fmt.Fprintf(buffer, "DEBUG: Inodo asociado: %d\n", inodeIndex)
					fmt.Fprintf(buffer, "DEBUG: Bytes del nombre: %v\n", folder.B_Content[j].B_Name)

					return nil
				}
			}
		}
	}
	return fmt.Errorf("no hay espacio disponible para agregar entrada en inodo padre")
}

func Mkfile(path string, p bool, content string, buffer *bytes.Buffer) {
	fmt.Fprintln(buffer, "=========== MKFILE ===========")

	if !isUserLoggedIn() {
		fmt.Fprintln(buffer, "Error MKFILE: No hay una sesión activa.")
		return
	}

	partitions := DiskManagement.GetMountedPartitions()
	var filepath, id string
	found := false
	for _, parts := range partitions {
		for _, part := range parts {
			if part.LoggedIn {
				filepath = part.Path
				id = part.ID
				found = true
				break
			}
		}
		if found {
			break
		}
	}

	if !found {
		fmt.Fprintln(buffer, "Error MKFILE: No se encontró partición activa.")
		return
	}

	file, err := Utilities.OpenFile(filepath, buffer)
	if err != nil {
		fmt.Fprintln(buffer, "Error MKFILE: No se pudo abrir el archivo.")
		return
	}
	defer file.Close()

	var mbr Structs.MRB
	if err := Utilities.ReadObject(file, &mbr, 0, buffer); err != nil {
		fmt.Fprintln(buffer, "Error MKFILE: No se pudo leer el MBR.")
		return
	}

	var index int = -1
	for i := 0; i < 4; i++ {
		if strings.Contains(string(mbr.MbrPartitions[i].ID[:]), id) && mbr.MbrPartitions[i].Status[0] == '1' {
			index = i
			break
		}
	}
	if index == -1 {
		fmt.Fprintln(buffer, "Error MKFILE: Partición no montada.")
		return
	}

	var sb Structs.Superblock
	if err := Utilities.ReadObject(file, &sb, int64(mbr.MbrPartitions[index].Start), buffer); err != nil {
		fmt.Fprintln(buffer, "Error MKFILE: No se pudo leer el Superblock.")
		return
	}

	pathParts := strings.Split(strings.Trim(path, "/"), "/")
	fileName := pathParts[len(pathParts)-1]
	parentDirs := pathParts[:len(pathParts)-1]

	rootInode := Structs.Inode{}
	if err := Utilities.ReadObject(file, &rootInode, int64(sb.SB_Inode_Start), buffer); err != nil {
		fmt.Fprintln(buffer, "Error MKFILE: No se pudo leer el Inodo raíz.")
		return
	}

	err = crearCarpetas(parentDirs, rootInode, file, sb, p, buffer)
	if err != nil {
		fmt.Fprintln(buffer, "Error MKFILE:", err)
		return
	}

	current := Structs.Inode{}
	entryInodeIndex := int32(0)
	if err := Utilities.ReadObject(file, &current, int64(sb.SB_Inode_Start), buffer); err != nil {
		fmt.Fprintln(buffer, "Error MKFILE: No se pudo leer el inodo raíz.")
		return
	}
	for _, part := range parentDirs {
		found := false
		for _, blk := range current.IN_Block {
			if blk == -1 {
				continue
			}
			var folder Structs.FolderBlock
			offset := int64(sb.SB_Block_Start + blk*int32(binary.Size(Structs.FolderBlock{})))
			if err := Utilities.ReadObject(file, &folder, offset, buffer); err != nil {
				continue
			}
			for _, entry := range folder.B_Content {
				name := strings.Trim(string(entry.B_Name[:]), "\x00")
				if name == part {
					entryInodeIndex = entry.B_Inode
					offset := int64(sb.SB_Inode_Start + entry.B_Inode*int32(binary.Size(Structs.Inode{})))
					if err := Utilities.ReadObject(file, &current, offset, buffer); err == nil {
						found = true
						break
					}
				}
			}
			if found {
				break
			}
		}
		if !found {
			fmt.Fprintln(buffer, "Error MKFILE: Carpeta padre no encontrada")
			return
		}
	}

	if !tienePermisoEscritura(current, User.Data.GetUID(), User.Data.GetGID()) {
		fmt.Fprintln(buffer, "Error MKFILE: El usuario no tiene permiso de escritura en la carpeta padre.")
		return
	}

	for _, blk := range current.IN_Block {
		if blk == -1 {
			continue
		}
		var folder Structs.FolderBlock
		offset := int64(sb.SB_Block_Start + blk*int32(binary.Size(Structs.FolderBlock{})))
		if err := Utilities.ReadObject(file, &folder, offset, buffer); err != nil {
			continue
		}
		for _, entry := range folder.B_Content {
			name := strings.Trim(string(entry.B_Name[:]), "\x00")
			if name == fileName {
				fmt.Fprintf(buffer, "Error MKFILE: Ya existe un archivo con el nombre '%s'\n", fileName)
				return
			}
		}
	}

	inodeIdx, blockIdx, err := reservarInodoYBloque(file, sb, buffer)
	if err != nil {
		fmt.Fprintln(buffer, "Error MKFILE:", err)
		return
	}

	var newInode Structs.Inode
	newInode.IN_Uid = int32(User.Data.GetUID())
	newInode.IN_Gid = int32(User.Data.GetGID())
	newInode.IN_Size = int32(len(content))
	copy(newInode.IN_Perm[:], "664")
	now := time.Now().Format("2006-01-02 15:04:05")
	copy(newInode.IN_Atime[:], now)
	copy(newInode.IN_Ctime[:], now)
	copy(newInode.IN_Mtime[:], now)
	for i := 0; i < 15; i++ {
		newInode.IN_Block[i] = -1
	}
	newInode.IN_Block[0] = blockIdx

	offsetInode := int64(sb.SB_Inode_Start + inodeIdx*int32(binary.Size(Structs.Inode{})))
	if err := Utilities.WriteObject(file, newInode, offsetInode, buffer); err != nil {
		fmt.Fprintln(buffer, "Error MKFILE: No se pudo escribir el inodo del archivo")
		return
	}

	var fileblock Structs.FileBlock
	copy(fileblock.B_Content[:], content)
	offsetBlock := int64(sb.SB_Block_Start + blockIdx*int32(binary.Size(Structs.FileBlock{})))
	if err := Utilities.WriteObject(file, fileblock, offsetBlock, buffer); err != nil {
		fmt.Fprintln(buffer, "Error MKFILE: No se pudo escribir el contenido del archivo")
		return
	}

	if err := agregarEntradaAFolderBlock(&current, fileName, inodeIdx, file, sb, buffer); err != nil {
		fmt.Fprintln(buffer, "Error MKFILE: No se pudo agregar la entrada al folder padre")
		return
	}

	// <-- CAMBIO AQUI
	offsetCurrent := int64(sb.SB_Inode_Start + entryInodeIndex*int32(binary.Size(Structs.Inode{})))
	if err := Utilities.WriteObject(file, current, offsetCurrent, buffer); err != nil {
		fmt.Fprintln(buffer, "Error MKFILE: No se pudo escribir el inodo padre actualizado")
		return
	}

	fmt.Fprintln(buffer, "Archivo creado exitosamente:", path)

	//fmt.Fprintln(buffer, "----------------------------------------------------------------------------")
	//fmt.Fprintf(buffer, "DEBUG: Archivo creado en inodo %d y bloque %d\n", inodeIdx, blockIdx)

	var createdInode Structs.Inode
	Utilities.ReadObject(file, &createdInode, int64(sb.SB_Inode_Start+inodeIdx*int32(binary.Size(Structs.Inode{}))), buffer)
	Structs.PrintInode(createdInode)

	var createdBlock Structs.FileBlock
	Utilities.ReadObject(file, &createdBlock, int64(sb.SB_Block_Start+blockIdx*int32(binary.Size(Structs.FileBlock{}))), buffer)
	Structs.PrintFileblock(createdBlock, buffer)

	//fmt.Fprintln(buffer, "\n--- DEBUG: FolderBlock del directorio padre ---")
	for _, blk := range current.IN_Block {
		if blk == -1 {
			continue
		}
		var folder Structs.FolderBlock
		offset := int64(sb.SB_Block_Start + blk*int32(binary.Size(Structs.FolderBlock{})))
		if err := Utilities.ReadObject(file, &folder, offset, buffer); err != nil {
			continue
		}
		Structs.PrintFolderblock(folder)
	}
}


func buscarInodoPorRuta(path string, file *os.File, sb Structs.Superblock, buffer *bytes.Buffer) int32 {
	ruta := strings.Split(strings.Trim(path, "/"), "/")
	inodoActual := int32(0) // raíz

	for _, nombre := range ruta {
		fmt.Fprintf(buffer, "DEBUG: buscando '%s' en inodo %d\n", nombre, inodoActual)

		var inode Structs.Inode
		offsetInodo := int64(sb.SB_Inode_Start + inodoActual*int32(binary.Size(Structs.Inode{})))
		if err := Utilities.ReadObject(file, &inode, offsetInodo, buffer); err != nil {
			fmt.Fprintf(buffer, "Error al leer inodo %d\n", inodoActual)
			return -1
		}

		encontrado := false
		for i := 0; i < 12; i++ { // Solo los bloques directos
			bloque := inode.IN_Block[i]
			if bloque == -1 {
				continue
			}

			var folder Structs.FolderBlock
			offsetBloque := int64(sb.SB_Block_Start + bloque*int32(binary.Size(Structs.FolderBlock{})))
			if err := Utilities.ReadObject(file, &folder, offsetBloque, buffer); err != nil {
				fmt.Fprintf(buffer, "Error al leer folder block %d\n", bloque)
				return -1
			}

			for _, entrada := range folder.B_Content {
				nombreEntrada := strings.Trim(string(entrada.B_Name[:]), "\x00")
				fmt.Fprintf(buffer, "DEBUG: comparando con entrada '%s'\n", nombreEntrada)
				if nombreEntrada == nombre && entrada.B_Inode != -1 {
					inodoActual = entrada.B_Inode
					encontrado = true
					break
				}
			}
			if encontrado {
				break
			}
		}

		if !encontrado {
			fmt.Fprintf(buffer, "No se encontró el elemento '%s'\n", nombre)
			return -1
		}
	}

	return inodoActual
}

//-------------metodos adicionales para reportes----------------
func BuscarInodoPorRutaREPORTE(path string, file *os.File, sb Structs.Superblock, buffer *bytes.Buffer) int32 {
	ruta := strings.Split(strings.Trim(path, "/"), "/")
	inodoActual := int32(0) // raíz

	for _, nombre := range ruta {
		fmt.Fprintf(buffer, "DEBUG: buscando '%s' en inodo %d\n", nombre, inodoActual)

		var inode Structs.Inode
		offsetInodo := int64(sb.SB_Inode_Start + inodoActual*int32(binary.Size(Structs.Inode{})))
		if err := Utilities.ReadObject(file, &inode, offsetInodo, buffer); err != nil {
			fmt.Fprintf(buffer, "Error al leer inodo %d\n", inodoActual)
			return -1
		}

		encontrado := false
		for i := 0; i < 12; i++ { // Solo los bloques directos
			bloque := inode.IN_Block[i]
			if bloque == -1 {
				continue
			}

			var folder Structs.FolderBlock
			offsetBloque := int64(sb.SB_Block_Start + bloque*int32(binary.Size(Structs.FolderBlock{})))
			if err := Utilities.ReadObject(file, &folder, offsetBloque, buffer); err != nil {
				fmt.Fprintf(buffer, "Error al leer folder block %d\n", bloque)
				return -1
			}

			for _, entrada := range folder.B_Content {
				nombreEntrada := strings.Trim(string(entrada.B_Name[:]), "\x00")
				fmt.Fprintf(buffer, "DEBUG: comparando con entrada '%s'\n", nombreEntrada)
				if nombreEntrada == nombre && entrada.B_Inode != -1 {
					inodoActual = entrada.B_Inode
					encontrado = true
					break
				}
			}
			if encontrado {
				break
			}
		}

		if !encontrado {
			fmt.Fprintf(buffer, "No se encontró el elemento '%s'\n", nombre)
			return -1
		}
	}

	return inodoActual
}

func IsUserLoggedInREPORTE() bool {
	fmt.Println("DEBUG: Usuario actual:", User.Data.GetIDUsuario())
	fmt.Println("DEBUG: Partición actual:", User.Data.GetIDPartition())
	return User.Data.GetIDUsuario() != "" && User.Data.GetIDPartition() != ""
}

func tienePermisoEscritura(inodo Structs.Inode, uid, gid int) bool {
	permStr := strings.Trim(string(inodo.IN_Perm[:]), "\x00")
	if len(permStr) != 3 {
		return false
	}

	if uid == 0 {
		return true // root siempre tiene permiso
	}

	u, _ := strconv.Atoi(string(permStr[0]))
	g, _ := strconv.Atoi(string(permStr[1]))
	o, _ := strconv.Atoi(string(permStr[2]))

	if uid == int(inodo.IN_Uid) {
		return u&2 != 0
	} else if gid == int(inodo.IN_Gid) {
		return g&2 != 0
	} else {
		return o&2 != 0
	}
}
