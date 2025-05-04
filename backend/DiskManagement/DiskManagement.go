
package DiskManagement

import (
	"fmt"
	"math/rand"
	"proyecto1/Structs"
	"proyecto1/Utilities"
	"time"
	"bytes"
	"strings"
	"encoding/binary"
)

//************ESTRUCTURAS************
type MountedPartition struct {
	Path     string
	Name   	 string
	ID       string
	Status   byte
	LoggedIn bool
}

var MountedPartitions = make(map[string][]MountedPartition)

//************************************

//********************PARTICIONES MONTADAS********************
// Función para imprimir las particiones montadas
func PrintMountedPartitions(path string, buffer *bytes.Buffer) {
	if len(MountedPartitions) == 0 {
		fmt.Println("No hay particiones montadas.")
		return
	}
	for DiscoID, partitions := range MountedPartitions {
		if DiscoID == path {
			fmt.Println("Disco:", DiscoID)
			fmt.Println("---------------------------")
			for _, Partition := range partitions {
				loginStatus := "No"
				if Partition.LoggedIn {
					loginStatus = "Sí"
				}
				fmt.Printf("Nombre: %v, ID: %v, Ruta: %v, Estado: %c, LoggedIn: %v\n",
					Partition.Name, Partition.ID, Partition.Path, Partition.Status, loginStatus)
			}
		}
		fmt.Println("---------------------------")
	}
}

// Función para obtener las particiones montadas
func GetMountedPartitions() map[string][]MountedPartition {
	return MountedPartitions
}
//****************************************************************

