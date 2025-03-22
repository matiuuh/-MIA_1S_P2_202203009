    package Report

import (
	"proyecto1/DiskManagement"
	"proyecto1/Structs"
	"proyecto1/Utilities"
	"bytes"
	"fmt"
	"html"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

func Rep(name string, path string, id string, path_file_ls string, buffer *bytes.Buffer) {
	if name == "" {
		fmt.Fprintf(buffer, "Error REP: El tipo de reporte es obligatorio.\n")
		return
	}
	if path == "" {
		fmt.Fprintf(buffer, "Error REP: La ruta del reporte es obligatoria.\n")
		return
	}
	if id == "" {
		fmt.Fprintf(buffer, "Error REP: El ID de la partición es obligatoria.\n")
		return
	}
	if name == "mbr" {
		ReporteMBR(id, path, buffer)
	} else if name == "disk" {
		ReporteDISK(id, path, buffer)
	} else if name == "sb" {
		ReporteSB(id, path, buffer)
	} else if name == "inode" {
		ReporteInode(id, path, buffer)//arreglar
	} else if name == "bm_inode" {
		Reporte_BitmapInode(id, path, buffer)//arreglar
	} else if name == "bm_bloc" {
		Reporte_BitmapBlock(id, path, buffer)//arreglar
	} else if name == "bloc" {
		ReportBloc(id, path, buffer)//desarrollar
	} else if name == "tree" {
		ReporteTree(id, path, buffer)//desarrollar
	} else if name == "ls" {
		ReporteLS(id, path, buffer)//desarrollar
	} else if name == "file" {
		ReporteFile(id, path, buffer)//desarrollar
	} else {
		fmt.Fprintf(buffer, "Error REP: El tipo de reporte no es válido.\n")
	}
}

func ReporteMBR(id string, path string, buffer *bytes.Buffer) {
	var ParticionesMontadas DiskManagement.MountedPartition
	var ParticionEncontrada bool

	for _, Particiones := range DiskManagement.GetMountedPartitions() {
		for _, Particion := range Particiones {
			if Particion.ID == id {
				ParticionesMontadas = Particion
				ParticionEncontrada = true
				break
			}
		}
		if ParticionEncontrada {
			break
		}
	}

	if !ParticionEncontrada {
		fmt.Fprintf(buffer, "Error REP MBR: No se encontró la partición con el ID: %s.\n", id)
		return
	}

	archivo, err := Utilities.OpenFile(ParticionesMontadas.Path, buffer)
	if err != nil {
		return
	}
	defer archivo.Close()

	var MBRTemporal Structs.MRB
	if err := Utilities.ReadObject(archivo, &MBRTemporal, 0, buffer); err != nil {
		return
	}

	dot := "digraph G {\n"
	dot += "node [shape=plaintext];\n"
	dot += "fontname=\"Courier New\";\n"
	dot += "title [label=\"REPORTE MBR\nMATEO DIEGO\n202203009\"];\n"
	dot += "mbrTable [label=<\n"
	dot += "<table border='1' cellborder='1' cellspacing='0'>\n"
	dot += "<tr><td bgcolor=\"blue\" colspan='2'>MBR</td></tr>\n"
	dot += fmt.Sprintf("<tr><td>Tamaño</td><td>%d</td></tr>\n", MBRTemporal.MbrSize)
	dot += fmt.Sprintf("<tr><td>Fecha De Creación</td><td>%s</td></tr>\n", string(MBRTemporal.CreationDate[:]))
	dot += fmt.Sprintf("<tr><td>Ajuste</td><td>%s</td></tr>\n", string(MBRTemporal.Fit[:]))
	dot += fmt.Sprintf("<tr><td>Signature</td><td>%d</td></tr>\n", MBRTemporal.Signature)
	dot += "</table>\n"
	dot += ">];\n"

	for i, Particion := range MBRTemporal.MbrPartitions {
		if Particion.Size != 0 {
			dot += fmt.Sprintf("PA%d [label=<\n", i+1)
			dot += "<table border='1' cellborder='1' cellspacing='0'>\n"
			dot += fmt.Sprintf("<tr><td bgcolor=\"red\" colspan='2'>Partición %d</td></tr>\n", i+1)
			dot += fmt.Sprintf("<tr><td>Estado</td><td>%s</td></tr>\n", string(Particion.Status[:]))
			dot += fmt.Sprintf("<tr><td>Tipo</td><td>%s</td></tr>\n", string(Particion.Type[:]))
			dot += fmt.Sprintf("<tr><td>Ajuste</td><td>%s</td></tr>\n", string(Particion.Fit[:]))
			dot += fmt.Sprintf("<tr><td>Incio</td><td>%d</td></tr>\n", Particion.Start)
			dot += fmt.Sprintf("<tr><td>Tamaño</td><td>%d</td></tr>\n", Particion.Size)
			dot += fmt.Sprintf("<tr><td>Nombre</td><td>%s</td></tr>\n", strings.Trim(string(Particion.Name[:]), "\x00"))
			dot += fmt.Sprintf("<tr><td>Correlativo</td><td>%d</td></tr>\n", Particion.Correlative)
			dot += "</table>\n"
			dot += ">];\n"
			if Particion.Type[0] == 'e' {
				var EBR Structs.EBR
				if err := Utilities.ReadObject(archivo, &EBR, int64(Particion.Start), buffer); err != nil {
					return
				}
				if EBR.PartSize != 0 {
					var ContadorLogicas int = 0
					dot += "subgraph cluster_0 {style=filled;color=lightgrey;label = \"Partición Extendida\";"
					dot += "fontname=\"Courier New\";"
					for {
						dot += fmt.Sprintf("EBR%d [label=<\n", EBR.PartStart)
						dot += "<table border='1' cellborder='1' cellspacing='0'>\n"
						dot += "<tr><td bgcolor=\"green\" colspan='2'>EBR</td></tr>\n"
						dot += fmt.Sprintf("<tr><td>Nombre</td><td>%s</td></tr>\n", strings.Trim(string(EBR.PartName[:]), "\x00"))
						dot += fmt.Sprintf("<tr><td>Ajuste</td><td>%s</td></tr>\n", string(EBR.PartFit[:]))
						dot += fmt.Sprintf("<tr><td>Inicio</td><td>%d</td></tr>\n", EBR.PartStart)
						dot += fmt.Sprintf("<tr><td>Tamaño</td><td>%d</td></tr>\n", EBR.PartSize)
						dot += fmt.Sprintf("<tr><td>Siguiente</td><td>%d</td></tr>\n", EBR.PartNext)
						dot += "</table>\n"
						dot += ">];\n"

						dot += fmt.Sprintf("Pl%d [label=<\n", ContadorLogicas)
						dot += "<table border='1' cellborder='1' cellspacing='0'>\n"
						dot += "<tr><td bgcolor=\"purple\" colspan='2'>Partición Lógica</td></tr>\n"
						dot += fmt.Sprintf("<tr><td>Estado</td><td>%s</td></tr>\n", string("0"))
						dot += fmt.Sprintf("<tr><td>Tipo</td><td>%s</td></tr>\n", string("l"))
						dot += fmt.Sprintf("<tr><td>Ajuste</td><td>%s</td></tr>\n", string(EBR.PartFit[:]))
						dot += fmt.Sprintf("<tr><td>Incio</td><td>%d</td></tr>\n", EBR.PartStart)
						dot += fmt.Sprintf("<tr><td>Tamaño</td><td>%d</td></tr>\n", EBR.PartSize)
						dot += fmt.Sprintf("<tr><td>Nombre</td><td>%s</td></tr>\n", strings.Trim(string(EBR.PartName[:]), "\x00"))
						dot += fmt.Sprintf("<tr><td>Correlativo</td><td>%d</td></tr>\n", ContadorLogicas+1)
						dot += "</table>\n"
						dot += ">];\n"
						if EBR.PartNext == -1 {
							break
						}
						if err := Utilities.ReadObject(archivo, &EBR, int64(EBR.PartNext), buffer); err != nil {
							fmt.Fprintf(buffer, "Error al leer siguiente EBR: %v\n", err)
							return
						}
						ContadorLogicas++
					}
					dot += "}\n"
				}
			}
		}
	}
	dot += "}\n"
	dotFilePath := "REPORTEMBR.dot"
	err = os.WriteFile(dotFilePath, []byte(dot), 0644)
	if err != nil {
		fmt.Fprintf(buffer, "Error REP MBR: Error al escribir el archivo DOT.\n")
		fmt.Println("Error REP MBR:", err)
		return
	}
	dir := filepath.Dir(path)
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		err = os.MkdirAll(dir, 0755)
		if err != nil {
			fmt.Fprintf(buffer, "Error REP MBR: Error al crear el directorio.\n")
			fmt.Println("Error REP MBR:", err)
			return
		}
	}
	cmd := exec.Command("dot", "-Tjpg", dotFilePath, "-o", path)
	var stderr bytes.Buffer
	cmd.Stderr = &stderr
	err = cmd.Run()
	if err != nil {
		fmt.Fprintf(buffer, "Error REP MBR: Error al ejecutar Graphviz.\n")
		fmt.Println("Error REP MBR:", err)
		return
	}
	fmt.Fprintf(buffer, "Reporte de MBR de la partición:%s generado con éxito en la ruta: %s\n", id, path)
}

