package Analyzer

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"proyecto1/DiskManagement"
	"proyecto1/FileSystem"
	"proyecto1/User"
	"regexp"
	"strings"
	"io"
	"bytes"
	"strconv"
	"proyecto1/Structs"
	"proyecto1/Utilities"
)

var re = regexp.MustCompile(`-(\w+)=("[^"]+"|\S+)`)

//-nombre=valor

//input := "mkdisk -size=3000 -unit=K -fit=BF -path=/home/bang/Disks/disk1.bin"

/*
parts[0] es "mkdisk"
*/

func getCommandAndParams(input string) (string, string) {
	parts := strings.Fields(input)
	if len(parts) > 0 {
		command := strings.ToLower(parts[0])
		for i := 1; i < len(parts); i++ {
			parts[i] = strings.ToLower(parts[i])
		}
		params := strings.Join(parts[1:], " ")
		return command, params
	}
	return "", input

	/*Después de procesar la entrada:
	command será "mkdisk".
	params será "-size=3000 -unit=K -fit=BF -path=/home/bang/Disks/disk1.bin".*/
}

//----------------Analizador de entrada----------------
func Analyze(entrada string) string {
	var buffer bytes.Buffer
	scanner := bufio.NewScanner(strings.NewReader(entrada))
	for scanner.Scan() {

		input := scanner.Text()
		
		if len(input) == 0 || input[0] == '#' {
			fmt.Fprintf(&buffer, "%s\n", input)
			continue
		}
		input = strings.TrimSpace(input)
		command, params := getCommandAndParams(input)

		fmt.Println("Comando: ", command, " - ", "Parametro: ", params)

		AnalyzeCommnad(command, params, &buffer)

		//mkdisk -size=3000 -unit=K -fit=BF -path="/home/bang/Disks/disk1.bin"
	}
	return buffer.String()
}

//--------------------Analizador de tipo de comando--------------------
func AnalyzeCommnad(command string, params string, buffer io.Writer) {

	if strings.Contains(command, "mkdisk") {
		fn_mkdisk(params, buffer)
	} else if strings.Contains(command, "rmdisk") {
		fn_rmdisk(params, buffer)
	} else if strings.Contains(command, "mounted") {
		fn_mounted(buffer)
	} else if strings.Contains(command, "fdisk") {
		fn_fdisk(params, buffer)
	} else if strings.Contains(command, "mount") {
		fn_mount(params, buffer)
	} else if strings.Contains(command, "mkfs") {
		fn_mkfs(params, buffer)
	} else if strings.Contains(command, "login") {
		fn_login(params, buffer)
	} else if strings.Contains(command, "logout") {
		fn_logout(params, buffer)
	} else if strings.Contains(command, "cat") {
		fn_cat(params, buffer)
	} else if strings.Contains(command, "mkgrp") {
		fn_mkgrp(params, buffer)
	} else if strings.Contains(command, "list") {
		fn_list(params, buffer)
	} else if strings.Contains(command, "mkusr") {
		fn_mkusr(params, buffer)
	} else if strings.Contains(command, "rmusr") {
		fn_rmusr(params, buffer)
	} else {
		fmt.Println("Error: Commando invalido o no encontrado")
	}

}

//--------------------Función para mkdisk--------------------
func fn_mkdisk(params string, buffer io.Writer) {
	// Definir flag
	fs := flag.NewFlagSet("mkdisk", flag.ExitOnError)
	size := fs.Int("size", 0, "Tamaño")
	fit := fs.String("fit", "ff", "Ajuste")
	unit := fs.String("unit", "m", "Unidad")
	path := fs.String("path", "", "Ruta")

	// Parse flag
	fs.Parse(os.Args[1:])

	// Encontrar la flag en el input
	matches := re.FindAllStringSubmatch(params, -1)

	// Process the input
	for _, match := range matches {
		flagName := match[1]                   // match[1]: Captura y guarda el nombre del flag (por ejemplo, "size", "unit", "fit", "path")
		flagValue := strings.ToLower(match[2]) //trings.ToLower(match[2]): Captura y guarda el valor del flag, asegurándose de que esté en minúsculas

		flagValue = strings.Trim(flagValue, "\"")

		switch flagName {
		case "size", "fit", "unit", "path":
			fs.Set(flagName, flagValue)
		default:
			fmt.Fprintf(buffer, "Error: Flag not found")
			return
		}
	}

	// Validaciones
	if *size <= 0 {
		fmt.Fprintf(buffer, "Error: Size must be greater than 0")
		return
	}

	if *fit != "bf" && *fit != "ff" && *fit != "wf" {
		fmt.Fprintf(buffer, "Error: Fit must be 'bf', 'ff', or 'wf'")
		return
	}

	if *unit != "k" && *unit != "m" {
		fmt.Fprintf(buffer, "Error: Unit must be 'k' or 'm'")
		return
	}

	if *path == "" {
		fmt.Fprintf(buffer, "Error: Path is required")
		return
	}

	// LLamamos a la funcion
	DiskManagement.Mkdisk(*size, *fit, *unit, *path, buffer.(*bytes.Buffer))
}

