
package DiskManagement

import (
	"fmt"
	"math/rand"
	"proyecto1/Structs"
	"proyecto1/Utilities"
	"time"
)

func Mkdisk(size int, fit string, unit string, path string) {
	fmt.Println("======INICIO MKDISK======")
	fmt.Println("Size:", size)
	fmt.Println("Fit:", fit)
	fmt.Println("Unit:", unit)
	fmt.Println("Path:", path)

	// Validar fit bf/ff/wf
	if fit != "bf" && fit != "wf" && fit != "ff" {
		fmt.Println("Error: Fit debe ser bf, wf or ff")
		return
	}

	// Validar size > 0
	if size <= 0 {
		fmt.Println("Error: Size debe ser mayo a  0")
		return
	}

	// Validar unidar k - m
	if unit != "k" && unit != "m" {
		fmt.Println("Error: Las unidades validas son k o m")
		return
	}

	// Create file
	err := Utilities.CreateFile(path)
	if err != nil {
		fmt.Println("Error: ", err)
		return
	}

	/*
		Si el usuario especifica unit = "k" (Kilobytes), el tamaño se multiplica por 1024 para convertirlo a bytes.
		Si el usuario especifica unit = "m" (Megabytes), el tamaño se multiplica por 1024 * 1024 para convertirlo a MEGA bytes.
	*/
	// Asignar tamanio
	if unit == "k" {
		size = size * 1024
	} else {
		size = size * 1024 * 1024
	}

	// Open bin file
	file, err := Utilities.OpenFile(path)
	if err != nil {
		return
	}

	// Escribir los 0 en el archivo

	// create array of byte(0)
	for i := 0; i < size; i++ {
		err := Utilities.WriteObject(file, byte(0), int64(i))
		if err != nil {
			fmt.Println("Error: ", err)
		}
	}

	// Crear MRB
	var newMRB Structs.MRB
	newMRB.MbrSize = int32(size)
	newMRB.Signature = rand.Int31() // Numero random rand.Int31() genera solo números no negativos
	copy(newMRB.Fit[:], fit)

	// Obtener la fecha del sistema en formato YYYY-MM-DD
	currentTime := time.Now()
	formattedDate := currentTime.Format("2006-01-02")
	copy(newMRB.CreationDate[:], formattedDate)

	/*
		newMRB.CreationDate[0] = '2'
		newMRB.CreationDate[1] = '0'
		newMRB.CreationDate[2] = '2'
		newMRB.CreationDate[3] = '4'
		newMRB.CreationDate[4] = '-'
		newMRB.CreationDate[5] = '0'
		newMRB.CreationDate[6] = '8'
		newMRB.CreationDate[7] = '-'
		newMRB.CreationDate[8] = '0'
		newMRB.CreationDate[9] = '8'
	*/

	// Escribir el archivo
	if err := Utilities.WriteObject(file, newMRB, 0); err != nil {
		return
	}

	var TempMBR Structs.MRB
	// Leer el archivo
	if err := Utilities.ReadObject(file, &TempMBR, 0); err != nil {
		return
	}

	// Print object
	Structs.PrintMBR(TempMBR)

	// Cerrar el archivo
	defer file.Close()

	fmt.Println("======FIN MKDISK======")
}

func CreateLogicalPartition(path string, size int, fit string, name string) {
	fmt.Println("======INICIO CREATE LOGICAL PARTITION======")

	// Abrir archivo del disco
	file, err := Utilities.OpenFile(path)
	if err != nil {
		fmt.Println("Error al abrir el disco:", err)
		return
	}
	defer file.Close()

	// Leer el MBR
	var mbr Structs.MRB
	if err := Utilities.ReadObject(file, &mbr, 0); err != nil {
		fmt.Println("Error al leer el MBR:", err)
		return
	}

	// Buscar partición extendida
	var extendedStart int64 = -1
	for _, partition := range mbr.MbrPartitions {
		if partition.PartType == 'E' {
			extendedStart = partition.PartStart
			break
		}
	}

	if extendedStart == -1 {
		fmt.Println("Error: No se encontró una partición extendida.")
		return
	}

	// Buscar el primer espacio libre en la partición extendida
	file.Seek(extendedStart, 0)
	var ebr Structs.EBR
	for {
		err := Utilities.ReadObject(file, &ebr, extendedStart)
		if err != nil {
			fmt.Println("Error al leer el EBR:", err)
			return
		}

		// Si encontramos un espacio libre o el último EBR, salimos del bucle
		if ebr.PartStart == 0 || ebr.PartNext == -1 {
			break
		}

		// Mover al siguiente EBR
		extendedStart = ebr.PartNext
	}

	// Crear nuevo EBR
	var newEBR Structs.EBR
	newEBR.PartMount = 0
	newEBR.PartFit = fit[0] // Guardamos solo la primera letra ('B', 'F' o 'W')
	newEBR.PartStart = extendedStart + 1
	newEBR.PartS = int64(size)
	newEBR.PartNext = -1
	copy(newEBR.PartName[:], name)

	// Escribir nuevo EBR en disco
	if err := Utilities.WriteObject(file, newEBR, extendedStart); err != nil {
		fmt.Println("Error al escribir el nuevo EBR:", err)
		return
	}

	fmt.Println("Partición lógica creada exitosamente en", path)
	fmt.Println("======FIN CREATE LOGICAL PARTITION======")
}