func ReporteDISK(id string, path string, buffer *bytes.Buffer) {
	var ParticionesMontadas DiskManagement.MountedPartition
	var ParticionEncontrada bool

	for _, Particiones := range DiskManagement.GetMountedPartitions() {
		for _, Particion := range Particiones {
			if Particion.ID == id {
				ParticionesMontadas = Particion
				ParticionEncontrada = true
				break
			}
		}
		if ParticionEncontrada {
			break
		}
	}

	if !ParticionEncontrada {
		fmt.Fprintf(buffer, "Error REP DISK: No se encontró la partición con el ID: %s.\n", id)
		return
	}

	archivo, err := Utilities.OpenFile(ParticionesMontadas.Path, buffer)
	if err != nil {
		return
	}
	defer archivo.Close()

	var MBRTemporal Structs.MRB
	if err := Utilities.ReadObject(archivo, &MBRTemporal, 0, buffer); err != nil {
		return
	}

	TamanoTotal := float64(MBRTemporal.MbrSize)
	EspacioUsado := 0.0

	dot := "digraph G {\n"
	dot += "labelloc=\"t\"\n"
	dot += "node [shape=plaintext];\n"
	dot += "fontname=\"Courier New\";\n"
	dot += "title [label=\"REPORTE DISK\nMATEO DIEGO\n202203009\"];\n"
	dot += "subgraph cluster1 {\n"
	dot += "fontname=\"Courier New\";\n"
	dot += "label=\"\"\n"
	dot += "disco [shape=none label=<\n"
	dot += "<TABLE border=\"0\" cellspacing=\"4\" cellpadding=\"5\" color=\"skyblue\">\n"
	dot += "<TR><TD bgcolor=\"#a7d0d2\" border=\"1\" cellpadding=\"65\">MBR</TD>\n"

	for i, Particion := range MBRTemporal.MbrPartitions {
		if Particion.Status[0] != 0 {
			partSize := float64(Particion.Size)
			EspacioUsado += partSize

			if Particion.Type[0] == 'e' || Particion.Type[0] == 'E' {
				dot += "<TD border=\"1\" width=\"75\">\n"
				dot += "<TABLE border=\"0\" cellspacing=\"4\" cellpadding=\"10\">\n"
				dot += "<TR><TD bgcolor=\"skyblue\" border=\"1\" colspan=\"5\" height=\"75\"> Partición Extendida<br/></TD></TR>\n"

				EspacioLibreExtendida := partSize
				finEbr := Particion.Start

				var EBR Structs.EBR
				if err := Utilities.ReadObject(archivo, &EBR, int64(Particion.Start), buffer); err != nil {
					return
				}
				if EBR.PartSize != 0 {
					for {
						var ebr Structs.EBR
						if err := Utilities.ReadObject(archivo, &ebr, int64(finEbr), buffer); err != nil {
							fmt.Println("Error al leer EBR:", err)
							fmt.Fprintf(buffer, "Error en linea : Error al leer EBR")
							break
						}

						TamanoEBR := float64(ebr.PartSize)
						EspacioUsado += TamanoEBR
						EspacioLibreExtendida -= TamanoEBR

						dot += "<TR>\n"
						dot += "<TD bgcolor=\"#264b5e\" border=\"1\" height=\"185\">EBR</TD>\n"
						dot += fmt.Sprintf("<TD bgcolor=\"#546eab\" border=\"1\" cellpadding=\"25\">Partición Lógica<br/>%.2f%% del Disco</TD>\n", (TamanoEBR/TamanoTotal)*100)
						dot += "</TR>\n"
						if ebr.PartNext <= 0 {
							break
						}
						finEbr = ebr.PartNext
					}
				}
				dot += "<TR>\n"
				dot += fmt.Sprintf("<TD bgcolor=\"#f1e6d2\" border=\"1\" colspan=\"5\"> Espacio Libre Dentro De La Partición Extendida<br/>%.2f%% del Disco</TD>\n", (EspacioLibreExtendida/TamanoTotal)*100)
				dot += "</TR>\n"

				dot += "</TABLE>\n</TD>\n"
			} else if Particion.Type[0] == 'p' || Particion.Type[0] == 'P' {
				dot += fmt.Sprintf("<TD bgcolor=\"#4697b4\" border=\"1\" cellpadding=\"20\">Partición Primaria %d<br/>%.2f%% del Disco</TD>\n", i+1, (partSize/TamanoTotal)*100)
			}
		}
	}

	Porcentaje := 100.0
	for _, partition := range MBRTemporal.MbrPartitions {
		if partition.Status[0] != 0 {
			partSize := float64(partition.Size)
			Porcentaje -= (partSize / TamanoTotal) * 100
		}
	}

	dot += fmt.Sprintf("<TD bgcolor=\"#f1e6d2\" border=\"1\" cellpadding=\"20\">Espacio Libre<br/>%.2f%% del Disco</TD>\n", Porcentaje)
	dot += "</TR>\n"
	dot += "</TABLE>\n"
	dot += ">];\n"
	dot += "}\n"
	dot += "}\n"

	RutaReporte := "REPORTEDISK.dot"
	err = os.WriteFile(RutaReporte, []byte(dot), 0644)
	if err != nil {
		fmt.Fprintf(buffer, "Error REP DISK: Error al escribir el archivo DOT.\n")
		fmt.Println("Error REP DISK:", err)
		return
	}
	dir := filepath.Dir(path)
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		err = os.MkdirAll(dir, 0755)
		if err != nil {
			fmt.Fprintf(buffer, "Error REP DISK: Error al crear el directorio.\n")
			fmt.Println("Error REP DISK:", err)
			return
		}
	}
	cmd := exec.Command("dot", "-Tjpg", RutaReporte, "-o", path)
	var stderr bytes.Buffer
	cmd.Stderr = &stderr
	err = cmd.Run()
	if err != nil {
		fmt.Fprintf(buffer, "Error REP DISK: Error al ejecutar Graphviz.")
		fmt.Println("Error REP DISK:", err)
		return
	}
	fmt.Fprintf(buffer, "Reporte de DISK de la partición:%s generado con éxito en la ruta: %s\n", id, path)
}