//--------------------Función para rmdisk--------------------
func fn_rmdisk(params string, buffer io.Writer) {
	fs := flag.NewFlagSet("rmdisk", flag.ExitOnError)
	ruta := fs.String("path", "", "Ruta")

	fs.Parse(os.Args[1:])
	matches := re.FindAllStringSubmatch(params, -1)

	for _, match := range matches {
		nombreFlag := match[1]
		valorFlag := strings.ToLower(match[2])
		valorFlag = strings.Trim(valorFlag, "\"")

		switch nombreFlag {
		case "path":
			fs.Set(nombreFlag, valorFlag)
		default:
			fmt.Println(buffer, "Error: comando 'rmdsik' inclyte parametros no asociados\n")
			return
		}
	}

	// Llamas a la función para borrar el disco aquí.
	DiskManagement.Rmdisk(*ruta, buffer.(*bytes.Buffer))
}

//--------------------Función para fdisk--------------------
func fn_fdisk(input string, buffer io.Writer) {
	// Definir flags
	fs := flag.NewFlagSet("fdisk", flag.ExitOnError)
	size := fs.Int("size", 0, "Tamaño")
	path := fs.String("path", "", "Ruta")
	name := fs.String("name", "", "Nombre")
	unit := fs.String("unit", "k", "Unidad") //por defecto en KiloBytes
	type_ := fs.String("type", "p", "Tipo")
	fit := fs.String("fit", "", "Ajuste") // Dejar fit vacío por defecto

	// Parsear los flags
	fs.Parse(os.Args[1:])

	// Encontrar los flags en el input
	matches := re.FindAllStringSubmatch(input, -1)

	// Procesar el input
	for _, match := range matches {
		flagName := match[1]
		flagValue := strings.ToLower(match[2])

		flagValue = strings.Trim(flagValue, "\"")

		switch flagName {
		case "size", "fit", "unit", "path", "name", "type":
			fs.Set(flagName, flagValue)
		default:
			fmt.Fprintf(buffer, "Error: El comando 'FDISK' incluye parámetros no asociados.\n")
			return
		}
	}

	// Validaciones
	if *size <= 0 {
		fmt.Println("Error: Size must be greater than 0")
		return
	}

	if *path == "" {
		fmt.Println("Error: Path is required")
		return
	}

	// Si no se proporcionó un fit, usar el valor predeterminado "w"
	if *fit == "" {
		*fit = "w"
	}

		// Validar fit (b/w/f)
	if *fit != "b" && *fit != "f" && *fit != "w" {
		fmt.Println("Error: Fit must be 'b', 'f', or 'w'")
		return
	}

	if *unit != "k" && *unit != "m" {
		fmt.Println("Error: Unit must be 'k' or 'm'")
		return
	}

	if *type_ != "p" && *type_ != "e" && *type_ != "l" {
		fmt.Println("Error: Type must be 'p', 'e', or 'l'")
		return
	}

	// Abrir disco para mostrar estado del MBR luego de fdisk
	file, err := Utilities.OpenFile(*path, buffer.(*bytes.Buffer))
	if err != nil {
		fmt.Fprintf(buffer, "Error abriendo disco después de FDISK: %s\n", err)
		return
	}
	defer file.Close()

	var mbr Structs.MRB
	if err := Utilities.ReadObject(file, &mbr, 0, buffer.(*bytes.Buffer)); err != nil {
		fmt.Fprintf(buffer, "Error leyendo MBR después de FDISK: %s\n", err)
		return
	}

	// Llamar a la función
	DiskManagement.Fdisk(*size, *path, *name, *unit, *type_, *fit, buffer.(*bytes.Buffer))
	Structs.PrintMBRP(mbr)
}

