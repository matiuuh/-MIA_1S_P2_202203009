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
	"encoding/binary"
	"proyecto1/FileSystem"
	"unicode"
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
		ReporteInode(id, path, buffer)
	} else if name == "bm_inode" {
		Reporte_BitmapInode(id, path, buffer)
	} else if name == "bm_bloc" {
		Reporte_BitmapBlock(id, path, buffer)
	} else if name == "block" {
		ReportBloc(id, path, buffer)
	} else if name == "tree" {
		ReporteTree(id, path, buffer)//desarrollar
	} else if name == "ls" {
		if path_file_ls == "" {
			fmt.Fprintf(buffer, "Error REP LS: Debes especificar el parámetro -path_file_ls\n")
			return
		}
		ReporteLS(id, path, buffer, path_file_ls)//desarrollar
	} else if name == "file" {
		if path_file_ls == "" {
			fmt.Fprintf(buffer, "Error REP FILE: Debes especificar el parámetro -path_file_ls\n")
			return
		}
		ReporteFile(id, path, buffer, path_file_ls)
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
		fmt.Fprintf(buffer, "Error REP INODE: No se encontró la partición con el ID: %s.\n", id)
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
					fmt.Fprintf(buffer, "Error REP INODE: La partición con el ID:%s no está montada.\n", id)
					return
				}
				break
			}
		}
	}

	if index == -1 {
		fmt.Fprintf(buffer, "Error REP INODE: No se encontró la partición con el ID: %s.\n", id)
		return
	}

	var sb Structs.Superblock
	if err := Utilities.ReadObject(archivo, &sb, int64(MBRTemporal.MbrPartitions[index].Start), buffer); err != nil {
		fmt.Fprintf(buffer, "Error REP INODE: Error al leer el SuperBloque.\n")
		return
	}

	var dot bytes.Buffer
	fmt.Fprintln(&dot, "digraph G {")
	fmt.Fprintln(&dot, "node [shape=none];")
	fmt.Fprintln(&dot, "fontname=\"Courier New\";")

	for i := 0; i < int(sb.SB_Inodes_Count); i++ {
		var inode Structs.Inode
		offset := int64(sb.SB_Inode_Start) + int64(i)*int64(binary.Size(Structs.Inode{}))
		if err := Utilities.ReadObject(archivo, &inode, offset, buffer); err != nil {
			continue
		}

		if inode.IN_Size > 0 {
			fmt.Fprintf(&dot, "inode%d [label=<\n", i)
			fmt.Fprintln(&dot, "<table border='1' cellborder='1' cellspacing='0'>")
			fmt.Fprintf(&dot, "<tr><td colspan='2' bgcolor='skyblue'>Inodo %d</td></tr>\n", i)
			fmt.Fprintf(&dot, "<tr><td>i_uid</td><td>%d</td></tr>\n", inode.IN_Uid)
			fmt.Fprintf(&dot, "<tr><td>i_gid</td><td>%d</td></tr>\n", inode.IN_Gid)
			fmt.Fprintf(&dot, "<tr><td>i_size</td><td>%d</td></tr>\n", inode.IN_Size)
		
			atime := html.EscapeString(strings.TrimSpace(string(bytes.Trim(inode.IN_Atime[:], "\x00"))))
			ctime := html.EscapeString(strings.TrimSpace(string(bytes.Trim(inode.IN_Ctime[:], "\x00"))))
			mtime := html.EscapeString(strings.TrimSpace(string(bytes.Trim(inode.IN_Mtime[:], "\x00"))))
		
			fmt.Fprintf(&dot, "<tr><td>i_atime</td><td>%s</td></tr>\n", atime)
			fmt.Fprintf(&dot, "<tr><td>i_ctime</td><td>%s</td></tr>\n", ctime)
			fmt.Fprintf(&dot, "<tr><td>i_mtime</td><td>%s</td></tr>\n", mtime)
		
			for j, blk := range inode.IN_Block {
				if blk != -1 {
					fmt.Fprintf(&dot, "<tr><td>i_block_%d</td><td>%d</td></tr>\n", j, blk)
				}
			}
		
			perm := html.EscapeString(strings.TrimSpace(string(bytes.Trim(inode.IN_Perm[:], "\x00"))))
			fmt.Fprintf(&dot, "<tr><td>i_perm</td><td>%s</td></tr>\n", perm)
		
			fmt.Fprintln(&dot, "</table>>];")
		}
	}

	fmt.Fprintln(&dot, "}")

	// Guardar archivo DOT
	dotFile := strings.ReplaceAll(path, ".jpg", ".dot")
	err = os.WriteFile(dotFile, dot.Bytes(), 0644)
	if err != nil {
		fmt.Fprintf(buffer, "Error REP INODE: No se pudo escribir el archivo DOT.\n")
		return
	}

	// Asegurar directorio
	dir := filepath.Dir(path)
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		_ = os.MkdirAll(dir, 0755)
	}

	// Ejecutar Graphviz
	cmd := exec.Command("dot", "-Tjpg", dotFile, "-o", path)
	var stderr bytes.Buffer
	cmd.Stderr = &stderr
	err = cmd.Run()
	if err != nil {
		fmt.Fprintf(buffer, "Error REP INODE: Graphviz falló: %s\n", stderr.String())
		return
	}

	fmt.Fprintf(buffer, "Reporte de INODE generado correctamente en %s\n", path)
}