func ReporteSB(id string, path string, buffer *bytes.Buffer) {
	var ParticionesMontadas DiskManagement.MountedPartition
	var ParticionEncontrada bool

	for _, Particiones := range DiskManagement.GetMountedPartitions() {
		for _, Particion := range Particiones {
			if Particion.ID == id {
				ParticionesMontadas = Particion
				ParticionEncontrada = true
				break
			}
		}
		if ParticionEncontrada {
			break
		}
	}

	if !ParticionEncontrada {
		fmt.Fprintf(buffer, "Error REP SB: No se encontró la partición con el ID: %s.\n", id)
		return
	}

	archivo, err := Utilities.OpenFile(ParticionesMontadas.Path, buffer)
	if err != nil {
		return
	}
	defer archivo.Close()

	var MBRTemporal Structs.MRB
	if err := Utilities.ReadObject(archivo, &MBRTemporal, 0, buffer); err != nil {
		return
	}

	var index int = -1
	for i := 0; i < 4; i++ {
		if MBRTemporal.MbrPartitions[i].Size != 0 {
			if strings.Contains(string(MBRTemporal.MbrPartitions[i].ID[:]), id) {
				if MBRTemporal.MbrPartitions[i].Status[0] == '1' {
					index = i
				} else {
					fmt.Fprintf(buffer, "Error REP SB: La partición con el ID:%s no está montada.\n", id)
					return
				}
				break
			}
		}
	}

	if index == -1 {
		fmt.Fprintf(buffer, "Error REP SB: No se encontró la partición con el ID: %s.\n", id)
		return
	}

	var TemporalSuperBloque = Structs.Superblock{}
	if err := Utilities.ReadObject(archivo, &TemporalSuperBloque, int64(MBRTemporal.MbrPartitions[index].Start), buffer); err != nil {
		fmt.Fprintf(buffer, "Error REP SB: Error al leer el SuperBloque.\n")
		return
	}

	dot := "digraph G {\n"
	dot += "node [shape=plaintext];\n"
	dot += "fontname=\"Courier New\";\n"
	dot += "title [label=\"REPORTE SB\nMATEO DIEGO\n202203009\"];\n"
	dot += "SBTable [label=<\n"
	dot += "<table border='1' cellborder='1' cellspacing='0'>\n"
	dot += "<tr><td bgcolor=\"skyblue\" colspan='2'>Super Bloque</td></tr>\n"
	dot += fmt.Sprintf("<tr><td>SB FileSystem Type</td><td>%d</td></tr>\n", int(TemporalSuperBloque.SB_FileSystem_Type))
	dot += fmt.Sprintf("<tr><td>SB Inodes Count</td><td>%d</td></tr>\n", int(TemporalSuperBloque.SB_Inodes_Count))
	dot += fmt.Sprintf("<tr><td>SB Blocks Count</td><td>%d</td></tr>\n", int(TemporalSuperBloque.SB_Blocks_Count))
	dot += fmt.Sprintf("<tr><td>SB Free Blocks Count</td><td>%d</td></tr>\n", int(TemporalSuperBloque.SB_Free_Blocks_Count))
	dot += fmt.Sprintf("<tr><td>SB Free Inodes Count</td><td>%d</td></tr>\n", int(TemporalSuperBloque.SB_Free_Inodes_Count))
	dot += fmt.Sprintf("<tr><td>SB Mtime</td><td>%s</td></tr>\n", TemporalSuperBloque.SB_Mtime[:])
	dot += fmt.Sprintf("<tr><td>SB Umtime</td><td>%s</td></tr>\n", TemporalSuperBloque.SB_Umtime[:])
	dot += fmt.Sprintf("<tr><td>SB Mnt Count</td><td>%d</td></tr>\n", int(TemporalSuperBloque.SB_Mnt_Count))
	dot += fmt.Sprintf("<tr><td>SB Magic</td><td>%d</td></tr>\n", int(TemporalSuperBloque.SB_Magic))
	dot += fmt.Sprintf("<tr><td>SB Inode Size</td><td>%d</td></tr>\n", int(TemporalSuperBloque.SB_Inode_Size))
	dot += fmt.Sprintf("<tr><td>SB Block Size</td><td>%d</td></tr>\n", int(TemporalSuperBloque.SB_Block_Size))
	dot += fmt.Sprintf("<tr><td>SB Fist Inode</td><td>%d</td></tr>\n", int(TemporalSuperBloque.SB_Fist_Ino))
	dot += fmt.Sprintf("<tr><td>SB First Block</td><td>%d</td></tr>\n", int(TemporalSuperBloque.SB_First_Blo))
	dot += fmt.Sprintf("<tr><td>SB Bm Inode Start</td><td>%d</td></tr>\n", int(TemporalSuperBloque.SB_Bm_Inode_Start))
	dot += fmt.Sprintf("<tr><td>SB Bm Block Start</td><td>%d</td></tr>\n", int(TemporalSuperBloque.SB_Bm_Block_Start))
	dot += fmt.Sprintf("<tr><td>SB Inode Start</td><td>%d</td></tr>\n", int(TemporalSuperBloque.SB_Inode_Start))
	dot += fmt.Sprintf("<tr><td>SB Block Start</td><td>%d</td></tr>\n", int(TemporalSuperBloque.SB_Block_Start))
	dot += "</table>\n"
	dot += ">];\n"
	dot += "}\n"

	RutaReporte := "REPORTESB.dot"
	err = os.WriteFile(RutaReporte, []byte(dot), 0644)
	if err != nil {
		fmt.Fprintf(buffer, "Error REP SB: Error al escribir el archivo DOT.\n")
		fmt.Println("Error REP DISK:", err)
		return
	}
	dir := filepath.Dir(path)
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		err = os.MkdirAll(dir, 0755)
		if err != nil {
			fmt.Fprintf(buffer, "Error REP SB: Error al crear el directorio.\n")
			fmt.Println("Error REP DISK:", err)
			return
		}
	}
	cmd := exec.Command("dot", "-Tjpg", RutaReporte, "-o", path)
	var stderr bytes.Buffer
	cmd.Stderr = &stderr
	err = cmd.Run()
	if err != nil {
		fmt.Fprintf(buffer, "Error REP SB: Error al ejecutar Graphviz.")
		fmt.Println("Error REP DISK:", err)
		return
	}
	fmt.Fprintf(buffer, "Reporte de SB de la partición:%s generado con éxito en la ruta: %s\n", id, path)
}

