package Structs

import (
	"fmt"
)

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

type Partition struct {
	PartStatus      [1]byte     // Indica si la partición está montada o no
	PartType        [1]byte     // 'P' para primaria, 'E' para extendida
	PartFit         [1]byte     // Tipo de ajuste: 'B' (Best), 'F' (First), 'W' (Worst)
	PartStart       int32    // Byte donde inicia la partición en el disco
	PartS           int32    // Tamaño total de la partición en bytes
	PartName        [16]byte // Nombre de la partición
	PartCorrelative int32    // Número correlativo, inicia en -1 y se incrementa al montar
	PartID          [4]byte  // ID de la partición generada al montar
}

func PrintPartition(data Partition) {
	fmt.Printf("Nombre: %s, Tipo: %s, Inicio: %d, Tamaño: %d, Estado: %s, ID: %s, Ajuste: %s, Correlativo: %d\n",
		string(data.PartName[:]), string(data.PartType[:]), data.PartStart, data.PartSize, string(data.PartStatus[:]),
		string(data.PartId[:]), string(data.PartFit[:]), data.PartCorrelative)
}

// Definir estructura EBR
type EBR struct {
	PartMount [1]byte     // Indica si la partición está montada
	PartFit   [1]byte     // Tipo de ajuste: 'B', 'F', 'W'
	PartStart int32    // Byte donde inicia la partición
	PartS     int32    // Tamaño de la partición
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