func Reporte_BitmapInode(id string, path string, buffer *bytes.Buffer) {
	path = corregirExtensionTxt(path) // 🔧 Forzar extensión .txt
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
	path = corregirExtensionTxt(path) // 🔧 Forzar extensión .txt
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

func ReportBloc(id string, path string, buffer *bytes.Buffer) {
	fmt.Fprintln(buffer, "=========== Reporte de Bloques ===========")

	// 1. Obtener partición montada
	particiones := DiskManagement.GetMountedPartitions()
	var diskPath string
	var encontrado bool
	for _, lista := range particiones {
		for _, part := range lista {
			if part.ID == id && part.LoggedIn {
				diskPath = part.Path
				encontrado = true
				break
			}
		}
		if encontrado {
			break
		}
	}

	if !encontrado {
		fmt.Fprintln(buffer, "Error: No se encontró partición montada con ID", id)
		return
	}

	// 2. Leer Superblock
	file, err := Utilities.OpenFile(diskPath, buffer)
	if err != nil {
		fmt.Fprintln(buffer, "Error al abrir disco:", err)
		return
	}
	defer file.Close()

	var mbr Structs.MRB
	if err := Utilities.ReadObject(file, &mbr, 0, buffer); err != nil {
		fmt.Fprintln(buffer, "Error al leer MBR:", err)
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
		fmt.Fprintln(buffer, "Error: No se encontró la partición activa en el MBR")
		return
	}

	var sb Structs.Superblock
	if err := Utilities.ReadObject(file, &sb, int64(mbr.MbrPartitions[index].Start), buffer); err != nil {
		fmt.Fprintln(buffer, "Error al leer el Superblock:", err)
		return
	}

	// 3. Iniciar DOT
	var dot bytes.Buffer
	dot.WriteString("digraph G {\n")
	dot.WriteString("rankdir=LR;\n")
	dot.WriteString("node [shape=record, fontname=Helvetica];\n")

	for i := int32(0); i < sb.SB_Blocks_Count; i++ {
		var b byte
		if err := Utilities.ReadObject(file, &b, int64(sb.SB_Bm_Block_Start+i), buffer); err != nil {
			continue
		}
		if b == 0 {
			continue
		}

		// Leer bloque
		blockOffset := int64(sb.SB_Block_Start + i*int32(binary.Size(Structs.FolderBlock{})))
		var folder Structs.FolderBlock
		if err := Utilities.ReadObject(file, &folder, blockOffset, buffer); err == nil {
			// FolderBlock (bloque de carpeta)
			dot.WriteString(fmt.Sprintf("block%d [label=\"Bloque Carpeta %d|", i, i))
			for _, entry := range folder.B_Content {
				name := strings.Trim(string(entry.B_Name[:]), "\x00")
				if name != "" && entry.B_Inode != -1 {
					dot.WriteString(fmt.Sprintf("%s | %d\\l", name, entry.B_Inode))
				}
			}
			dot.WriteString("\"];\n")
			continue
		}

		var fileblock Structs.FileBlock
		blockOffset = int64(sb.SB_Block_Start + i*int32(binary.Size(Structs.FileBlock{})))
		if err := Utilities.ReadObject(file, &fileblock, blockOffset, buffer); err == nil {
			// FileBlock (bloque de archivo)
			content := strings.ReplaceAll(string(fileblock.B_Content[:]), "\"", "'")
			content = strings.ReplaceAll(content, "\n", " ")
			dot.WriteString(fmt.Sprintf("block%d [label=\"Bloque Archivo %d\\n%s\"];\n", i, i, content))
			continue
		}

		// Si implementas apuntadores indirectos, también puedes detectar tipo Structs.PointerBlock
	}

	dot.WriteString("}\n")

	// 4. Guardar archivo .dot
	outputDot := strings.ReplaceAll(path, ".jpg", ".dot")
	err = os.WriteFile(outputDot, dot.Bytes(), 0644)
	if err != nil {
		fmt.Fprintf(buffer, "Error al guardar archivo DOT: %v\n", err)
		return
	}

	// 5. Generar imagen con Graphviz
	cmd := exec.Command("dot", "-Tjpg", outputDot, "-o", path)
	var stderr bytes.Buffer
	cmd.Stderr = &stderr
	err = cmd.Run()
	if err != nil {
		fmt.Fprintf(buffer, "Error al ejecutar Graphviz: %v\n", err)
		fmt.Fprintf(buffer, "Detalles: %s\n", stderr.String())
		return
	}

	fmt.Fprintf(buffer, "Reporte de bloques generado con éxito en la ruta: %s\n", path)
}


func ReporteTree(id string, path string, buffer *bytes.Buffer) {
	fmt.Fprintln(buffer, "=========== Reporte Tree ===========")

	// 1. Obtener partición montada
	particiones := DiskManagement.GetMountedPartitions()
	var diskPath string
	var encontrado bool
	for _, lista := range particiones {
		for _, part := range lista {
			if part.ID == id && part.LoggedIn {
				diskPath = part.Path
				encontrado = true
				break
			}
		}
		if encontrado {
			break
		}
	}

	if !encontrado {
		fmt.Fprintln(buffer, "Error: No se encontró partición montada con ID", id)
		return
	}

	// 2. Abrir archivo y leer MBR y Superblock
	file, err := Utilities.OpenFile(diskPath, buffer)
	if err != nil {
		fmt.Fprintln(buffer, "Error al abrir disco:", err)
		return
	}
	defer file.Close()

	var mbr Structs.MRB
	if err := Utilities.ReadObject(file, &mbr, 0, buffer); err != nil {
		fmt.Fprintln(buffer, "Error al leer MBR:", err)
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
		fmt.Fprintln(buffer, "Error: No se encontró la partición activa en el MBR")
		return
	}

	var sb Structs.Superblock
	if err := Utilities.ReadObject(file, &sb, int64(mbr.MbrPartitions[index].Start), buffer); err != nil {
		fmt.Fprintln(buffer, "Error al leer el Superblock:", err)
		return
	}

	// 3. Generar contenido del DOT
	var dot bytes.Buffer
	dot.WriteString("digraph G {\n")
	dot.WriteString("rankdir=LR;\n")
	dot.WriteString("node [shape=record, fontname=Helvetica];\n")

	// Recorremos el sistema desde el inodo raíz
	generarArbolInodos(0, &sb, file, &dot, buffer)

	dot.WriteString("}\n")

	// 4. Guardar archivo .dot
	outputDotPath := strings.ReplaceAll(path, ".jpg", ".dot")
	err = os.WriteFile(outputDotPath, dot.Bytes(), 0644)
	if err != nil {
		fmt.Fprintf(buffer, "Error al escribir archivo DOT: %v\n", err)
		return
	}

	// Crear directorio de destino si no existe
	dir := filepath.Dir(path)
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		err = os.MkdirAll(dir, 0755)
		if err != nil {
			fmt.Fprintf(buffer, "Error al crear directorio destino: %v\n", err)
			return
		}
	}

	// Ejecutar Graphviz para generar la imagen
	cmd := exec.Command("dot", "-Tjpg", outputDotPath, "-o", path)
	var stderr bytes.Buffer
	cmd.Stderr = &stderr
	err = cmd.Run()
	if err != nil {
		fmt.Fprintf(buffer, "Error al ejecutar Graphviz: %v\n", err)
		fmt.Fprintf(buffer, "Detalles: %s\n", stderr.String())
		return
	}

	fmt.Fprintf(buffer, "Reporte Tree generado con éxito en la ruta: %s\n", path)
}