func ReporteInode(id string, path string, buffer *bytes.Buffer) {
	var ParticionesMontadas DiskManagement.MountedPartition
	var ParticionEncontrada bool

	for _, Particiones := range DiskManagement.GetMountedPartitions() {
		for _, Particion := range Particiones {
			if Particion.ID == id {
				ParticionesMontadas = Particion
				ParticionEncontrada = true
				break
			}
		}
		if ParticionEncontrada {
			break
		}
	}

	if !ParticionEncontrada {
		fmt.Fprintf(buffer, "Error REP SB: No se encontró la partición con el ID: %s.\n", id)
		return
	}

	archivo, err := Utilities.OpenFile(ParticionesMontadas.Path, buffer)
	if err != nil {
		return
	}
	defer archivo.Close()

	var MBRTemporal Structs.MRB
	if err := Utilities.ReadObject(archivo, &MBRTemporal, 0, buffer); err != nil {
		return
	}

	var index int = -1
	for i := 0; i < 4; i++ {
		if MBRTemporal.MbrPartitions[i].Size != 0 {
			if strings.Contains(string(MBRTemporal.MbrPartitions[i].ID[:]), id) {
				if MBRTemporal.MbrPartitions[i].Status[0] == '1' {
					index = i
				} else {
					fmt.Fprintf(buffer, "Error REP Inode: La partición con el ID:%s no está montada.\n", id)
					return
				}
				break
			}
		}
	}

	if index == -1 {
		fmt.Fprintf(buffer, "Error REP Inode: No se encontró la partición con el ID: %s.\n", id)
		return
	}

	var TemporalSuperBloque = Structs.Superblock{}
	if err := Utilities.ReadObject(archivo, &TemporalSuperBloque, int64(MBRTemporal.MbrPartitions[index].Start), buffer); err != nil {
		fmt.Fprintf(buffer, "Error REP Inode: Error al leer el SuperBloque.\n")
		return
	}

	var dot bytes.Buffer

	fmt.Fprintln(&dot, "digraph G {")
	fmt.Fprintln(&dot, "node [shape=none];")
	fmt.Fprintln(&dot, "fontname=\"Courier New\";")
	fmt.Fprintln(&dot, "title [label=\"REPORTE INODE\nMATEO DIEGO\n202203009\"];")

	for i := 0; i < int(TemporalSuperBloque.SB_Inodes_Count); i++ {
		var inode Structs.Inode

		if err := Utilities.ReadObject(archivo, &inode, int64(TemporalSuperBloque.SB_Inode_Start)+int64(i)*int64(TemporalSuperBloque.SB_Inode_Size), buffer); err != nil {
			fmt.Println("Error al leer el inodo:", err)
			continue
		}

		if inode.IN_Size > 0 {
			fmt.Fprintf(&dot, "inode%d [label=<\n", i)
			fmt.Fprintf(&dot, "<table border='0' cellborder='1' cellspacing='0' cellpadding='10'>\n")
			fmt.Fprintf(&dot, "<tr><td colspan='2' bgcolor='skyblue'>Inode %d</td></tr>\n", i)
			fmt.Fprintf(&dot, "<tr><td>UID</td><td>%d</td></tr>\n", inode.IN_Uid)
			fmt.Fprintf(&dot, "<tr><td>GID</td><td>%d</td></tr>\n", inode.IN_Gid)
			fmt.Fprintf(&dot, "<tr><td>Size</td><td>%d</td></tr>\n", inode.IN_Size)
			fmt.Fprintf(&dot, "<tr><td>ATime</td><td>%s</td></tr>\n", html.EscapeString(string(inode.IN_Atime[:])))
			fmt.Fprintf(&dot, "<tr><td>CTime</td><td>%s</td></tr>\n", html.EscapeString(string(inode.IN_Ctime[:])))
			fmt.Fprintf(&dot, "<tr><td>MTime</td><td>%s</td></tr>\n", html.EscapeString(string(inode.IN_Mtime[:])))
			fmt.Fprintf(&dot, "<tr><td>Blocks</td><td>%v</td></tr>\n", inode.IN_Block)
			//fmt.Fprintf(&dot, "<tr><td>Type</td><td>%c</td></tr>\n", inode.IN_Type[0])
			fmt.Fprintf(&dot, "<tr><td>Perms</td><td>%v</td></tr>\n", inode.IN_Perm)
			fmt.Fprintf(&dot, "</table>\n")
			fmt.Fprintf(&dot, " >];\n")
		}
	}
	fmt.Fprintln(&dot, "}")

	RutaReporte := "REPORTEINODE.dot"
	err = os.WriteFile(RutaReporte, dot.Bytes(), 0644)
	if err != nil {
		fmt.Fprintf(buffer, "Error REP INODE: Error al escribir el archivo DOT.\n")
		fmt.Println("Error REP DISK:", err)
		return
	}
	dir := filepath.Dir(path)
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		err = os.MkdirAll(dir, 0755)
		if err != nil {
			fmt.Fprintf(buffer, "Error REP INODE: Error al crear el directorio.\n")
			fmt.Println("Error REP INODE:", err)
			return
		}
	}
	cmd := exec.Command("dot", "-Tjpg", RutaReporte, "-o", path)
	var stderr bytes.Buffer
	cmd.Stderr = &stderr
	err = cmd.Run()
	if err != nil {
		fmt.Fprintf(buffer, "Error REP INODE: Error al ejecutar Graphviz.")
		fmt.Println("Error REP INODE:", err)
		return
	}
	fmt.Fprintf(buffer, "Reporte de INODE de la partición:%s generado con éxito en la ruta: %s\n", id, path)
}

