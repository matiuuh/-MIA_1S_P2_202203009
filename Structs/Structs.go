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
	PartStatus      byte     // Indica si la partición está montada o no
	PartType        byte     // 'P' para primaria, 'E' para extendida
	PartFit         byte     // Tipo de ajuste: 'B' (Best), 'F' (First), 'W' (Worst)
	PartStart       int64    // Byte donde inicia la partición en el disco
	PartS           int64    // Tamaño total de la partición en bytes
	PartName        [16]byte // Nombre de la partición
	PartCorrelative int32    // Número correlativo, inicia en -1 y se incrementa al montar
	PartID          [4]byte  // ID de la partición generada al montar
}


// Definir estructura EBR
type EBR struct {
	PartMount byte     // Indica si la partición está montada
	PartFit   byte     // Tipo de ajuste: 'B', 'F', 'W'
	PartStart int64    // Byte donde inicia la partición
	PartS     int64    // Tamaño de la partición
	PartNext  int64    // Byte donde está el siguiente EBR (-1 si no hay)
	PartName  [16]byte // Nombre de la partición
}