func generarArbolInodos(index int32, sb *Structs.Superblock, file *os.File, dot *bytes.Buffer, buffer *bytes.Buffer) {
	inodeOffset := int64(sb.SB_Inode_Start + index*int32(binary.Size(Structs.Inode{})))
	var inode Structs.Inode
	if err := Utilities.ReadObject(file, &inode, inodeOffset, buffer); err != nil {
		fmt.Fprintf(buffer, "Error al leer inodo %d: %v\n", index, err)
		return
	}

	// Construir tabla tipo HTML para el inodo
	dot.WriteString(fmt.Sprintf("inode%d [label=<\n<TABLE BORDER='1' CELLBORDER='1' CELLSPACING='0'>\n", index))
	dot.WriteString("<TR><TD COLSPAN='2'><B>Inodo " + fmt.Sprint(index) + "</B></TD></TR>\n")
	dot.WriteString("<TR><TD>Type</TD><TD>0</TD></TR>\n")
	dot.WriteString("<TR><TD>ap0</TD><TD>" + fmt.Sprint(inode.IN_Block[0]) + "</TD></TR>\n")
	dot.WriteString("<TR><TD>ap1</TD><TD>" + fmt.Sprint(inode.IN_Block[1]) + "</TD></TR>\n")
	dot.WriteString("<TR><TD>ap2</TD><TD>" + fmt.Sprint(inode.IN_Block[2]) + "</TD></TR>\n")
	dot.WriteString("</TABLE>> shape=plain];\n")

	for i, block := range inode.IN_Block {
		if block == -1 {
			continue
		}

		dot.WriteString(fmt.Sprintf("inode%d -> block%d_%d;\n", index, index, i))

		if i < 12 {
			if string(inode.IN_Type[:]) == "0" || string(inode.IN_Type[:]) == "" {
				var folder Structs.FolderBlock
				blockOffset := int64(sb.SB_Block_Start + block*int32(binary.Size(Structs.FolderBlock{})))
				if err := Utilities.ReadObject(file, &folder, blockOffset, buffer); err != nil {
					fmt.Fprintf(buffer, "Error al leer FolderBlock %d: %v\n", block, err)
					continue
				}

				dot.WriteString(fmt.Sprintf("block%d_%d [label=<\n<TABLE BORDER='1' CELLBORDER='1' CELLSPACING='0'>\n", index, i))
				dot.WriteString("<TR><TD COLSPAN='2'><B>FolderBlock " + fmt.Sprint(block) + "</B></TD></TR>\n")
				for _, content := range folder.B_Content {
					name := strings.Trim(string(content.B_Name[:]), "\x00")
					if name != "" {
						dot.WriteString("<TR><TD>" + name + "</TD><TD>" + fmt.Sprint(content.B_Inode) + "</TD></TR>\n")
					}
				}
				dot.WriteString("</TABLE>> shape=plain fillcolor=salmon style=filled];\n")

				for _, content := range folder.B_Content {
					name := strings.Trim(string(content.B_Name[:]), "\x00")
					if name != "" && content.B_Inode != -1 && name != "." && name != ".." {
						generarArbolInodos(content.B_Inode, sb, file, dot, buffer)
					}
				}

			} else {
				var fileblock Structs.FileBlock
				blockOffset := int64(sb.SB_Block_Start + block*int32(binary.Size(Structs.FileBlock{})))
				if err := Utilities.ReadObject(file, &fileblock, blockOffset, buffer); err != nil {
					fmt.Fprintf(buffer, "Error al leer FileBlock %d: %v\n", block, err)
					continue
				}
				contenido := sanitizeDOTContent(strings.TrimRight(string(fileblock.B_Content[:]), "\x00"))
				dot.WriteString(fmt.Sprintf("block%d_%d [label=<\n<TABLE BORDER='1' CELLBORDER='1' CELLSPACING='0'>\n<TR><TD><B>FileBlock %d</B></TD></TR><TR><TD>%s</TD></TR></TABLE>> shape=plain fillcolor=khaki style=filled];\n", index, i, block, contenido))
			}
		} else {
			var pointerBlock Structs.PointerBlock
			blockOffset := int64(sb.SB_Block_Start + block*int32(binary.Size(Structs.PointerBlock{})))
			if err := Utilities.ReadObject(file, &pointerBlock, blockOffset, buffer); err != nil {
				fmt.Fprintf(buffer, "Error al leer PointerBlock %d: %v\n", block, err)
				continue
			}

			dot.WriteString(fmt.Sprintf("block%d_%d [label=<\n<TABLE BORDER='1' CELLBORDER='1' CELLSPACING='0'>\n<TR><TD COLSPAN='2'><B>PointerBlock %d</B></TD></TR>\n", index, i, block))
			for j, ptr := range pointerBlock.B_Pointers {
				if ptr != -1 {
					dot.WriteString("<TR><TD>ptr" + fmt.Sprint(j) + "</TD><TD>" + fmt.Sprint(ptr) + "</TD></TR>\n")
				}
			}
			dot.WriteString("</TABLE>> shape=plain fillcolor=gray style=filled];\n")

			for _, ptr := range pointerBlock.B_Pointers {
				if ptr != -1 {
					generarArbolInodos(ptr, sb, file, dot, buffer)
				}
			}
		}
	}
}