func Reporte_BitmapInode(id string, path string, buffer *bytes.Buffer) {
	dir := filepath.Dir(path)
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		err = os.MkdirAll(dir, 0755)
		if err != nil {
			fmt.Fprintf(buffer, "Error REP BITMAP INODE: Error al crear el directorio: %v", err)
			return
		}
	}
	var ParticionesMontadas DiskManagement.MountedPartition
	var ParticionEncontrada bool

	for _, Particiones := range DiskManagement.GetMountedPartitions() {
		for _, Particion := range Particiones {
			if Particion.ID == id {
				ParticionesMontadas = Particion
				ParticionEncontrada = true
				break
			}
		}
		if ParticionEncontrada {
			break
		}
	}

	if !ParticionEncontrada {
		fmt.Fprintf(buffer, "Error REP BITMAP INODE: No se encontró la partición con el ID: %s.\n", id)
		return
	}

	archivo, err := Utilities.OpenFile(ParticionesMontadas.Path, buffer)
	if err != nil {
		return
	}
	defer archivo.Close()

	var MBRTemporal Structs.MRB
	if err := Utilities.ReadObject(archivo, &MBRTemporal, 0, buffer); err != nil {
		return
	}

	var index int = -1
	for i := 0; i < 4; i++ {
		if MBRTemporal.MbrPartitions[i].Size != 0 {
			if strings.Contains(string(MBRTemporal.MbrPartitions[i].ID[:]), id) {
				if MBRTemporal.MbrPartitions[i].Status[0] == '1' {
					index = i
				} else {
					fmt.Fprintf(buffer, "Error REP BITMAP INODE: La partición con el ID:%s no está montada.\n", id)
					return
				}
				break
			}
		}
	}

	if index == -1 {
		fmt.Fprintf(buffer, "Error REP BITMAP INODE: No se encontró la partición con el ID: %s.\n", id)
		return
	}

	var TemporalSuperBloque = Structs.Superblock{}
	if err := Utilities.ReadObject(archivo, &TemporalSuperBloque, int64(MBRTemporal.MbrPartitions[index].Start), buffer); err != nil {
		fmt.Fprintf(buffer, "Error REP BITMAP INODE: Error al leer el SuperBloque.\n")
		return
	}

	BitMapInode := make([]byte, TemporalSuperBloque.SB_Inodes_Count)
	if _, err := archivo.ReadAt(BitMapInode, int64(TemporalSuperBloque.SB_Bm_Inode_Start)); err != nil {
		fmt.Fprint(buffer, "Error REP BITMAP INODE: No se pudo leer el bitmap de inodos:", err)
		return
	}

	SalidaArchivo, err := os.Create(path)
	if err != nil {
		fmt.Fprint(buffer, "Error REP BITMAP INODE: No se pudo crear el archivo de reporte:", err)
		return
	}
	defer SalidaArchivo.Close()

	fmt.Fprintln(SalidaArchivo, "REPORTE BITMAP INODE\nMATEO DIEGO\n202203009")
	fmt.Fprintln(SalidaArchivo, "---------------------------------------")

	for i, bit := range BitMapInode {
		if i > 0 && i%20 == 0 {
			fmt.Fprintln(SalidaArchivo)
		}
		fmt.Fprintf(SalidaArchivo, "%d ", bit)
	}

	fmt.Fprintf(buffer, "Reporte de BITMAP INODE de la partición:%s generado con éxito en la ruta: %s\n", id, path)
}

