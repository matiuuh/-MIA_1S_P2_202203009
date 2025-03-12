package Structs

import (
	"fmt"
	"bytes"
)

//*********MBR*********
//Definir estructura MBR
type MRB struct {
	MbrSize      int32
	CreationDate [10]byte
	Signature    int32
	Fit          [1]byte
	MbrPartitions [4]Partition // Agregar esta línea
}

// Imprimir MBR
func PrintMBR(data MRB) {

	fmt.Println(fmt.Sprintf("Creation Date: %s, Fit: %s, Size: %d, Signature: %d ",
		string(data.CreationDate[:]),
		string(data.Fit[:]),
		data.MbrSize,
		data.Signature))

}

//*********PARTITION*********
type Partition struct {
	Status      [1]byte     // Indica si la partición está montada o no
	Type        [1]byte     // 'P' para primaria, 'E' para extendida
	Fit         [1]byte     // Tipo de ajuste: 'B' (Best), 'F' (First), 'W' (Worst)
	Start       int32    // Byte donde inicia la partición en el disco
	Size        int32    // Tamaño total de la partición en bytes
	Name        [16]byte // Nombre de la partición
	Correlative int32    // Número correlativo, inicia en -1 y se incrementa al montar
	ID          [4]byte  // ID de la partición generada al montar
}

func PrintPartition(data Partition) {
	fmt.Printf("Nombre: %s, Tipo: %s, Inicio: %d, Tamaño: %d, Estado: %s, ID: %s, Ajuste: %s, Correlativo: %d\n",
		string(data.Name[:]), string(data.Type[:]), data.Start, data.Size, string(data.Status[:]),
		string(data.ID[:]), string(data.Fit[:]), data.Correlative)
}

//*********EBR*********
// Definir estructura EBR
type EBR struct {
	PartMount [1]byte     // Indica si la partición está montada
	PartFit   [1]byte     // Tipo de ajuste: 'B', 'F', 'W'
	PartStart int32    // Byte donde inicia la partición
	PartSize     int32    // Tamaño de la partición
	PartNext  int32    // Byte donde está el siguiente EBR (-1 si no hay)
	PartName  [16]byte // Nombre de la partición
}

func PrintEBR(data EBR) {
	fmt.Printf("Name: %s, fit: %c, start: %d, size: %d, next: %d, mount: %c\n",
		string(data.PartName[:]),
		data.PartFit,
		data.PartStart,
		data.PartSize,
		data.PartNext,
		data.PartMount)
}

//*********SUPER BLOQUE*********
// Superblock
type Superblock struct {
	SB_FileSystem_Type   int32		// Guarda el número que identifica el sistema de archivos utilizado
	SB_Inodes_Count      int32		// Guarda el número total de inodos
	SB_Blocks_Count      int32		// Guarda el número total de bloques
	SB_Free_Blocks_Count int32		// Contiene el número de bloques libres
	SB_Free_Inodes_Count int32		// Contiene el número de inodos libres
	SB_Mtime             [17]byte	// Última fecha en el que el sistema fue montado
	SB_Umtime            [17]byte	// Última fecha en que el sistema fue desmontado
	SB_Mnt_Count         int32		// Indica cuantas veces se ha montado el sistema
	SB_Magic             int32		// Valor que identifica al sistema de archivos, tendrá el valor 0xEF53
	SB_Inode_Size        int32		// Tamaño del inodo
	SB_Block_Size        int32		// Tamaño del bloque
	SB_Fist_Ino          int32		// Primer inodo libre (dirección del inodo)
	SB_First_Blo         int32		// Primer bloque libre (dirección del inodo)
	SB_Bm_Inode_Start    int32		// Guardará el inicio del bitmap de inodos
	SB_Bm_Block_Start    int32		// Guardará el inicio del bitmap de bloques
	SB_Inode_Start       int32		// Guardará el inicio de la tabla de inodos
	SB_Block_Start       int32		// Guardará el inicio de la tabla de bloques
}