//--------------------Función para mount--------------------
func fn_mount(input string, buffer io.Writer) {
	fs := flag.NewFlagSet("mount", flag.ExitOnError)
	path := fs.String("path", "", "Ruta")
	name := fs.String("name", "", "Nombre")

	fs.Parse(os.Args[1:])
	matches := re.FindAllStringSubmatch(input, -1)

	for _, match := range matches {
		nameFlag := match[1]
		valueFlag := strings.ToLower(match[2])

		valueFlag = strings.Trim(valueFlag, "\"")

		switch nameFlag {
		case "path", "name":
			fs.Set(nameFlag, valueFlag)
		default:
			fmt.Fprintf(buffer, "Error: El comando 'MOUNT' incluye parámetros no asociados.\n")
			return
		}
	}
	DiskManagement.Mount(*path, *name, buffer.(*bytes.Buffer))
}

//--------------------Función para mkfs--------------------
func fn_mkfs(input string, buffer io.Writer) {
	fs := flag.NewFlagSet("mkfs", flag.ExitOnError)
	id := fs.String("id", "", "Id")
	type_ := fs.String("type", "", "Tipo")
	fs_ := fs.String("fs", "2fs", "Fs")

	// Parse the input string, not os.Args
	matches := re.FindAllStringSubmatch(input, -1)

	for _, match := range matches {
		flagName := match[1]
		flagValue := match[2]

		flagValue = strings.Trim(flagValue, "\"")

		switch flagName {
		case "id", "type", "fs":
			fs.Set(flagName, flagValue)
		default:
			fmt.Println("Error: Flag not found")
		}
	}

	// Verifica que se hayan establecido todas las flags necesarias
	if *id == "" {
		fmt.Println("Error: id es un parámetro obligatorio.")
		return
	}

	if *type_ == "" { //2fs 3fs
		fmt.Println("Error: type es un parámetro obligatorio.")
		return
	}

	// Llamar a la función
	FileSystem.Mkfs(*id, *type_, *fs_, buffer.(*bytes.Buffer))
}

//--------------------Función para login--------------------
func fn_login(input string, buffer io.Writer) {
	fs := flag.NewFlagSet("login", flag.ExitOnError)
	user := fs.String("user", "", "Usuario")
	pass := fs.String("pass", "", "Contraseña")
	id := fs.String("id", "", "Id")

	fs.Parse(os.Args[1:])
	matches := re.FindAllStringSubmatch(input, -1)

	for _, match := range matches {
		flagName := match[1]
		flagValue := match[2]

		flagValue = strings.Trim(flagValue, "\"")

		switch flagName {
		case "user", "pass", "id":
			fs.Set(flagName, flagValue)
		default:
			fmt.Println("Error: Flag not found")
		}
	}

	User.Login(*user, *pass, *id, buffer.(*bytes.Buffer))

}

//--------------------Función para logout--------------------
func fn_logout(input string, buffer io.Writer) {
	input = strings.TrimSpace(input)
	if len(input) > 0 {
		fmt.Fprintf(buffer, "Error: El comando 'LOGOUT' incluye parámetros no asociados.\n")
		return
	}
	User.LogOut(buffer.(*bytes.Buffer))
}

//--------------------Función para cat--------------------
func fn_cat(params string, buffer io.Writer) {
	files := make(map[int]string)
	matches := re.FindAllStringSubmatch(params, -1)

	for _, match := range matches {
		flagName := match[1]
		flagValue := strings.ToLower(match[2])

		flagValue = strings.Trim(flagValue, "\"")

		if strings.HasPrefix(flagName, "file") {

			NUmber, err := strconv.Atoi(strings.TrimPrefix(flagName, "file"))
			if err != nil {
				fmt.Fprintf(buffer, "Error: Nombre de archivo inválido")
				return
			}
			files[NUmber] = flagValue
		} else {
			fmt.Fprintf(buffer, "Error: Flag not found")
		}
	}
	var orden []string
	for i := 1; i <= len(files); i++ {
		if file, exists := files[i]; exists {
			orden = append(orden, file)
		} else {
			fmt.Fprintf(buffer, "Error: Falta un archivo en la secuencia")
			return
		}
	}
	if len(orden) == 0 {
		fmt.Fprintf(buffer, "Error: No se encontraron archivos")
		return
	}
	FileSystem.CAT(orden, buffer.(*bytes.Buffer))
}