func sanitizeDOTContent(input string) string {
	safe := ""
	for _, r := range input {
		if unicode.IsLetter(r) || unicode.IsDigit(r) || strings.ContainsRune(".,:_- ", r) {
			safe += string(r)
		} else {
			safe += " "
		}
	}
	return safe
}

func ReporteLS(id string, path string, buffer *bytes.Buffer, pathFileLs string) {
	fmt.Fprintln(buffer, "=========== REPORTE LS ===========")

	if !FileSystem.IsUserLoggedInREPORTE() {
		fmt.Fprintln(buffer, "Error REP LS: No hay una sesión activa.")
		return
	}

	// Obtener partición montada
	particiones := DiskManagement.GetMountedPartitions()
	var diskPath string
	found := false
	for _, parts := range particiones {
		for _, part := range parts {
			if part.ID == id && part.LoggedIn {
				diskPath = part.Path
				found = true
				break
			}
		}
		if found {
			break
		}
	}
	if !found {
		fmt.Fprintf(buffer, "Error REP LS: No se encontró la partición con ID %s\n", id)
		return
	}

	// Abrir disco
	file, err := Utilities.OpenFile(diskPath, buffer)
	if err != nil {
		fmt.Fprintln(buffer, "Error REP LS: No se pudo abrir el disco.")
		return
	}
	defer file.Close()

	// Leer MBR y Superbloque
	var mbr Structs.MRB
	if err := Utilities.ReadObject(file, &mbr, 0, buffer); err != nil {
		fmt.Fprintln(buffer, "Error REP LS: No se pudo leer el MBR.")
		return
	}

	var index = -1
	for i := 0; i < 4; i++ {
		if strings.Contains(string(mbr.MbrPartitions[i].ID[:]), id) && mbr.MbrPartitions[i].Status[0] == '1' {
			index = i
			break
		}
	}
	if index == -1 {
		fmt.Fprintf(buffer, "Error REP LS: Partición con ID %s no válida.\n", id)
		return
	}

	var sb Structs.Superblock
	if err := Utilities.ReadObject(file, &sb, int64(mbr.MbrPartitions[index].Start), buffer); err != nil {
		fmt.Fprintln(buffer, "Error REP LS: No se pudo leer el SuperBloque.")
		return
	}

	// Obtener inodo objetivo
	indexInodo := FileSystem.BuscarInodoPorRutaREPORTE(pathFileLs, file, sb, buffer)
	if indexInodo == -1 {
		fmt.Fprintf(buffer, "Error REP LS: No se encontró la ruta: %s\n", pathFileLs)
		return
	}

	var inode Structs.Inode
	offsetInode := int64(sb.SB_Inode_Start + indexInodo*int32(binary.Size(Structs.Inode{})))
	if err := Utilities.ReadObject(file, &inode, offsetInode, buffer); err != nil {
		fmt.Fprintln(buffer, "Error REP LS: No se pudo leer el inodo objetivo.")
		return
	}

	// Generar el archivo DOT
	var dot bytes.Buffer
	fmt.Fprintln(&dot, "digraph G {")
	fmt.Fprintln(&dot, "node [shape=plaintext];")
	fmt.Fprintln(&dot, "ls [label=<")
	fmt.Fprintln(&dot, "<table border='1' cellborder='1' cellspacing='0'>")
	fmt.Fprintln(&dot, "<tr><td><b>Permisos</b></td><td><b>Owner</b></td><td><b>Grupo</b></td><td><b>Size (en Bytes)</b></td><td><b>Fecha</b></td><td><b>Hora</b></td><td><b>Tipo</b></td><td><b>Name</b></td></tr>")

	// Leer y mostrar información de archivos/carpetas del inodo objetivo
	for _, blk := range inode.IN_Block {
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
			if name == "" || name == "." || name == ".." {
				continue
			}
			var hijo Structs.Inode
			hijoOffset := int64(sb.SB_Inode_Start + entry.B_Inode*int32(binary.Size(Structs.Inode{})))
			if err := Utilities.ReadObject(file, &hijo, hijoOffset, buffer); err != nil {
				continue
			}

			perm := string(hijo.IN_Perm[:])
			uid := hijo.IN_Uid
			gid := hijo.IN_Gid
			size := hijo.IN_Size
			fecha := strings.Trim(string(hijo.IN_Ctime[:]), "\x00")
			hora := strings.Trim(string(hijo.IN_Mtime[:]), "\x00")
			tipo := "Archivo"
			if hijo.IN_Block[0] != -1 {
				tipo = "Carpeta"
			}

			fmt.Fprintf(&dot, "<tr><td>%s</td><td>User%d</td><td>Group%d</td><td>%d</td><td>%s</td><td>%s</td><td>%s</td><td>%s</td></tr>\n",
				perm, uid, gid, size, fecha, hora, tipo, name)
		}
	}

	fmt.Fprintln(&dot, "</table>")
	fmt.Fprintln(&dot, ">];")
	fmt.Fprintln(&dot, "}")

	// Guardar el archivo .dot
	dotPath := strings.ReplaceAll(path, ".jpg", ".dot")
	if err := os.WriteFile(dotPath, dot.Bytes(), 0644); err != nil {
		fmt.Fprintf(buffer, "Error REP LS: No se pudo escribir el archivo DOT.\n")
		return
	}

	// Asegurar carpeta de destino
	dir := filepath.Dir(path)
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		_ = os.MkdirAll(dir, 0755)
	}

	// Ejecutar Graphviz
	cmd := exec.Command("dot", "-Tjpg", dotPath, "-o", path)
	var stderr bytes.Buffer
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		fmt.Fprintf(buffer, "Error REP LS: Graphviz falló: %s\n", stderr.String())
		return
	}

	fmt.Fprintf(buffer, "Reporte LS generado exitosamente en: %s\n", path)
}