func PrintSuperblock(sb Superblock) {
	fmt.Println("====== Superblock ======")
	fmt.Printf("S_filesystem_type: %d\n", sb.SB_FileSystem_Type)
	fmt.Printf("S_inodes_count: %d\n", sb.SB_Inodes_Count)
	fmt.Printf("S_blocks_count: %d\n", sb.SB_Blocks_Count)
	fmt.Printf("S_free_blocks_count: %d\n", sb.SB_Free_Blocks_Count)
	fmt.Printf("S_free_inodes_count: %d\n", sb.SB_Free_Inodes_Count)
	fmt.Printf("S_mtime: %s\n", string(sb.SB_Mtime[:]))
	fmt.Printf("S_umtime: %s\n", string(sb.SB_Umtime[:]))
	fmt.Printf("S_mnt_count: %d\n", sb.SB_Mnt_Count)
	fmt.Printf("S_magic: 0x%X\n", sb.SB_Magic)
	fmt.Printf("S_inode_size: %d\n", sb.SB_Inode_Size)
	fmt.Printf("S_block_size: %d\n", sb.SB_Block_Size)
	fmt.Printf("S_fist_ino: %d\n", sb.SB_Fist_Ino)
	fmt.Printf("S_first_blo: %d\n", sb.SB_First_Blo)
	fmt.Printf("S_bm_inode_start: %d\n", sb.SB_Bm_Inode_Start)
	fmt.Printf("S_bm_block_start: %d\n", sb.SB_Bm_Block_Start)
	fmt.Printf("S_inode_start: %d\n", sb.SB_Inode_Start)
	fmt.Printf("S_block_start: %d\n", sb.SB_Block_Start)
	fmt.Println("========================")
}

//*********INODO*********
// Inode
type Inode struct {
	IN_Uid   int32
	IN_Gid   int32
	IN_Size  int32
	IN_Atime [17]byte
	IN_Ctime [17]byte
	IN_Mtime [17]byte
	IN_Block [15]int32
	IN_Type  [1]byte
	IN_Perm  [3]byte
}

func PrintInode(inode Inode) {
	fmt.Println("====== Inode ======")
	fmt.Printf("I_uid: %d\n", inode.IN_Uid)
	fmt.Printf("I_gid: %d\n", inode.IN_Gid)
	fmt.Printf("I_size: %d\n", inode.IN_Size)
	fmt.Printf("I_atime: %s\n", string(inode.IN_Atime[:]))
	fmt.Printf("I_ctime: %s\n", string(inode.IN_Ctime[:]))
	fmt.Printf("I_mtime: %s\n", string(inode.IN_Mtime[:]))
	fmt.Printf("I_type: %s\n", string(inode.IN_Type[:]))
	fmt.Printf("I_perm: %s\n", string(inode.IN_Perm[:]))
	fmt.Printf("I_block: %v\n", inode.IN_Block)
	fmt.Println("===================")
}

//*********BLOQUE DE CARPETAS*********
// Bloque De Carpetas
type FolderBlock struct {
	B_Content [4]Content
}

type Content struct {
	B_Name  [12]byte
	B_Inode int32
}

func PrintFolderblock(folderblock FolderBlock) {
	fmt.Println("====== Folderblock ======")
	for i, content := range folderblock.B_Content {
		fmt.Printf("Content %d: Name: %s, Inodo: %d\n", i, string(content.B_Name[:]), content.B_Inode)
	}
	fmt.Println("=========================")
}

//*********BLOQUE DE ARCHIVOS*********
// Bloque De Archivos
type FileBlock struct {
	B_Content [64]byte
}

func PrintFileblock(fileblock FileBlock, buffer *bytes.Buffer) {
	fmt.Fprintf(buffer, "====== Fileblock ======\n")
	fmt.Fprintf(buffer, "%s\n", string(fileblock.B_Content[:]))
}

//*********BLOQUE DE APUNTADORES*********
type PointerBlock struct {
	B_Pointers [16]int32
}

func PrintPointerblock(pointerblock PointerBlock) {
	fmt.Println("====== Pointerblock ======")
	for i, pointer := range pointerblock.B_Pointers {
		fmt.Printf("Pointer %d: %d\n", i, pointer)
	}
	fmt.Println("=========================")
}