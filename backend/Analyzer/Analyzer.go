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
	"proyecto1/Report"
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
		buffer.WriteString("\n") // ← Este es el que faltaba

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
	} else if strings.Contains(command, "mkusr") {
		fn_mkusr(params, buffer)
	} else if strings.Contains(command, "rmusr") {
		fn_rmusr(params, buffer)
	} else if strings.Contains(command, "rep") {
		fn_rep(params, buffer)
	} else if strings.Contains(command, "mkdir") {
		fn_mkdir(params, buffer)
	} else if strings.Contains(command, "mkfile") {
		fn_mkfile(params, buffer)
	} else if strings.Contains(command, "rmgrp") {
		fn_rmgrp(params, buffer)
	} else if strings.Contains(command, "chgrp") {
		fn_chgrp(params, buffer)
	}else {
		fmt.Fprintf(buffer, "Error: Commando invalido o no encontrado")
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

	if *unit == "" {
		*unit = "k"
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
			fmt.Fprintf(buffer, "Error: comando 'rmdsik' inclyte parametros no asociados\n")
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
		fmt.Fprintf(buffer, "Error: Size must be greater than 0")
		return
	}

	if *path == "" {
		fmt.Fprintf(buffer, "Error: Path is required")
		return
	}

	// Si no se proporcionó un fit, usar el valor predeterminado "w"
	if *fit == "" {
		*fit = "w"
	}

	if *fit == "bf" {
		*fit = "b"
	}

	if *fit == "wf" {
		*fit = "w"
	}

	if *fit == "ff" {
		*fit = "f"
	}

	if *unit == "" {
		*unit = "k"
	}

	fmt.Fprintf(buffer, "EL FIT ES : %s\n", *fit)

	// Validar fit (b/w/f)
	if *fit != "b" && *fit != "f" && *fit != "w" {
		fmt.Fprintf(buffer, "Error: Fit must be 'b', 'f', or 'w'")
		return
	}

	if *unit != "k" && *unit != "m" && *unit != "b" {
		fmt.Fprintf(buffer, "Error: Unit must be 'k' or 'm' or 'b'")
		return
	}

	if *type_ != "p" && *type_ != "e" && *type_ != "l" {
		fmt.Fprintf(buffer, "Error: Type must be 'p', 'e', or 'l'")
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

	fmt.Println("MBR después de FDISK:")
	// Llamar a la función
	DiskManagement.Fdisk(*size, *path, *name, *unit, *type_, *fit, buffer.(*bytes.Buffer))
	//Structs.PrintMBRP(mbr)
	fmt.Println("=============DESPUES DE FDISK===================\n")
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
			fmt.Fprintf(buffer, "Error: Flag not found")
		}
	}

	// Verifica que se hayan establecido todas las flags necesarias
	if *id == "" {
		fmt.Fprintf(buffer, "Error: id es un parámetro obligatorio.")
		return
	}

	if *type_ == "" { //2fs 3fs
		*type_ = "full"
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
			fmt.Fprintf(buffer, "Error: Flag not found")
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
	fmt.Fprintf(buffer, "=========== CAT ===========\n")
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
	fmt.Fprintf(buffer, "==============FIN CAT===============\n")
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

func fn_rep(input string, buffer io.Writer) {
	fs := flag.NewFlagSet("rep", flag.ExitOnError)
	nombre := fs.String("name", "", "Nombre")
	ruta := fs.String("path", "full", "Ruta")
	ID := fs.String("id", "", "IDParticion")
	path_file_ls := fs.String("path_file_ls", "", "PathFile")

	fs.Parse(os.Args[1:])
	matches := re.FindAllStringSubmatch(input, -1)

	for _, match := range matches {
		nombreFlag := match[1]
		valorFlag := strings.ToLower(match[2])

		valorFlag = strings.Trim(valorFlag, "\"")

		switch nombreFlag {
		case "name", "path", "id", "path_file_ls":
			fs.Set(nombreFlag, valorFlag)
		default:
			fmt.Fprintf(buffer, "Error: El comando 'REP' incluye parámetros no asociados.\n")
			return
		}
	}
	Report.Rep(*nombre, *ruta, *ID, *path_file_ls, buffer.(*bytes.Buffer))
}

//--------------------Función para mkdir--------------------
func fn_mkdir(input string, buffer io.Writer) {
	//fmt.Fprintf(buffer, "DEBUG: Entrando a fn_mkdir con input: %s\n", input)
	var path string
	var p bool = false

	// Dividir en tokens por espacios
	tokens := strings.Fields(input)

	for _, token := range tokens {
		if strings.HasPrefix(token, "-path=") {
			path = strings.Trim(strings.SplitN(token, "=", 2)[1], "\"")
		} else if token == "-p" {
			p = true
		} else if strings.HasPrefix(token, "-") {
			fmt.Fprintf(buffer, "Error: El comando 'MKDIR' incluye parámetros no asociados: %s\n", token)
			return
		}
	}

	//fmt.Fprintf(buffer, "DEBUG: flag -p = %v\n", p)

	if path == "" {
		fmt.Fprintf(buffer, "Error: MKDIR requiere parámetro obligatorio -path.\n")
		return
	}
		
		//fmt.Fprintf(buffer, "DEBUG: flag -p = %v\n", p)

		// Llamar a la función final
		FileSystem.Mkdir(path, p, buffer.(*bytes.Buffer))
}

func fn_mkfile(input string, buffer io.Writer) {
	fmt.Fprintf(buffer, "=========== MKFILE ===========")
	var path, cont string
	var p bool
	var size int = 0

	re := regexp.MustCompile(`-(\w+)=?("[^"]*"|\S*)`)
	matches := re.FindAllStringSubmatch(input, -1)

	for _, match := range matches {
		flagName := strings.ToLower(match[1])
		originalValue := match[2]                  // <-- Conserva el original
		flagValue := strings.Trim(match[2], "\"")  // <-- El que usarás

		switch flagName {
		case "path":
			// Usa el valor original, no lo bajes a minúscula
			path = strings.Trim(originalValue, "\"")
		case "cont":
			cont = strings.Trim(originalValue, "\"")
		case "size":
			val, err := strconv.Atoi(flagValue)
			if err != nil || val < 0 {
				fmt.Fprintf(buffer, "Error: el parámetro -size debe ser un número entero no negativo.\n")
				return
			}
			size = val
		case "p":
			if match[2] == "" {
				p = true
			} else {
				fmt.Fprintf(buffer, "Error: el flag -p no debe llevar valor.\n")
				return
			}
		case "r":
			if match[2] == "" {
				p = true
			} else {
				fmt.Fprintf(buffer, "Error: el flag -r no debe llevar valor.\n")
				return
			}
		
		default:
			fmt.Fprintf(buffer, "Error: parámetro no reconocido -%s\n", flagName)
			return
		}
	}

	if path == "" {
		fmt.Fprintf(buffer, "Error: el parámetro -path es obligatorio para MKFILE.\n")
		return
	}

	if cont != "" {

		aliasMap := map[string]string{
			"/home/matius/escritorio": "/home/matius/Escritorio",
		}
	
		for k, v := range aliasMap {
			if strings.HasPrefix(cont, k) {
				cont = strings.Replace(cont, k, v, 1)
				break
			}
		}
		
		if _, err := os.Stat(cont); os.IsNotExist(err) {
			fmt.Fprintf(buffer, "Error: el archivo %s no existe en el sistema local.\n", cont)
			return
		}
		// Aquí podrías cargar el contenido real del archivo si deseas
		// contentBytes, _ := os.ReadFile(cont)
		// content = string(contentBytes)
	} else if size > 0 {
		// Si no hay -cont, se genera contenido según size
		var builder strings.Builder
		for i := 0; i < size; i++ {
			builder.WriteByte(byte('0' + i%10))
		}
		cont = builder.String()
	}

	// Llama a la función en el FileSystem
	FileSystem.Mkfile(path, p, cont, buffer.(*bytes.Buffer))
}

//--------------------Función para rmgrp--------------------
func fn_rmgrp(input string, buffer io.Writer) {
	fs := flag.NewFlagSet("rmgrp", flag.ExitOnError)
	name := fs.String("name", "", "Nombre")

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
			fmt.Fprintf(buffer, "Error: El comando 'RMGRP' incluye parámetros no asociados.\n")
			return
		}
	}
	User.Rmgrp(*name, buffer.(*bytes.Buffer))
}

func fn_chgrp(params string, buffer io.Writer) {
	// Convertir buffer a *bytes.Buffer para pasar al método del paquete User
	buff := buffer.(*bytes.Buffer)

	// Inicializar variables
	var user string
	var group string

	// Separar los parámetros por espacios
	paramList := strings.Fields(params)

	for _, param := range paramList {
		param = strings.ToLower(param)

		if strings.HasPrefix(param, "-user=") {
			user = strings.TrimPrefix(param, "-user=")
		} else if strings.HasPrefix(param, "-grp=") {
			group = strings.TrimPrefix(param, "-grp=")
		}
	}

	// Validar que ambos parámetros hayan sido proporcionados
	if user == "" || group == "" {
		fmt.Fprint(buff, "Error CHGRP: Faltan parámetros obligatorios (-user y -grp).\n")
		return
	}

	// Llamar al método que realiza la lógica
	User.Chgrp(user, group, buff)
}