func ReporteFile(id string, outputPath string, buffer *bytes.Buffer, filePathLs string) {
	fmt.Fprintln(buffer, "=========== REPORTE FILE ===========")

	if !FileSystem.IsUserLoggedInREPORTE() {
		fmt.Fprintln(buffer, "Error REP FILE: No hay una sesión activa.")
		return
	}

	// Obtener partición activa
	particiones := DiskManagement.GetMountedPartitions()
	var diskPath string
	found := false
	for _, parts := range particiones {
		for _, part := range parts {
			if part.ID == id && part.LoggedIn {
				diskPath = part.Path
				found = true
				break
			}
		}
		if found {
			break
		}
	}
	if !found {
		fmt.Fprintf(buffer, "Error REP FILE: No se encontró la partición con ID %s\n", id)
		return
	}

	// Abrir archivo del disco
	file, err := Utilities.OpenFile(diskPath, buffer)
	if err != nil {
		fmt.Fprintln(buffer, "Error REP FILE: No se pudo abrir el archivo del disco.")
		return
	}
	defer file.Close()

	// Leer MBR y Superbloque
	var mbr Structs.MRB
	if err := Utilities.ReadObject(file, &mbr, 0, buffer); err != nil {
		fmt.Fprintln(buffer, "Error REP FILE: No se pudo leer el MBR.")
		return
	}

	var partitionIndex = -1
	for i := 0; i < 4; i++ {
		if strings.Contains(string(mbr.MbrPartitions[i].ID[:]), id) && mbr.MbrPartitions[i].Status[0] == '1' {
			partitionIndex = i
			break
		}
	}
	if partitionIndex == -1 {
		fmt.Fprintln(buffer, "Error REP FILE: Partición no válida o no montada.")
		return
	}

	var sb Structs.Superblock
	if err := Utilities.ReadObject(file, &sb, int64(mbr.MbrPartitions[partitionIndex].Start), buffer); err != nil {
		fmt.Fprintln(buffer, "Error REP FILE: No se pudo leer el superbloque.")
		return
	}

	// Buscar inodo del archivo
	indexInode := FileSystem.BuscarInodoPorRutaREPORTE(filePathLs, file, sb, buffer)
	if indexInode == -1 {
		fmt.Fprintf(buffer, "Error REP FILE: No se encontró el archivo %s\n", filePathLs)
		return
	}

	var inode Structs.Inode
	offsetInode := int64(sb.SB_Inode_Start + indexInode*int32(binary.Size(Structs.Inode{})))
	if err := Utilities.ReadObject(file, &inode, offsetInode, buffer); err != nil {
		fmt.Fprintln(buffer, "Error REP FILE: No se pudo leer el inodo del archivo.")
		return
	}

	// Leer contenido del archivo desde los bloques
	var contenido strings.Builder
	for _, block := range inode.IN_Block {
		if block == -1 {
			continue
		}
		var fb Structs.FileBlock
		offset := int64(sb.SB_Block_Start + block*int32(binary.Size(Structs.FileBlock{})))
		if err := Utilities.ReadObject(file, &fb, offset, buffer); err != nil {
			continue
		}
		contenido.WriteString(strings.TrimRight(string(fb.B_Content[:]), "\x00"))
	}

	// Crear archivo de salida
	dir := filepath.Dir(outputPath)
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		_ = os.MkdirAll(dir, 0755)
	}

	if err := os.WriteFile(outputPath, []byte(contenido.String()), 0644); err != nil {
		fmt.Fprintln(buffer, "Error REP FILE: No se pudo escribir el archivo de salida.")
		return
	}

	fmt.Fprintf(buffer, "Reporte de archivo generado exitosamente en: %s\n", outputPath)
}

func corregirExtensionTxt(path string) string {
	if filepath.Ext(path) != ".txt" {
		return strings.TrimSuffix(path, filepath.Ext(path)) + ".txt"
	}
	return path
}

func escapeDOTString(s string) string {
	s = strings.TrimSpace(s)
	s = strings.ReplaceAll(s, "\\", "\\\\")
	s = strings.ReplaceAll(s, "\"", "\\\"")
	s = strings.ReplaceAll(s, "\n", "\\n")
	s = strings.ReplaceAll(s, "\r", "\\r")
	s = strings.ReplaceAll(s, "{", "\\{")
	s = strings.ReplaceAll(s, "}", "\\}")
	s = strings.ReplaceAll(s, "|", "\\|")
	s = strings.ReplaceAll(s, "<", "\\<")
	s = strings.ReplaceAll(s, ">", "\\>")
	return s
}