//--------------------Función para mkgrp--------------------
func fn_mkgrp(input string, buffer io.Writer) {
	fs := flag.NewFlagSet("mkgrp ", flag.ExitOnError)
	nombre := fs.String("name", "", "Nombre")

	fs.Parse(os.Args[1:])
	matches := re.FindAllStringSubmatch(input, -1)

	for _, match := range matches {
		nombreFlag := match[1]
		valorFlag := strings.ToLower(match[2])

		valorFlag = strings.Trim(valorFlag, "\"")

		switch nombreFlag {
		case "name":
			fs.Set(nombreFlag, valorFlag)
		default:
			fmt.Fprintf(buffer, "Error: El comando 'MKGRP' incluye parámetros no asociados.\n")
			return
		}
	}
	User.Mkgrp(*nombre, buffer.(*bytes.Buffer))
}

//--------------------Función para comando_list--------------------
func fn_list(input string, buffer io.Writer) {
	input = strings.TrimSpace(input)
	if len(input) > 0 {
		fmt.Fprintf(buffer, "Error: El comando 'LIST' incluye parámetros no asociados.\n")
		return
	}
	DiskManagement.List(buffer.(*bytes.Buffer))
}

//--------------------Función para comando_mounted--------------------
func fn_mounted(buffer io.Writer) {
	fmt.Fprintf(buffer, "===== PARTICIONES MONTADAS =====\n")
	if len(DiskManagement.MountedPartitions) == 0 {
		fmt.Fprintf(buffer, "No hay particiones montadas.\n")
		return
	}

	// Iterar sobre particiones montadas e imprimir sus IDs
	for _, partitions := range DiskManagement.MountedPartitions {
		for _, particion := range partitions {
			fmt.Fprintf(buffer, "- %s\n", particion.ID)
		}
	}
	fmt.Fprintf(buffer, "================================\n")
}

//--------------------Función para mkusr--------------------
func fn_mkusr(input string, buffer io.Writer) {
	fs := flag.NewFlagSet("mkusr", flag.ExitOnError)
	user := fs.String("user", "", "Usuario")
	pass := fs.String("pass", "", "Contraseña")
	grp := fs.String("grp", "", "Grupo")

	matches := re.FindAllStringSubmatch(input, -1)

	for _, match := range matches {
		flagName := strings.ToLower(match[1])
		flagValue := strings.Trim(match[2], "\"")

		switch flagName {
		case "user", "pass", "grp":
			fs.Set(flagName, flagValue)
		default:
			fmt.Fprintf(buffer, "Error: El comando 'MKUSR' incluye parámetros no asociados.\n")
			return
		}
	}

	if *user == "" || *pass == "" || *grp == "" {
		fmt.Fprintf(buffer, "Error: MKUSR requiere obligatoriamente parámetros -user, -pass y -grp.\n")
		return
	}

	User.Mkusr(*user, *pass, *grp, buffer.(*bytes.Buffer))
}

//--------------------Función para rmusr--------------------
func fn_rmusr(input string, buffer io.Writer) {
	fs := flag.NewFlagSet("rmusr", flag.ExitOnError)
	user := fs.String("user", "", "Usuario a eliminar")

	matches := re.FindAllStringSubmatch(input, -1)

	// Procesar los parámetros de entrada
	for _, match := range matches {
		flagName := strings.ToLower(match[1])
		flagValue := strings.Trim(match[2], "\"")

		switch flagName {
		case "user":
			fs.Set(flagName, flagValue)
		default:
			fmt.Fprintf(buffer, "Error: El comando 'RMUSR' incluye parámetros no asociados.\n")
			return
		}
	}

	// Verificar que se ha pasado el parámetro -user
	if *user == "" {
		fmt.Fprintf(buffer, "Error: El comando 'RMUSR' requiere el parámetro -user.\n")
		return
	}

	// Llamar a la función correspondiente en User.go
	User.Rmusr(*user, buffer.(*bytes.Buffer))
}

//rmgrp
//rmusr
//chgrp
//mkdir
//rep
//mkfile