func Reporte_BitmapBlock(id string, path string, buffer *bytes.Buffer) {
	dir := filepath.Dir(path)
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		err = os.MkdirAll(dir, 0755)
		if err != nil {
			fmt.Fprintf(buffer, "Error REP BITMAP BLOCK: Error al crear el directorio: %v", err)
			return
		}
	}

	var ParticionesMontadas DiskManagement.MountedPartition
	var ParticionEncontrada bool

	for _, Particiones := range DiskManagement.GetMountedPartitions() {
		for _, Particion := range Particiones {
			if Particion.ID == id {
				ParticionesMontadas = Particion
				ParticionEncontrada = true
				break
			}
		}
		if ParticionEncontrada {
			break
		}
	}

	if !ParticionEncontrada {
		fmt.Fprintf(buffer, "Error REP BITMAP BLOCK: No se encontró la partición con el ID: %s.\n", id)
		return
	}

	archivo, err := Utilities.OpenFile(ParticionesMontadas.Path, buffer)
	if err != nil {
		return
	}
	defer archivo.Close()

	var MBRTemporal Structs.MRB
	if err := Utilities.ReadObject(archivo, &MBRTemporal, 0, buffer); err != nil {
		return
	}

	var index int = -1
	for i := 0; i < 4; i++ {
		if MBRTemporal.MbrPartitions[i].Size != 0 {
			if strings.Contains(string(MBRTemporal.MbrPartitions[i].ID[:]), id) {
				if MBRTemporal.MbrPartitions[i].Status[0] == '1' {
					index = i
				} else {
					fmt.Fprintf(buffer, "Error REP BITMAP BLOCK: La partición con el ID:%s no está montada.\n", id)
					return
				}
				break
			}
		}
	}

	if index == -1 {
		fmt.Fprintf(buffer, "Error REP BITMAP BLOCK: No se encontró la partición con el ID: %s.\n", id)
		return
	}

	var TemporalSuperBloque = Structs.Superblock{}
	if err := Utilities.ReadObject(archivo, &TemporalSuperBloque, int64(MBRTemporal.MbrPartitions[index].Start), buffer); err != nil {
		fmt.Fprintf(buffer, "Error REP BITMAP BLOCK: Error al leer el SuperBloque.\n")
		return
	}

	BitMaBlock := make([]byte, TemporalSuperBloque.SB_Blocks_Count)
	if _, err := archivo.ReadAt(BitMaBlock, int64(TemporalSuperBloque.SB_Bm_Block_Start)); err != nil {
		fmt.Fprint(buffer, "Error REP BITMAP BLOCK: No se pudo leer el bitmap de bloque:", err)
		return
	}

	SalidaArchivo, err := os.Create(path)
	if err != nil {
		fmt.Fprint(buffer, "Error REP BITMAP BLOCK: No se pudo crear el archivo de reporte:", err)
		return
	}
	defer SalidaArchivo.Close()

	fmt.Fprintln(SalidaArchivo, "REPORTE BITMAP BLOCK\nMATEO DIEGO\n202203009")
	fmt.Fprintln(SalidaArchivo, "---------------------------------------")

	for i, bit := range BitMaBlock {
		if i > 0 && i%20 == 0 {
			fmt.Fprintln(SalidaArchivo)
		}
		fmt.Fprintf(SalidaArchivo, "%d ", bit)
	}

	fmt.Fprintf(buffer, "Reporte de BITMAP BLOCK de la partición:%s generado con éxito en la ruta: %s\n", id, path)
}

func ReportBloc(id string, path string, buffer *bytes.Buffer){
	fmt.Println("===========Reporte de Bloques===========\n")
}

func ReporteTree(id string, path string, buffer *bytes.Buffer){
	fmt.Println("===========Reporte de arbol===========\n")
}

func ReporteLS(id string, path string, buffer *bytes.Buffer){
	fmt.Println("===========Reporte ls===========\n")
}

func ReporteFile(id string, path string, buffer *bytes.Buffer){
	fmt.Println("===========Reporte File===========\n")
}