//********************COMANDOS********************
func Mkdisk(size int, fit string, unit string, path string, buffer *bytes.Buffer ) {
	fmt.Fprintf(buffer, "======INICIO MKDISK======\n")
	fmt.Println("Size:", size)
	fmt.Println("Fit:", fit)
	fmt.Println("Unit:", unit)
	fmt.Println("Path:", path)

	if unit == "" {
		unit = "m"
	}

	if fit == "" {
		fit = "ff"
	}	

	// Validar fit bf/ff/wf
	if fit != "bf" && fit != "wf" && fit != "ff" {
		fmt.Fprintf(buffer, "Error: Fit debe ser bf, wf or ff\n")
		return
	}

	// Validar size > 0
	if size <= 0 {
		fmt.Fprintf(buffer, "Error: Size debe ser mayo a  0\n")
		return
	}

	// Validar unidar k - m
	if unit != "k" && unit != "m" {
		fmt.Fprintf(buffer, "Error: Las unidades validas son k o m\n")
		return
	}

	// Validar la ruta (path)
	if path == "" {
		fmt.Fprintf(buffer, "Error MKDISK: La ruta del disco es obligatoria.\n")
		return
	}

	// Create file
	err := Utilities.CreateFile(path, buffer)
	if err != nil {
		fmt.Fprintf(buffer, "Error: ", err)
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
	file, err := Utilities.OpenFile(path, buffer)
	if err != nil {
		return
	}

	// Escribir los 0 en el archivo

	// create array of byte(0)
	for i := 0; i < size; i++ {
		err := Utilities.WriteObject(file, byte(0), int64(i), buffer)
		if err != nil {
			fmt.Fprintf(buffer, "Error: ", err)
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

	// Escribir el archivo
	if err := Utilities.WriteObject(file, newMRB, 0, buffer); err != nil {
		return
	}

	var TempMBR Structs.MRB
	// Leer el archivo
	if err := Utilities.ReadObject(file, &TempMBR, 0, buffer); err != nil {
		return
	}

	fmt.Println("----------------------------")
	// Print object
	Structs.PrintMBR(TempMBR)
	fmt.Println("----------------------------")

	// Cerrar el archivo
	defer file.Close()
	fmt.Fprintf(buffer, "Dico creado con exito en la ruta: %s con el tamaño: %d.\n", path, size)
	fmt.Fprintf(buffer, "======FIN MKDISK======")
}

func Rmdisk(path string, buffer *bytes.Buffer) {
	fmt.Fprintf(buffer, "======RMDISK======\n")
	// Validar la ruta (path)
	if path == "" {
		fmt.Fprintf(buffer, "Error RMDISK: La ruta del disco es obligatoria.\n")
		return
	}
	err := Utilities.DeleteFile(path, buffer)
	if err != nil {
		return
	}
	DeleteDiscWithPath(path, buffer)
	fmt.Fprintf(buffer, "Disco eliminado con éxito en la ruta: %s.\n", path)
}

func Fdisk(size int, path string, name string, unit string, type_ string, fit string, buffer *bytes.Buffer) {
	fmt.Fprintf(buffer, "======Start FDISK======\n")
	// Validar el tamaño (size)
	if size <= 0 {
		fmt.Fprintf(buffer, "Error FDISK: EL tamaño de la partición debe ser mayor que 0.\n")
		return
	}
	// Validar la unidad (unit)
	if unit != "b" && unit != "k" && unit != "m" {
		fmt.Fprintf(buffer, "Error FDISK: La unidad de tamaño debe ser Bytes, Kilobytes, Megabytes.\n")
		return
	}
	// Validar la ruta (path)
	if path == "" {
		fmt.Fprintf(buffer, "Error FDISK: La ruta del disco es obligatoria.\n")
		return
	}
	// Validar el tipo (type)
	if type_ != "p" && type_ != "e" && type_ != "l" {
		fmt.Fprintf(buffer, "Error FDISK: El tipo de partición debe ser Primaria, Extendida, Lógica.\n")
		return
	}
	// Validar el ajuste (fit)
	if fit != "b" && fit != "f" && fit != "w" {
		fmt.Fprintf(buffer, "Error FDISK: El ajuste de la partición debe ser b, f o w.\n")
		return
	}
	// Validar el nombre (name)
	if name == "" {
		fmt.Fprintf(buffer, "Error FDISK: El nombre de la partición es obligatorio.\n")
		return
	}

	// Convertir el tamaño a bytes
	if unit == "k" {
		size = size * 1024
	} else if unit == "m" {
		size = size * 1024 * 1024
	}

	// Abrir archivo binario
	file, err := Utilities.OpenFile(path, buffer)
	if err != nil {
		return
	}

	var TempMBR Structs.MRB
	if err := Utilities.ReadObject(file, &TempMBR, 0, buffer); err != nil {
		return
	}
	
	
	for i := 0; i < 4; i++ {
		if strings.Contains(string(TempMBR.MbrPartitions[i].Name[:]), name) {
			fmt.Fprintf(buffer, "Error FDISK: El nombre: %s ya está en uso en las particiones.\n", name)
			return
		}
	}

	var ContadorPrimaria, ContadorExtendida, TotalParticiones int
	var EspacioUtilizado int32 = 0
	var maxEnd int32 = int32(binary.Size(TempMBR))

	for i := 0; i < 4; i++ {
		if TempMBR.MbrPartitions[i].Size != 0 {
			TotalParticiones++
			EspacioUtilizado += TempMBR.MbrPartitions[i].Size
			end := TempMBR.MbrPartitions[i].Start + TempMBR.MbrPartitions[i].Size
			if end > maxEnd {
				maxEnd = end
			}
			if TempMBR.MbrPartitions[i].Type[0] == 'p' {
				ContadorPrimaria++
			} else if TempMBR.MbrPartitions[i].Type[0] == 'e' {
				ContadorExtendida++
			}
		}
	}

	if TotalParticiones >= 4 && type_ != "l" {
		fmt.Fprintf(buffer, "Error FDISK: No se pueden crear más de 4 particiones primarias o extendidas en total.\n")
		return
	}
	if type_ == "e" && ContadorExtendida > 0 {
		fmt.Fprintf(buffer, "Error FDISK: Solo se permite una partición extendida por disco.\n")
		return
	}
	if type_ == "l" && ContadorExtendida == 0 {
		fmt.Fprintf(buffer, "Error FDISK: No se puede crear una partición lógica sin una partición extendida.\n")
		return
	}
	if EspacioUtilizado+int32(size) > TempMBR.MbrSize {
		fmt.Fprintf(buffer, "Error FDISK: No hay suficiente espacio en el disco para crear esta partición.\n")
		return
	}

	// Inicializar el archivo con ceros
	for i := 0; i < size; i++ {
		err := Utilities.WriteObject(file, byte(0), int64(i), buffer)
		if err != nil {
			return
		}
	}


	var vacio int32 = int32(binary.Size(TempMBR))
	if TotalParticiones > 0 {
		vacio = TempMBR.MbrPartitions[TotalParticiones-1].Start + TempMBR.MbrPartitions[TotalParticiones-1].Size
	}

	for i := 0; i < 4; i++ {
		if TempMBR.MbrPartitions[i].Size == 0 {
			if type_ == "p" || type_ == "e" {
				TempMBR.MbrPartitions[i].Size = int32(size)
				TempMBR.MbrPartitions[i].Start = vacio
				copy(TempMBR.MbrPartitions[i].Name[:], name)
				copy(TempMBR.MbrPartitions[i].Fit[:], fit)
				copy(TempMBR.MbrPartitions[i].Status[:], "0")
				copy(TempMBR.MbrPartitions[i].Type[:], type_)
				TempMBR.MbrPartitions[i].Correlative = int32(TotalParticiones + 1)
				if type_ == "e" {
					EBRInicio := vacio
					EBRNuevo := Structs.EBR{
						PartFit:   [1]byte{fit[0]},//revisar
						PartStart: EBRInicio,
						PartSize:  0,
						PartNext:  -1,
					}
					copy(EBRNuevo.PartName[:], "")
					fmt.Println("Creando EBR inicial con Next = -1")
					if err := Utilities.WriteObject(file, EBRNuevo, int64(EBRInicio), buffer); err != nil {
						return
					}
				}
				fmt.Fprintf(buffer, "Partición creada tipo: %s exitosamente en la ruta: %s con el nombre: %s.\n", type_, path, name)
				break
			}
		}
	}

	if type_ == "l" {
		fmt.Println("Entrando a creación de partición lógica")
		var ParticionExtendida *Structs.Partition
		for i := 0; i < 4; i++ {
			fmt.Println("Revisando partición ", i)
			if TempMBR.MbrPartitions[i].Type[0] == 'e' {
				ParticionExtendida = &TempMBR.MbrPartitions[i]
				fmt.Println("Partición extendida encontrada en índice", i)
				break
			}
		}
		if ParticionExtendida == nil {
			fmt.Fprintf(buffer, "Error FDISK: No se encontró una partición extendida para crear la partición lógica.\n")
			return
		}

		EBRPosterior := ParticionExtendida.Start
		fmt.Println("Inicio extendida:", EBRPosterior)
		var EBRUltimo Structs.EBR
		for {
			fmt.Println("Leyendo EBR en posición:", EBRPosterior)
			if err := Utilities.ReadObject(file, &EBRUltimo, int64(EBRPosterior), buffer); err != nil {
				fmt.Println("Error al leer EBR")
				return
			}
			fmt.Println("EBR leido. Nombre:", string(EBRUltimo.PartName[:]), " Next:", EBRUltimo.PartNext)
			if strings.Contains(string(EBRUltimo.PartName[:]), name) {
				fmt.Fprintf(buffer, "Error FDISK: El nombre: %s ya está en uso en las particiones.\n", name)
				return
			}
			if EBRUltimo.PartNext == -1 || (EBRUltimo.PartSize == 0 && EBRUltimo.PartNext == 0) {
				fmt.Println("No hay más EBRs. Se creará el nuevo aquí.")
				break
			}
			EBRPosterior = EBRUltimo.PartNext
		}

		var EBRNuevoPosterior int32
		if EBRUltimo.PartSize == 0 {
			EBRNuevoPosterior = EBRPosterior
			fmt.Println("EBR vacío, se usará la misma posición:", EBRNuevoPosterior)
		} else {
			EBRNuevoPosterior = EBRUltimo.PartStart + EBRUltimo.PartSize
			fmt.Println("Nueva posición para EBR:", EBRNuevoPosterior)
		}

		if EBRNuevoPosterior+int32(size)+int32(binary.Size(Structs.EBR{})) > ParticionExtendida.Start+ParticionExtendida.Size {
			fmt.Println("Error: No hay espacio suficiente en la extendida")
			fmt.Fprintf(buffer, "Error FDISK: No hay suficiente espacio en la partición extendida para esta partición lógica.\n")
			return
		}

		if EBRUltimo.PartSize != 0 {
			EBRUltimo.PartNext = EBRNuevoPosterior
			fmt.Println("Actualizando EBR previo con Next =", EBRNuevoPosterior)
			if err := Utilities.WriteObject(file, EBRUltimo, int64(EBRPosterior), buffer); err != nil {
				fmt.Println("Error al escribir EBR anterior")
				return
			}
		}

		newEBR := Structs.EBR{
			PartFit:   [1]byte{fit[0]},
			PartStart: EBRNuevoPosterior + int32(binary.Size(Structs.EBR{})),
			PartSize:  int32(size),
			PartNext:  -1,
		}
		copy(newEBR.PartName[:], name)
		fmt.Println("Escribiendo nuevo EBR en:", EBRNuevoPosterior)
		if err := Utilities.WriteObject(file, newEBR, int64(EBRNuevoPosterior), buffer); err != nil {
			return
		}
		fmt.Fprintf(buffer, "Partición lógica creada exitosamente en la ruta: %s con el nombre: %s.\n", path, name)
		fmt.Println("---------------------------------------------")
		EBRActual := ParticionExtendida.Start
		for {
			var EBRTemp Structs.EBR
			if err := Utilities.ReadObject(file, &EBRTemp, int64(EBRActual), buffer); err != nil {
				fmt.Fprintf(buffer, "Error leyendo EBR: %v\n", err)
				return
			}
			Structs.PrintEBR(EBRTemp)
			if EBRTemp.PartNext == -1 {
				break
			}
			EBRActual = EBRTemp.PartNext
		}
		fmt.Println("---------------------------------------------")
	}
	//ImprimirEBRsExtendida(file, ParticionExtendida.Start, buffer)

	if err := Utilities.WriteObject(file, TempMBR, 0, buffer); err != nil {
		fmt.Println("Error al escribir nuevo EBR")
		return
	}
	var TempMRB Structs.MRB
	if err := Utilities.ReadObject(file, &TempMRB, 0, buffer); err != nil {
		return
	}
	fmt.Println("---------------------------------------------")
	Structs.PrintMBR(TempMRB)
	fmt.Println("---------------------------------------------")
	defer file.Close()
	fmt.Println("Partición lógica creada con éxito")
}

func Mount(path string, name string, buffer *bytes.Buffer) {
	fmt.Fprintf(buffer, "=========MOUNT=========\n")
	fmt.Print(path)

	if path == "" {
		fmt.Fprintf(buffer, "Error MOUNT: La ruta del disco es obligatoria.\n")
		return
	}
	if name == "" {
		fmt.Fprintf(buffer, "Error MOUNT: El nombre de la partición es obligatorio.\n")
		return
	}

	file, err := Utilities.OpenFile(path, buffer)
	if err != nil {
		return
	}
	defer file.Close()

	var TempMBR Structs.MRB
	if err := Utilities.ReadObject(file, &TempMBR, 0, buffer); err != nil {
		return
	}

	var ParticionExiste = false
	var IndiceParticion = -1
	NameBytes := [16]byte{}
	copy(NameBytes[:], []byte(name))

	for i := 0; i < 4; i++ {
		if TempMBR.MbrPartitions[i].Type[0] == 'e' && bytes.Equal(TempMBR.MbrPartitions[i].Name[:], NameBytes[:]) {
			fmt.Fprintf(buffer, "Error MOUNT: No se puede montar una partición extendida.\n")
			return
		}
	}

	for i := 0; i < 4; i++ {
		if TempMBR.MbrPartitions[i].Type[0] == 'p' && bytes.Equal(TempMBR.MbrPartitions[i].Name[:], NameBytes[:]) {
			ParticionExiste = true
			IndiceParticion = i
			break
		}
	}

	if !ParticionExiste {
		fmt.Fprintf(buffer, "Error MOUNT: No se encontró la partición con el nombre especificado. Solo se pueden montar particiones primarias.\n")
		return
	}

	DiscoID := GeneratorDiscID(path)
	for _, p := range MountedPartitions[DiscoID] {
		if p.Name == name {
			fmt.Fprintf(buffer, "Error MOUNT: La partición ya está montada.\n")
			return
		}
	}

	MountedPartitionsOnDisc := MountedPartitions[DiscoID]
	var Letra byte

	if len(MountedPartitionsOnDisc) == 0 {
		if len(MountedPartitions) == 0 {
			Letra = 'a'
		} else {
			UltimoDiscoID := getLastDiskID()
			UltimaLetra := MountedPartitions[UltimoDiscoID][0].ID[len(MountedPartitions[UltimoDiscoID][0].ID)-1]
			Letra = UltimaLetra + 1
		}
	} else {
		Letra = MountedPartitionsOnDisc[0].ID[len(MountedPartitionsOnDisc[0].ID)-1]
	}

	carnet := "202203009"
	UltimosDigitos := carnet[len(carnet)-2:]
	MountedCount := len(MountedPartitions[DiscoID]) + 1
	IDParticion := fmt.Sprintf("%s%d%c", UltimosDigitos, MountedCount, Letra)

	// Establecer el ID y Status en la partición dentro del MBR (simulado)
	copy(TempMBR.MbrPartitions[IndiceParticion].ID[:], IDParticion)
	TempMBR.MbrPartitions[IndiceParticion].Status[0] = '1'
	if err := Utilities.WriteObject(file, TempMBR, 0, buffer); err != nil {
		fmt.Fprintf(buffer, "Error MOUNT: No se pudo actualizar el MBR: %v\n", err)
		return
	}

	MountedPartitions[DiscoID] = append(MountedPartitions[DiscoID], MountedPartition{
		Path:   path,
		Name:   name,
		ID:     IDParticion,
		Status: '1',
	})

	fmt.Fprintf(buffer, "Partición montada con éxito en la ruta: %s con el nombre: %s y ID: %s.\n", path, name, IDParticion)

	//fmt.Println("---------------------------------------------")
	//PrintMountedPartitions(path, buffer)
	//fmt.Println("---------------------------------------------")

	var TempMRB Structs.MRB
	if err := Utilities.ReadObject(file, &TempMRB, 0, buffer); err != nil {
		return
	}
	Structs.PrintMBR(TempMRB)
	fmt.Println("---------------------------------------------")
}

//***********************Metodos auxiliares***********************

// EliminarDiscoPorRuta Elimina un disco por su ruta
func DeleteDiscWithPath(path string, buffer *bytes.Buffer) {
	discID := GeneratorDiscID(path)
	if _, existe := MountedPartitions[discID]; existe {
		delete(MountedPartitions, discID)
		fmt.Fprintf(buffer, "El disco con ruta '%s' y sus particiones asociadas han sido eliminados.\n", path)
	}
}

//--------------Funciones con discos montados----------------
// GenerarDiscoID Genera un ID único para un disco
func GeneratorDiscID(path string) string {
	return strings.ToLower(path)
}

//función para obtener el último disco montado
func getLastDiskID() string {
	var UltimoDiscoID string
	for DiscoID := range MountedPartitions {
		UltimoDiscoID = DiscoID
	}
	return UltimoDiscoID
}

//-----------------------------------------------------
func MarkPartitionAsLoggedIn(id string) {
	for DiscoID, partitions := range MountedPartitions {
		for i, Particion := range partitions {
			if Particion.ID == id {
				MountedPartitions[DiscoID][i].LoggedIn = true
				return
			}
		}
	}
}

func MarkPartitionAsLoggedOut(id string) {
	for DiscoID, partitions := range MountedPartitions {
		for i, Particion := range partitions {
			if Particion.ID == id {
				MountedPartitions[DiscoID][i].LoggedIn = false
				return
			}
		}
	}
}
