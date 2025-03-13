
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

		indexInode := buscarStart(filePath, file, tempSuperblock, buffer)
		if indexInode == -1 {
			fmt.Printf("Error: No se pudo encontrar el archivo %s\n", filePath)
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
	ParticionesMount := DiskManagement.GetMountedPartitions()

	for _, partitions := range ParticionesMount {
		for _, partition := range partitions {
			if partition.LoggedIn {
				return true
			}
		}
	}

	return false
}

// Función para verificar si el usuario tiene permisos
func tienePermiso(buffer *bytes.Buffer) bool {
	ParticionesMount := DiskManagement.GetMountedPartitions()
	var filepath string
	var id string

	for _, partitions := range ParticionesMount {
		for _, partition := range partitions {
			// Verifica si alguna partición tiene un usuario logueado
			if partition.LoggedIn {
				filepath = partition.Path
				id = partition.ID
				break
			}
		}
	}

	file, err := Utilities.OpenFile(filepath, buffer)
	if err != nil {
		fmt.Println("Error: No se pudo abrir el archivo:", err)
		return false
	}
	defer file.Close()

	var TempMBR Structs.MRB

	if err := Utilities.ReadObject(file, &TempMBR, 0, buffer); err != nil {
		fmt.Println("Error: No se pudo leer el MBR:", err)
		return false
	}

	var index int = -1

	for i := 0; i < 4; i++ {
		if TempMBR.MbrPartitions[i].Size != 0 {
			if strings.Contains(string(TempMBR.MbrPartitions[i].ID[:]), id) {
				if TempMBR.MbrPartitions[i].Status[0] == '1' {
					index = i
				} else {
					return false
				}
				break
			}
		}
	}

	if index == -1 {
		return false
	}

	var tempSuperblock Structs.Superblock
	if err := Utilities.ReadObject(file, &tempSuperblock, int64(TempMBR.MbrPartitions[index].Start), buffer); err != nil {
		return false
	}

	indexInode := buscarStart("/user.txt", file, tempSuperblock, buffer)

	var crrInode Structs.Inode

	if err := Utilities.ReadObject(file, &crrInode, int64(tempSuperblock.SB_Inode_Start+indexInode*int32(binary.Size(Structs.Inode{}))), buffer); err != nil {
		return false
	}

	perm := string(crrInode.IN_Perm[:])
	return strings.Contains(perm, "664")
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
	SearchedName := strings.Replace(pop(&RutaPasada), " ", "", -1)

	for _, block := range Inode.IN_Block {
		if block != -1 {
			if len(RutaPasada) == 0 { // Caso donde encontramos el archivo
				var fileblock Structs.FileBlock
				if err := Utilities.ReadObject(file, &fileblock, int64(tempSuperblock.SB_Block_Start+block*int32(binary.Size(Structs.FileBlock{}))), buffer); err != nil {
					return -1
				}

				//Structs.PrintFileblock(fileblock) // Imprime el contenido del FileBlock
				return 1
			} else {
				// En este caso seguimos buscando en los bloques de carpetas
				var crrFolderBlock Structs.FolderBlock
				if err := Utilities.ReadObject(file, &crrFolderBlock, int64(tempSuperblock.SB_Block_Start+block*int32(binary.Size(Structs.FolderBlock{}))), buffer); err != nil {
					return -1
				}

				for _, folder := range crrFolderBlock.B_Content {
					if strings.Contains(string(folder.B_Name[:]), SearchedName) {
						var NextInode Structs.Inode
						if err := Utilities.ReadObject(file, &NextInode, int64(tempSuperblock.SB_Inode_Start+folder.B_Inode*int32(binary.Size(Structs.Inode{}))), buffer); err != nil {
							return -1
						}

						return BuscarInodoRuta(RutaPasada, NextInode, file, tempSuperblock, buffer)
					}
				}
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