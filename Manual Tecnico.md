# Manual Técnico de Sistema de archivos

Este documento proporciona detalles técnicos sobre el sistema de archivos, incluyendo la estructura del código, la lógica detrás de las funcionalidades principales y la forma en que se gestionan los datos. Está dirigido a desarrolladores que desean entender, mantener o expandir el sistema.

## Requisitos Técnicos

- **Lenguaje de programación:** Go
- **Compilador recomendado:** Visual Studio Code
- **Entorno de desarrollo:** Go
- **Herramientas adicionales:**
  - [Graphviz](https://graphviz.org/): Para la generación de gráficos en formato `.dot`.

## Estructura del Proyecto

El presente es la fase 1 del proyecto del curso Manejo e Implementacion de archivos, dicho proyecto conserva tiene como objetivo desarrollar una aplicación web para la interacción y gestión de un sistema de archivos EXT2. Esta herramienta web moderna permitirá acceder y administrar el sistema de archivos desde cualquier lugar y en cualquier sistema operativo. En el backend, el sistema de archivos se gestionará en una distribución Linux, que atenderá todas las solicitudes provenientes del frontend.

El proyecto está dividido en múltiples archivos para mejorar la modularidad y organización, sin embargo a continuación se definen los más importantes:

#### Main
- **`main`:** Es el archivo main, por medio de este se llama al metodo de analyzar comando con el cual se ejecuta todo el funcionamiento del programa principal.

Dicho módulo fue programado para el uso exclusivo del administrador de la red social. A dicho módulo se puede acceder por medio de las credenciales **`Correo: admin@gmail.com`** y con contraseña **`contraseña: EDD2S2024`**. Al momento de ingresar el administrador tiene diferentes espacios y funciones con los que puede interactuar, espacios que serás descritos a continuación:


#### Analyzer
Dicho archivo contiene la logica principal del programa, posee varios metodos encargados del procesamiento de los comandos que son utilizados a lo largo del proyecto.

ESte posee un metodo principal llamado **`AnalyzeCommand`** el cual se encarga de leer y procesar los comandos que sean enviados desde el frontend.

```go
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
		fmt.Println("Error: Commando invalido o no encontrado")
	}
}
```
Como se puede observar esto se realiza con ayuda de condiciones, una vez se obtenga el nombre del comando que se este recibiendo se llama al metodo que se encarga de procesar dicha operacion.

A continuacion, se presenta un ejemplo del funcionamiento de uno de los comandos. Se muestra que pasaria si se ingresa el comando mkdisk, el metodo que sera ejecutado es **`fn_mkdisk`**, dichos metodos comparten una estructura muy parecida, por lo tanto se toma como ejemplo este metodo.

```` go
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
````
Dicho metodo se encarga de verificar los parametros que posee cada comando, que venga todo en orden, y en caso de que alguno de ellos venga mal escrito o no se encuentre entonces lanza un mensaje de error informando al usuario, en el caso de que todo este en orden, se llama al metodo **`Mkdisk`** del archivo  **`DiskManagement`**, el cual se encarga de llevar a cabo las instrucciones correspondientes para realizar las acciones del respectivo comando.

Cada comando hace uso de algun archivo para realizar sus operaciones, por lo tanto, seria un error pensar que DiskManagement es el que se encarga de procesar cada uno de los metodos usados para el funcionamiento del programa.

#### DiskManagement
Este archivo contiene codigo sobre funciones del disco, asi como de sus particiones.

Como se menciono anteriormente existen metodos de los comandos del sistema repartidos en difentes archivos, los comandos que se encuentran en este son los siguientes:

- Mkdisk
- Rmdisk
- Fdisk
- Mount
- List

Para su correcta implementacion se usaron metodos auxiliares que facilitaron ciertas secciones, entre ellos se pueden mencionar los metodos para eliminar un disk por medio de su ruta, generar id de un disco, obtener el ultimo disco montado, entre otros.

##### Mkdisk
```` go
func Mkdisk(size int, fit string, unit string, path string, buffer *bytes.Buffer ) {
	fmt.Fprintf(buffer, "======INICIO MKDISK======\n")
	fmt.Println("Size:", size)
	fmt.Println("Fit:", fit)
	fmt.Println("Unit:", unit)
	fmt.Println("Path:", path)

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

	fmt.Fprintf(buffer, "======FIN MKDISK======")
}

````

Dicho metodo es el encargado de crear un disco binario simulando un disco físico, escribiendo en él un MBR (Master Boot Record) y llenándolo inicialmente con bytes vacíos (0).

###### Parámetros
- size int: Tamaño del disco a crear. Obligatorio. Debe ser mayor a 0.

- fit string: Tipo de ajuste de partición. Debe ser uno de: bf, ff, wf.

- unit string: Unidad del tamaño. Puede ser "k" para kilobytes o "m" para megabytes.

- path string: Ruta completa (con nombre de archivo) donde se creará el disco. Obligatoria.

- buffer *bytes.Buffer: Buffer de log para registrar mensajes (útil para redirigir mensajes a consola o archivos).

##### Funcionamiento

Conversión del tamaño:

- Si unit es "k", se multiplica por 1024. 
- Si unit es "m", se multiplica por 1024 * 1024.

Creación del archivo binario:

- Se usa Utilities. **`CreateFile(path)`** para crear el archivo.

- Luego se abre con Utilities.**`OpenFile`**.

Inicialización del contenido del disco:

- Se escribe byte por byte con valor 0 hasta llenar el tamaño del disco.

Creación del MBR (Master Boot Record):

- Se inicializa una estructura Structs.MRB (debería contener: MbrSize, Signature, Fit, CreationDate).

- Se escribe al inicio del archivo binario (posición 0).

Lectura y validación del MBR:

- Se lee nuevamente el MBR desde el disco y se imprime con Structs.PrintMBR.

##### Rmdisk
```` go
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
````
Este metodo se encarga de eliminar un disco del sistema, mediante una ruta especifica.

###### Parámetros
- path string: Ruta absoluta del archivo .mia a eliminar. Obligatorio.

- buffer *bytes.Buffer: Buffer de salida para registrar mensajes de log.

Verifica que el parámetro path no esté vacío:

- Llama a Utilities.DeleteFile() para eliminar el archivo, si el archivo no existe o ocurre un error, se detiene el flujo.

- Llama a DeleteDiscWithPath():

Elimina el disco del sistema de monitoreo o almacenamiento temporal.

##### Fdisk
```` go
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
````

Dicho método es el encargado de administrar las particiones dentro de un disco binario previamente creado. Permite crear particiones primarias, extendidas y lógicas dentro del disco, respetando las reglas de la teoría de particiones (máximo 4 primarias/extendidas, solo una extendida, y lógicas dentro de la extendida).

---

### Parámetros

- `size int`: Tamaño de la partición a crear. **Obligatorio**. Debe ser mayor a 0.  
- `path string`: Ruta del archivo de disco sobre el cual se agregará la partición. **Obligatorio**.  
- `name string`: Nombre de la partición. **Obligatorio**. No puede repetirse en el disco.  
- `unit string`: Unidad del tamaño. Puede ser `"b"` (bytes), `"k"` (kilobytes) o `"m"` (megabytes). **Opcional**, por defecto `"k"`.  
- `type string`: Tipo de partición. Puede ser `"p"` (primaria), `"e"` (extendida), o `"l"` (lógica). **Opcional**, por defecto `"p"`.  
- `fit string`: Tipo de ajuste de espacio. Puede ser `"b"` (best fit), `"f"` (first fit) o `"w"` (worst fit). **Opcional**, por defecto `"w"`.  
- `buffer *bytes.Buffer`: Buffer de log para registrar mensajes (útil para consola o archivos).

---

### Funcionamiento

#### Validación de parámetros

Los parámetros se validan y se les asignan valores por defecto desde la función `fn_fdisk`. Esta también se encarga de limpiar la entrada, verificar errores y preparar la ejecución del comando `Fdisk`.

---

#### Conversión del tamaño según unidad

- `"k"`: se multiplica por `1024`.  
- `"m"`: se multiplica por `1024 * 1024`.  
- `"b"`: no se transforma (ya está en bytes).

---

#### Apertura y lectura del disco

- Se abre el archivo con `Utilities.OpenFile(path)`  
- Se lee el MBR (Master Boot Record) desde el inicio del archivo con `Utilities.ReadObject`.

---

#### Verificación de restricciones

- El nombre de la partición no debe estar en uso.
- Se verifica la cantidad de particiones primarias y extendidas existentes (**máximo 4**).
- Solo se permite **una** partición extendida.
- No se puede crear una partición lógica si no existe una extendida.
- Se verifica el **espacio disponible** en el disco para crear la nueva partición.

---

#### Creación de particiones

- Para **primarias y extendidas**, se utiliza el arreglo de 4 particiones del MBR.
- Para una partición **extendida**, se inicializa un primer **EBR vacío** al inicio de su espacio.
- Para **particiones lógicas**, se navega por la lista enlazada de EBRs dentro de la extendida hasta encontrar el final, y se inserta un nuevo EBR con la información de la partición lógica.

---

#### Actualización del disco

- El MBR modificado se reescribe en el inicio del disco.
- En caso de lógica, también se escriben los EBRs actualizados en las posiciones correspondientes.
- Se imprimen los datos del MBR o EBRs en consola con `Structs.PrintMBR()` y `Structs.PrintEBR()` para validación visual.

---
## Comando `Mount`

```` go
func Mount(path string, name string, buffer *bytes.Buffer) {
	fmt.Fprintf(buffer, "=========MOUNT=========\n")
	fmt.Print(path)

	// Validar la ruta (path)
	if path == "" {
		fmt.Fprintf(buffer, "Error MOUNT: La ruta del disco es obligatoria.\n")
		return
	}
	// Validar el nombre (name)
	if name == "" {
		fmt.Fprintf(buffer, "Error MOUNT: El nombre de la partición es obligatorio.\n")
		return
	}

	// Abrir archivo binario
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
	var IndiceParticion int
	NameBytes := [16]byte{}
	copy(NameBytes[:], []byte(name))

	// Rechazar partición extendida
	for i := 0; i < 4; i++ {
		if TempMBR.MbrPartitions[i].Type[0] == 'e' && bytes.Equal(TempMBR.MbrPartitions[i].Name[:], NameBytes[:]) {
			fmt.Fprintf(buffer, "Error MOUNT: No se puede montar una partición extendida.\n")
			return
		}
	}

	// Buscar partición primaria con ese nombre
	for i := 0; i < 4; i++ {
		if TempMBR.MbrPartitions[i].Type[0] == 'p' && bytes.Equal(TempMBR.MbrPartitions[i].Name[:], NameBytes[:]) {
			if TempMBR.MbrPartitions[i].Status[0] == '1' {
				fmt.Fprintf(buffer, "Error MOUNT: La partición ya está montada.\n")
				return
			}
			IndiceParticion = i
			ParticionExiste = true
			break
		}
	}

	if !ParticionExiste {
		fmt.Fprintf(buffer, "Error MOUNT: No se encontró la partición con el nombre especificado. Solo se pueden montar particiones primarias.\n")
		return
	}

	// Generar ID
	DiscoID := GeneratorDiscID(path)
	
	// Verificar si ya está montada en memoria
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
	IDParticion := fmt.Sprintf("%s%d%c", UltimosDigitos, IndiceParticion+1, Letra)

	// Guardar partición en memoria RAM
	MountedPartitions[DiscoID] = append(MountedPartitions[DiscoID], MountedPartition{
		Path:   path,
		Name:   name,
		ID:     IDParticion,
		Status: '1',
	})

	fmt.Fprintf(buffer, "Partición montada con éxito en la ruta: %s con el nombre: %s y ID: %s.\n", path, name, IDParticion)

	fmt.Println("---------------------------------------------")
	PrintMountedPartitions(path, buffer)
	fmt.Println("---------------------------------------------")

	// Solo imprimir el MBR, no modificarlo
	var TempMRB Structs.MRB
	if err := Utilities.ReadObject(file, &TempMRB, 0, buffer); err != nil {
		return
	}
	Structs.PrintMBR(TempMRB)
	fmt.Println("---------------------------------------------")
}

````

Este método es el encargado de montar una partición primaria de un disco virtual en el sistema, identificándola mediante un ID único generado con base en el número de carné del estudiante, el número de partición, y una letra asignada por disco.

Este montaje se realiza **únicamente en memoria RAM** (es decir, en estructuras internas del programa) y no modifica el contenido físico del archivo binario `.mia`.

---

### Parámetros

- `path string`: Ruta del archivo de disco (.mia) donde se encuentra la partición. **Obligatorio**.  
- `name string`: Nombre de la partición que se desea montar. **Obligatorio**.  
- `buffer *bytes.Buffer`: Buffer de log para registrar mensajes (útil para consola o archivos).

---

### Funcionamiento

#### Validaciones

- Verifica que el archivo de disco exista.
- Verifica que el nombre de la partición sea válido y corresponda a una partición **primaria**.
- No se permite montar particiones **extendidas** ni inexistentes.
- No se permite montar una misma partición más de una vez.

---

#### Generación del ID de Montaje

Cada partición montada se identifica con un ID generado con la siguiente estructura:

## Comando `mounted`

Este comando permite visualizar todas las particiones que han sido montadas actualmente en el sistema. Las particiones se almacenan en memoria RAM mientras el programa se encuentra en ejecución, y desaparecen al cerrarse.

---

### Parámetros

Este comando **no recibe parámetros**.

---

### Funcionamiento

- Verifica si existen particiones montadas en la estructura de memoria `MountedPartitions`.
- Si no hay particiones montadas, muestra un mensaje informativo.
- Si existen particiones montadas, recorre la estructura e imprime el **ID de cada una**.
- Los IDs son únicos por partición montada y siguen el formato.

```` go
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
````

#### FileSystem
Este módulo contiene la implementación principal del sistema de archivos EXT2 simulado sobre discos virtuales. Su propósito es permitir el manejo de estructuras típicas de un sistema de archivos real (inodos, bloques, superbloque, permisos, jerarquía, etc.) dentro del entorno simulado del proyecto.

---

### Funcionalidades principales

- **`Mkfs`**: Formatea una partición previamente montada, creando la estructura base de EXT2. Inicializa el superbloque, bitmaps, inodos y bloques, y crea automáticamente la carpeta raíz y el archivo `users.txt` con el usuario `root`.
- **`CAT`**: Permite leer y mostrar el contenido de uno o varios archivos dentro del sistema de archivos, validando que el usuario tenga permisos de lectura.
- **`Mkdir`**: Crea directorios en cualquier parte del árbol de archivos, respetando permisos, jerarquía, y permitiendo el uso de la bandera `-p`.
- **`Mkfile`**: Crea archivos con contenido arbitrario, también respetando la jerarquía de carpetas y los permisos del usuario logueado.
- **Permisos**: Incluye funciones para validar si un usuario tiene permiso de lectura o escritura sobre un archivo o directorio, simulando permisos estilo Unix (UGO).
- **Navegación**: Implementa búsqueda recursiva de rutas para identificar inodos y bloques correspondientes a carpetas o archivos específicos.
- **Reportes**: Incluye funciones auxiliares específicas para la generación de reportes del sistema de archivos (`BuscarInodoPorRutaREPORTE`, `IsUserLoggedInREPORTE`, etc.).

---

### Estructuras clave usadas

- **Superbloque (`Superblock`)**: Contiene metainformación crítica como la cantidad total y libre de inodos y bloques, posición de bitmaps, inicios de áreas de datos, etc.
- **Inodos (`Inode`)**: Describen archivos o carpetas, almacenan metadatos como UID, GID, permisos, tamaño, y apuntadores a bloques.
- **Bloques (`FileBlock`, `FolderBlock`)**: Contienen los datos reales de archivos o directorios.
- **Bitmaps**: Representan el uso libre/ocupado de inodos y bloques.
- **`users.txt`**: Archivo del sistema que guarda la información de usuarios y grupos, creado automáticamente al formatear.

---

### Organización general

Este archivo trabaja directamente sobre archivos binarios `.mia`, manipulando sus bytes para simular las operaciones típicas de un sistema de archivos real. Todas las operaciones respetan la sesión del usuario logueado (almacenada en `User.Data`) y el ID de la partición montada.

---

### Dependencias

Este módulo depende de:
- `Utilities`: para abrir, leer y escribir objetos binarios.
- `Structs`: para definir todas las estructuras usadas en el sistema de archivos.
- `DiskManagement`: para acceder a las particiones montadas y validar sesión.
- `User`: para acceder al usuario logueado actual.

---

## Comando MKFS

```` go
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
````

El comando `mkfs` es responsable de **formatear completamente una partición montada** bajo el sistema de archivos EXT2 simulado. Su ejecución inicializa todas las estructuras esenciales del sistema de archivos, y crea automáticamente el archivo `/users.txt` en el inodo 1, donde se almacenarán los usuarios y grupos del sistema.

---

### Parámetros

| Parámetro | Tipo       | Descripción |
|----------|------------|-------------|
| `-id`    | Obligatorio | ID generado con `mount`. Identifica la partición a formatear. |
| `-type`  | Opcional    | Tipo de formateo. Solo se acepta `full` (por defecto si se omite). |
| `-fs`    | Opcional (default: `2fs`) | Define el tipo de sistema (interno). Se fija en EXT2 en la práctica. |

---

### Funcionalidad Interna

1. **Verificación de la partición:**
   - Se obtiene la lista de particiones montadas mediante `DiskManagement.GetMountedPartitions()`.
   - Se verifica que la partición especificada esté montada y activa (`Status == '1'`).

2. **Cálculo del número de estructuras (`n`):**
   - Se calcula la cantidad máxima de estructuras (`n`) posibles dentro del tamaño disponible de la partición, considerando el espacio que ocupan el superbloque, los bitmaps, inodos y bloques.

3. **Inicialización del `Superblock`:**
   - Se crea un `Superblock` con los campos definidos en EXT2, incluyendo:
     - Cantidad de inodos y bloques (`n`, `3n`)
     - Contadores de libres
     - Fechas de montaje y desmontaje
     - Tamaños de estructuras (`inode_size`, `block_size`)
     - Posiciones de inicio de bitmaps e inicios reales de inodos y bloques

4. **Creación de estructuras base:**
   - Se invoca la función `SistemaEXT2(...)`, que:
     - Inicializa los bitmaps de inodos y bloques
     - Escribe estructuras vacías (array de 0s)
     - Reserva el inodo raíz y el `users.txt`
     - Inserta el contenido inicial `1,G,root\n1,U,root,root,123\n` en `users.txt`.

---

### Estructuras creadas por `mkfs`

- Superbloque (estructura base del sistema)
- Bitmap de inodos
- Bitmap de bloques
- Tabla de inodos
- Tabla de bloques
- Inodo raíz
- Archivo `users.txt` con sus respectivos inodo y bloque de datos

---

## Comando CAT

```` go 
func CAT(files []string, buffer *bytes.Buffer) {
	fmt.Println("Usuario activo:", User.Data.GetIDUsuario())
	fmt.Println("Partición activa:", User.Data.GetIDPartition())
	//fmt.Fprintf(buffer, "Error: No hay un usuario logueado")
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

		indexInode := buscarInodoPorRuta(filePath, file, tempSuperblock, buffer)
		if indexInode == -1 {
			fmt.Printf("Error: No se pudo encontrar el archivo %s\n", filePath)
			//Structs.PrintFileblock(Structs.FileBlock{}, buffer)
			//fmt.Fprintf(buffer, "------esto despues de que no se pudo encontrar------")
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
````

El comando `CAT` permite visualizar el contenido de uno o más archivos dentro del sistema de archivos simulado. Este comando se puede ejecutar únicamente si existe una sesión activa de usuario, y se valida que el usuario tenga permisos de lectura.

### Parámetros

| Parámetro | Tipo       | Descripción                                                                 |
|-----------|------------|-----------------------------------------------------------------------------|
| -fileN    | Obligatorio| Se permite una lista de archivos a mostrar, numerados como `-file1=...`.    |

### Reglas de uso
- Se debe estar logueado en una partición.
- El usuario debe tener permisos de lectura (`7**`, `*7*`, `**7`).
- Se permite mostrar múltiples archivos, cuyo contenido será concatenado en el mismo orden de entrada.
- Si el archivo no existe o el usuario no tiene permisos, se mostrará un error específico.

### Funcionamiento Interno

1. **Parámetros**: Se parsean dinámicamente los parámetros `-fileN` y se almacenan en orden.
2. **Sesión Activa**: Se valida si hay un usuario logueado (`isUserLoggedIn`).
3. **Permisos**: Se valida si el usuario tiene permisos de lectura con `tienePermiso()`.
4. **Montaje**: Se localiza la partición montada correspondiente al usuario.
5. **Lectura de Archivos**:
    - Se obtiene el inodo del archivo con `buscarInodoPorRuta`.
    - Se accede a los bloques correspondientes.
    - Se lee el contenido y se imprime en el buffer.
6. **Errores**:
    - Si el archivo no se encuentra, se imprime un mensaje.
    - Si no hay permisos, se muestra una advertencia.


#### Structs
Este módulo define las **estructuras base** del sistema de archivos simulado, siendo fundamentales para representar el estado y los componentes del sistema de disco. Aquí se encuentran todas las estructuras utilizadas en operaciones como `mkdisk`, `fdisk`, `mkfs`, `mkdir`, `mkfile`, `cat`, y en los distintos reportes.

---

### Contenido del módulo

Incluye las definiciones y funciones para imprimir las siguientes estructuras:

- **`MRB`**: Master Boot Record (MBR), contiene metadatos del disco y las 4 particiones principales.
- **`Partition`**: Estructura que representa una partición primaria o extendida del disco.
- **`EBR`**: Extended Boot Record, usado para manejar particiones lógicas dentro de una extendida.
- **`Superblock`**: Estructura principal del sistema de archivos EXT2. Contiene punteros, contadores y parámetros esenciales.
- **`Inode`**: Representa un archivo o directorio en el sistema, con permisos, UID/GID, bloques asignados, timestamps, etc.
- **`FolderBlock` y `Content`**: Bloque que almacena hasta 4 entradas de carpeta, cada una con un nombre y apuntador a inodo.
- **`FileBlock`**: Bloque que almacena el contenido de un archivo regular (hasta 64 bytes por bloque).
- **`PointerBlock`**: Bloque especial que almacena punteros a otros bloques (directos, indirectos simples, dobles o triples).

---

### Funcionalidades auxiliares

Cada estructura tiene asociada una función `Print` que facilita la depuración o generación de reportes al imprimir de forma clara sus atributos:
- `PrintMBR`, `PrintMBRP`, `PrintPartition`, `PrintEBR`
- `PrintSuperblock`, `PrintInode`, `PrintFolderblock`, `PrintFileblock`, `PrintPointerblock`

Estas funciones se utilizan ampliamente en las salidas por consola o en la generación de reportes visuales como `.dot`, `.txt`, `.png`.

---

## Mkfile

```` go
func Mkfile(path string, p bool, content string, buffer *bytes.Buffer) {
	fmt.Fprintln(buffer, "=========== MKFILE ===========")

	if !isUserLoggedIn() {
		fmt.Fprintln(buffer, "Error MKFILE: No hay una sesión activa.")
		return
	}

	partitions := DiskManagement.GetMountedPartitions()
	var filepath, id string
	found := false
	for _, parts := range partitions {
		for _, part := range parts {
			if part.LoggedIn {
				filepath = part.Path
				id = part.ID
				found = true
				break
			}
		}
		if found {
			break
		}
	}

	if !found {
		fmt.Fprintln(buffer, "Error MKFILE: No se encontró partición activa.")
		return
	}

	file, err := Utilities.OpenFile(filepath, buffer)
	if err != nil {
		fmt.Fprintln(buffer, "Error MKFILE: No se pudo abrir el archivo.")
		return
	}
	defer file.Close()

	var mbr Structs.MRB
	if err := Utilities.ReadObject(file, &mbr, 0, buffer); err != nil {
		fmt.Fprintln(buffer, "Error MKFILE: No se pudo leer el MBR.")
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
		fmt.Fprintln(buffer, "Error MKFILE: Partición no montada.")
		return
	}

	var sb Structs.Superblock
	if err := Utilities.ReadObject(file, &sb, int64(mbr.MbrPartitions[index].Start), buffer); err != nil {
		fmt.Fprintln(buffer, "Error MKFILE: No se pudo leer el Superblock.")
		return
	}

	pathParts := strings.Split(strings.Trim(path, "/"), "/")
	fileName := pathParts[len(pathParts)-1]
	parentDirs := pathParts[:len(pathParts)-1]

	rootInode := Structs.Inode{}
	if err := Utilities.ReadObject(file, &rootInode, int64(sb.SB_Inode_Start), buffer); err != nil {
		fmt.Fprintln(buffer, "Error MKFILE: No se pudo leer el Inodo raíz.")
		return
	}

	err = crearCarpetas(parentDirs, rootInode, file, sb, p, buffer)
	if err != nil {
		fmt.Fprintln(buffer, "Error MKFILE:", err)
		return
	}

	current := Structs.Inode{}
	entryInodeIndex := int32(0)
	if err := Utilities.ReadObject(file, &current, int64(sb.SB_Inode_Start), buffer); err != nil {
		fmt.Fprintln(buffer, "Error MKFILE: No se pudo leer el inodo raíz.")
		return
	}
	for _, part := range parentDirs {
		found := false
		for _, blk := range current.IN_Block {
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
				if name == part {
					entryInodeIndex = entry.B_Inode
					offset := int64(sb.SB_Inode_Start + entry.B_Inode*int32(binary.Size(Structs.Inode{})))
					if err := Utilities.ReadObject(file, &current, offset, buffer); err == nil {
						found = true
						break
					}
				}
			}
			if found {
				break
			}
		}
		if !found {
			fmt.Fprintln(buffer, "Error MKFILE: Carpeta padre no encontrada")
			return
		}
	}

	if !tienePermisoEscritura(current, User.Data.GetUID(), User.Data.GetGID()) {
		fmt.Fprintln(buffer, "Error MKFILE: El usuario no tiene permiso de escritura en la carpeta padre.")
		return
	}

	for _, blk := range current.IN_Block {
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
			if name == fileName {
				fmt.Fprintf(buffer, "Error MKFILE: Ya existe un archivo con el nombre '%s'\n", fileName)
				return
			}
		}
	}

	inodeIdx, blockIdx, err := reservarInodoYBloque(file, sb, buffer)
	if err != nil {
		fmt.Fprintln(buffer, "Error MKFILE:", err)
		return
	}

	var newInode Structs.Inode
	newInode.IN_Uid = int32(User.Data.GetUID())
	newInode.IN_Gid = int32(User.Data.GetGID())
	newInode.IN_Size = int32(len(content))
	copy(newInode.IN_Perm[:], "664")
	now := time.Now().Format("2006-01-02 15:04:05")
	copy(newInode.IN_Atime[:], now)
	copy(newInode.IN_Ctime[:], now)
	copy(newInode.IN_Mtime[:], now)
	for i := 0; i < 15; i++ {
		newInode.IN_Block[i] = -1
	}
	newInode.IN_Block[0] = blockIdx

	offsetInode := int64(sb.SB_Inode_Start + inodeIdx*int32(binary.Size(Structs.Inode{})))
	if err := Utilities.WriteObject(file, newInode, offsetInode, buffer); err != nil {
		fmt.Fprintln(buffer, "Error MKFILE: No se pudo escribir el inodo del archivo")
		return
	}

	var fileblock Structs.FileBlock
	copy(fileblock.B_Content[:], content)
	offsetBlock := int64(sb.SB_Block_Start + blockIdx*int32(binary.Size(Structs.FileBlock{})))
	if err := Utilities.WriteObject(file, fileblock, offsetBlock, buffer); err != nil {
		fmt.Fprintln(buffer, "Error MKFILE: No se pudo escribir el contenido del archivo")
		return
	}

	if err := agregarEntradaAFolderBlock(&current, fileName, inodeIdx, file, sb, buffer); err != nil {
		fmt.Fprintln(buffer, "Error MKFILE: No se pudo agregar la entrada al folder padre")
		return
	}

	offsetCurrent := int64(sb.SB_Inode_Start + entryInodeIndex*int32(binary.Size(Structs.Inode{})))
	if err := Utilities.WriteObject(file, current, offsetCurrent, buffer); err != nil {
		fmt.Fprintln(buffer, "Error MKFILE: No se pudo escribir el inodo padre actualizado")
		return
	}

	fmt.Fprintln(buffer, "Archivo creado exitosamente:", path)


	var createdInode Structs.Inode
	Utilities.ReadObject(file, &createdInode, int64(sb.SB_Inode_Start+inodeIdx*int32(binary.Size(Structs.Inode{}))), buffer)
	Structs.PrintInode(createdInode)

	var createdBlock Structs.FileBlock
	Utilities.ReadObject(file, &createdBlock, int64(sb.SB_Block_Start+blockIdx*int32(binary.Size(Structs.FileBlock{}))), buffer)
	Structs.PrintFileblock(createdBlock, buffer)

	for _, blk := range current.IN_Block {
		if blk == -1 {
			continue
		}
		var folder Structs.FolderBlock
		offset := int64(sb.SB_Block_Start + blk*int32(binary.Size(Structs.FolderBlock{})))
		if err := Utilities.ReadObject(file, &folder, offset, buffer); err != nil {
			continue
		}
		Structs.PrintFolderblock(folder)
	}
}
````

Este comando permite crear un archivo dentro del sistema de archivos EXT2 simulado. El propietario del archivo será el usuario actualmente logueado, y por defecto se asignan permisos `664`. Si el usuario en sesión es `root`, se considera que tiene permisos `777` sobre todos los archivos, ignorando restricciones.

### Parámetros:

| Parámetro | Tipo       | Descripción |
|----------|------------|-------------|
| `-path`  | Obligatorio | Ruta absoluta donde se creará el archivo. Si contiene espacios, debe ir entre comillas. |
| `-r`     | Opcional   | Si se especifica, se crean carpetas padre que no existan. No debe llevar ningún valor. |
| `-size`  | Opcional   | Indica el tamaño en bytes del archivo. El contenido será una secuencia de `0123456789...` repetida. Se ignora si se especifica `-cont`. |
| `-cont`  | Opcional   | Ruta a un archivo en el sistema real desde el cual se cargará el contenido. Tiene prioridad sobre `-size`. |

### Consideraciones:
- Si ya existe un archivo con el mismo nombre, el sistema notifica que no puede sobrescribirlo.
- Si no existen las carpetas padres, se mostrará un error, a menos que se especifique `-r`.
- El contenido se guarda en bloques asignados dinámicamente desde el bitmap.
- Los permisos UGO del archivo se definen como `664`, a menos que el usuario sea `root`.

## Mkdir

```` go
func Mkdir(path string, p bool, buffer *bytes.Buffer) {
	fmt.Fprintln(buffer, "=========== MKDIR ===========")

	// 1. Validar sesión activa
	if !isUserLoggedIn() {
		fmt.Fprintln(buffer, "Error MKDIR: No hay una sesión activa.")
		return
	}

	// 2. Obtener partición montada con sesión activa
	partitions := DiskManagement.GetMountedPartitions()
	var filepath, id string
	found := false
	for _, parts := range partitions {
		for _, part := range parts {
			if part.LoggedIn {
				filepath = part.Path
				id = part.ID
				found = true
				break
			}
		}
		if found {
			break
		}
	}
	if !found {
		fmt.Fprintln(buffer, "Error MKDIR: No se encontró partición activa.")
		return
	}

	// 3. Abrir disco
	file, err := Utilities.OpenFile(filepath, buffer)
	if err != nil {
		fmt.Fprintln(buffer, "Error MKDIR: No se pudo abrir el archivo.")
		return
	}
	defer file.Close()

	// 4. Leer MBR y Superblock
	var mbr Structs.MRB
	if err := Utilities.ReadObject(file, &mbr, 0, buffer); err != nil {
		fmt.Fprintln(buffer, "Error MKDIR: No se pudo leer el MBR.")
		return
	}

	// Buscar la partición correspondiente
	var index int = -1
	for i := 0; i < 4; i++ {
		if strings.Contains(string(mbr.MbrPartitions[i].ID[:]), id) && mbr.MbrPartitions[i].Status[0] == '1' {
			index = i
			break
		}
	}
	if index == -1 {
		fmt.Fprintln(buffer, "Error MKDIR: Partición no montada.")
		return
	}

	var sb Structs.Superblock
	if err := Utilities.ReadObject(file, &sb, int64(mbr.MbrPartitions[index].Start), buffer); err != nil {
		fmt.Fprintln(buffer, "Error MKDIR: No se pudo leer el Superblock.")
		return
	}

	// 5. Procesar ruta
	pathParts := strings.Split(strings.Trim(path, "/"), "/")
	if len(pathParts) == 0 {
		fmt.Fprintln(buffer, "Error MKDIR: Ruta inválida.")
		return
	}

	// 6. Buscar carpeta padre y crear carpetas si no existen
	rootInode := Structs.Inode{}
	if err := Utilities.ReadObject(file, &rootInode, int64(sb.SB_Inode_Start), buffer); err != nil {
		fmt.Fprintln(buffer, "Error MKDIR: No se pudo leer el Inodo raíz.")
		return
	}

	err = crearCarpetas(pathParts, rootInode, file, sb, p, buffer)
	if err != nil {
		fmt.Fprintln(buffer, "Error MKDIR:", err)
		return
	}

	fmt.Fprintln(buffer, "Carpeta creada exitosamente:", path)
}
````

Este comando permite crear carpetas dentro del sistema de archivos EXT2 simulado. El propietario de la carpeta será el usuario que ha iniciado sesión. Las carpetas nuevas se crean con permisos predeterminados `664`.

### Parámetros:

| Parámetro | Tipo       | Descripción |
|----------|------------|-------------|
| `-path`  | Obligatorio | Ruta absoluta donde se creará la carpeta. Si contiene espacios debe ir entre comillas. |
| `-p`     | Opcional   | Si se especifica, se crean las carpetas padres que no existan. No debe llevar ningún valor. |

### Consideraciones:
- Si el usuario no tiene permiso de escritura sobre la carpeta padre, se mostrará un error.
- Si la carpeta ya existe, no se sobrescribe ni lanza error.
- Si no existen carpetas padre y no se usa `-p`, se mostrará un error.
- El usuario `root` siempre tendrá permiso para crear carpetas, sin importar permisos.


#### User

Este módulo gestiona todo lo relacionado con **usuarios, grupos y sesiones activas** en el sistema de archivos simulado. Implementa la lógica necesaria para operar sobre el archivo especial `/users.txt`, donde se almacenan de forma persistente los registros de usuarios y grupos del sistema.

---

### Contenido del módulo

#### Control de sesión:
- `PartitionUser`: estructura global `Data` que mantiene información del usuario actualmente logueado (usuario, partición, UID, GID).
- `Login(...)`: Verifica credenciales y establece una sesión activa.
- `LogOut(...)`: Finaliza la sesión activa actual.
- Métodos `Get` y `Set` para obtener o actualizar los campos de la sesión.

#### Gestión de Grupos y Usuarios:
- `Mkgrp(...)`: Crea un nuevo grupo (solo lo puede ejecutar `root`).
- `Mkusr(...)`: Crea un nuevo usuario dentro de un grupo (solo `root`).
- `Rmusr(...)`: Elimina lógicamente un usuario (marca el ID con `0`).
- `Rmgrp(...)`: Elimina lógicamente un grupo.
- `Chgrp(...)`: Cambia de grupo a un usuario existente.

#### Manejo de `/users.txt`:
- `InitSearch(...)` y `SarchInodeByPath(...)`: Navegan por la estructura de carpetas para encontrar el inodo de `/users.txt`.
- `GetInodeFileData(...)`: Recupera todo el contenido del archivo `users.txt` desde los bloques del inodo.
- `AppendToFileBlock(...)`: Añade contenido nuevo (grupo o usuario) al archivo `users.txt`, reutilizando bloques existentes o solicitando nuevos si es necesario.
- `OverwriteFileBlock(...)`: Borra bloques actuales y sobrescribe por completo el contenido de `users.txt`.

#### Utilidades:
- `obtenerIDGrupo(...)`: Extrae el ID numérico de un grupo según su nombre.
- `pop(...)`: Extrae el último elemento de una lista de strings (tipo stack).

## Comando `LOGIN`

El comando `LOGIN` permite iniciar sesión en el sistema de archivos. Este comando es obligatorio para ejecutar la mayoría de los comandos que manipulan archivos, carpetas y usuarios, exceptuando `MKFS` y el mismo `LOGIN`.

```` Go
func Login(user string, pass string, id string, buffer *bytes.Buffer) {
	fmt.Println("======Start LOGIN======")
	fmt.Println("User:", user)
	fmt.Println("Pass:", pass)
	fmt.Println("Id:", id)

	mountedPartitions := DiskManagement.GetMountedPartitions()
	var filepath string
	var partitionFound bool
	var login bool = false

	for _, partitions := range mountedPartitions {
		for _, Partition := range partitions {
			if Partition.ID == id && Partition.LoggedIn {
				fmt.Fprintf(buffer, "Error LOGIN: Ya existe un usuario logueado en la partición:%s\n", id)
				return
			}
			if Partition.ID == id {
				filepath = Partition.Path
				partitionFound = true
				break
			}
		}
		if partitionFound {
			break
		}
	}

	if !partitionFound {
		fmt.Fprintf(buffer, "Error en LOGIN: No se encontró ninguna partición montada con el ID proporcionado %s\n", id)
		return
	}

	file, err := Utilities.OpenFile(filepath, buffer)
	if err != nil {
		fmt.Println("Error: No se pudo abrir el archivo:", err)
		return
	}
	defer file.Close()

	var TempMBR Structs.MRB
	if err := Utilities.ReadObject(file, &TempMBR, 0, buffer); err != nil {
		fmt.Println("Error: No se pudo leer el MBR:", err)
		return
	}

	var index int = -1
	for i := 0; i < 4; i++ {
		if TempMBR.MbrPartitions[i].Size != 0 {
			if strings.Contains(string(TempMBR.MbrPartitions[i].ID[:]), id) {
				if TempMBR.MbrPartitions[i].Status[0] == '1' {
					fmt.Println("particion montada\n")
					index = i
				} else {
					fmt.Println("particion no montada\n")
					return
				}
				break
			}
		}
	}

	if index == -1 {
		fmt.Fprintf(buffer, "Error en LOGIN: no se encontró ninguna partición con el ID %s\n", id)
		return
	}

	var tempSuperblock Structs.Superblock
	if err := Utilities.ReadObject(file, &tempSuperblock, int64(TempMBR.MbrPartitions[index].Start), buffer); err != nil {
		fmt.Println("Error: No se pudo leer el Superblock:", err)
		return
	}

	indexInode := InitSearch("/users.txt", file, tempSuperblock, buffer)
	if indexInode == -1 {
		fmt.Fprintln(buffer, "Error: No se encontró /users.txt, asegúrese de haber ejecutado mkfs.")
		return
	}

	var crrInode Structs.Inode
	if err := Utilities.ReadObject(file, &crrInode, int64(tempSuperblock.SB_Inode_Start+indexInode*int32(binary.Size(Structs.Inode{}))), buffer); err != nil {
		fmt.Fprintf(buffer, "Error: No se pudo leer el Inodo: %v\n", err)
		return
	}

	data := GetInodeFileData(crrInode, file, tempSuperblock, buffer)
	lines := strings.Split(data, "\n")

	var uidEncontrado, gidEncontrado int

	for _, line := range lines {
		words := strings.Split(line, ",")
		if len(words) == 5 {
			if words[3] == user && words[4] == pass {
				login = true
				uidEncontrado, _ = strconv.Atoi(words[0])          // UID del usuario
				gidEncontrado = obtenerIDGrupo(words[2], lines)    // GID según nombre grupo
				break
			}
		}
	}

	if login {
		fmt.Fprintf(buffer, "Usuario logueado con éxito en la partición:%s\n", id)
		fmt.Println("Usuario logueado con éxito")
		DiskManagement.MarkPartitionAsLoggedIn(id)
		Data.SetIDPartition(id)
		Data.SetIDUsuario(user)
		Data.SetUID(uidEncontrado)
		Data.SetGID(gidEncontrado)
	} else {
		fmt.Fprintf(buffer, "Error LOGIN: Credenciales incorrectas o usuario no encontrado.\n")
	}

	fmt.Println("======End LOGIN======")
}
````

1. **Inicio de depuración**
   - Se imprimen por consola los parámetros recibidos: usuario, contraseña e ID de partición.

2. **Búsqueda de partición montada**
   - Se recorre la lista de particiones montadas usando `DiskManagement.GetMountedPartitions()`.
   - Si se encuentra una partición con el ID especificado que **ya tiene una sesión activa (`LoggedIn`)**, se lanza un error y se aborta.
   - Si se encuentra la partición pero no tiene sesión activa, se guarda su ruta (`Path`) para usarla más adelante.

3. **Validación de existencia**
   - Si no se encontró ninguna partición montada con el ID especificado, se imprime un mensaje de error.

4. **Apertura del archivo del disco**
   - Se abre el archivo `.mia` correspondiente a la partición.

5. **Lectura del MBR**
   - Se lee el Master Boot Record (MBR) desde el inicio del archivo binario.

6. **Verificación de estado de la partición**
   - Se busca dentro de las 4 particiones del MBR una que tenga el ID coincidente y esté marcada como activa (`Status == '1'`).
   - Si no se encuentra una partición válida, se imprime error.

7. **Lectura del Superblock**
   - Se accede al `Superblock` a partir del byte de inicio de la partición y se carga en memoria para acceder a las estructuras del sistema EXT2.

8. **Búsqueda del archivo `users.txt`**
   - Se utiliza `InitSearch("/users.txt")` para obtener el índice del inodo correspondiente a ese archivo lógico.
   - Si el archivo no existe, probablemente no se ha ejecutado el comando `mkfs`, y se lanza un mensaje de error.

9. **Lectura del inodo del archivo `users.txt`**
   - Se accede al inodo para luego extraer el contenido del archivo usando `GetInodeFileData`.

10. **Validación de credenciales**
    - Se recorre cada línea del archivo `users.txt` en busca de una línea con el formato `UID,U,Grupo,Usuario,Contraseña`.
    - Si hay coincidencia exacta de `user` y `pass`, se considera login exitoso.
    - Además, se extraen el UID y el GID asociados para almacenarlos en la sesión.

11. **Inicio de sesión**
    - Si las credenciales fueron válidas:
      - Se marca la partición como activa (`MarkPartitionAsLoggedIn`).
      - Se actualiza la variable global `User.Data` con `IDPartition`, `IDUsuario`, `UID` y `GID`.
      - Se imprime un mensaje indicando que el login fue exitoso.
    - Si las credenciales no coinciden:
      - Se informa al usuario que la autenticación falló.

---
### Comando `LOGOUT`

```` go
func LogOut(buffer *bytes.Buffer) {
	fmt.Fprint(buffer, "==========LOGOUT==========\n")
	mountedPartitions := DiskManagement.GetMountedPartitions()
	var SesionActiva bool

	if len(mountedPartitions) == 0 {
		fmt.Fprintf(buffer, "Error LOGOUT: No hay ninguna partición montada.\n")
		return
	}

	for _, Particiones := range mountedPartitions {
		for _, Particion := range Particiones {
			if Particion.LoggedIn {
				SesionActiva = true
				break
			}
		}
		if SesionActiva {
			break
		}
	}
	if !SesionActiva {
		fmt.Fprintf(buffer, "Error LOGOUT: No hay ninguna sesión activa.\n")
		return
	} else {
		DiskManagement.MarkPartitionAsLoggedOut(Data.GetIDPartition())
		fmt.Fprintf(buffer, "Sesión cerrada con éxito de la partición:%s\n", Data.GetIDPartition())
	}
	Data.SetIDPartition("")
	Data.SetIDUsuario("")
}
````

Este comando se encarga de cerrar la sesión actualmente activa en el sistema de archivos. Es obligatorio tener una sesión activa para poder ejecutar este comando; de lo contrario, debe mostrar un mensaje de error.

---

#### Descripción del funcionamiento del método `LogOut(buffer *bytes.Buffer)`

1. **Impresión del encabezado**
   - Se imprime un mensaje inicial: `==========LOGOUT==========`.

2. **Verificación de particiones montadas**
   - Se obtiene la lista de particiones montadas mediante `DiskManagement.GetMountedPartitions()`.
   - Si no hay particiones montadas, se lanza el error:  
     `Error LOGOUT: No hay ninguna partición montada.`

3. **Verificación de sesión activa**
   - Se recorre cada partición montada y se verifica si alguna tiene la bandera `LoggedIn = true`.
   - Si no se encuentra una partición con sesión activa, se muestra:  
     `Error LOGOUT: No hay ninguna sesión activa.`

4. **Cierre de sesión**
   - Se invoca `DiskManagement.MarkPartitionAsLoggedOut(id)` con el ID actual de la sesión.
   - Se imprime un mensaje de confirmación como:  
     `Sesión cerrada con éxito de la partición:062A`
   - Se limpian las variables globales de sesión `User.Data`:
     - `SetIDPartition("")`
     - `SetIDUsuario("")`

---

#### Validación de parámetros

La función `fn_logout` valida que **no se reciban parámetros adicionales**. Si se encuentra algo como `logout extra`, se devuelve:
#### Consideraciones importantes

- Solo puede haber **una sesión activa por partición**.
- El archivo `users.txt` debe haber sido creado previamente por `mkfs` y contener al menos al usuario `root` con su contraseña.
- Se respetan **mayúsculas y minúsculas** en `user` y `pass`.
- No se permiten múltiples sesiones concurrentes.

---

### Comando `MKGRP`

```` go 
func Mkgrp(name string, buffer *bytes.Buffer) {
	fmt.Fprint(buffer, "=============MKGRP=============\n")

	mountedPartitions := DiskManagement.GetMountedPartitions()
	var filePath string
	var PartitionFound bool

	// Buscar la partición donde se ha iniciado sesión
	for _, Particiones := range mountedPartitions {
		for _, Particion := range Particiones {
			if Particion.ID == Data.GetIDPartition() {
				filePath = Particion.Path
				PartitionFound = true
				break
			}
		}
		if PartitionFound {
			break
		}
	}

	if !PartitionFound {
		fmt.Fprintf(buffer, "Error MKGRP: No se encontró ninguna partición montada con el ID: %s\n", Data.GetIDPartition())
		return
	}

	// Validación del usuario "root"
	if Data.GetIDUsuario() != "root" {
		fmt.Fprintf(buffer, "Error MKGRP: Solo el usuario 'root' puede crear grupos.\n")
		return
	}

	// Abrir el archivo de la partición
	file, err := Utilities.OpenFile(filePath, buffer)
	if err != nil {
		fmt.Fprintf(buffer, "Error MKGRP: no se pudo abrir disco: %v\n", err)
		return
	}
	defer file.Close()

	// Leer el MBR de la partición
	var TempMBR Structs.MRB
	if err := Utilities.ReadObject(file, &TempMBR, 0, buffer); err != nil {
		fmt.Fprintf(buffer, "Error MKGRP: No se pudo leer el MBR: %v\n", err)
		return
	}

	var partitionIndex int = -1
	for i, part := range TempMBR.MbrPartitions {
		if strings.Trim(string(part.ID[:]), "\x00") == Data.GetIDPartition() && part.Status[0] == '1' {
			partitionIndex = i
			break
		}
	}
	if partitionIndex == -1 {
		fmt.Fprintf(buffer, "Error MKGRP: La partición no está montada o no existe.\n")
		return
	}

	var tempSuperblock Structs.Superblock
	sbStart := int64(TempMBR.MbrPartitions[partitionIndex].Start)
	if err := Utilities.ReadObject(file, &tempSuperblock, sbStart, buffer); err != nil {
		fmt.Fprintf(buffer, "Error MKGRP: No se pudo leer el Superblock: %v\n", err)
		return
	}

	// Buscar archivo /users.txt (obtienes índice de Inodo)
	indexInode := InitSearch("/users.txt", file, tempSuperblock, buffer)
	if indexInode == -1 {
		fmt.Fprintf(buffer, "Error MKGRP: No se encontró el archivo /users.txt\n")
		return
	}

	var usersInode Structs.Inode
	inodePos := int64(tempSuperblock.SB_Inode_Start + indexInode*int32(binary.Size(Structs.Inode{})))
	if err := Utilities.ReadObject(file, &usersInode, inodePos, buffer); err != nil {
		fmt.Fprintf(buffer, "Error MKGRP: No se pudo leer el Inodo: %v\n", err)
		return
	}

	data := GetInodeFileData(usersInode, file, tempSuperblock, buffer)
	lines := strings.Split(data, "\n")

	var groupID int = 1
	for _, line := range lines {
		if strings.TrimSpace(line) == "" {
			continue
		}
		fields := strings.Split(line, ",")
		if len(fields) == 3 && fields[1] == "G" {
			if fields[2] == name {
				fmt.Fprintf(buffer, "Error MKGRP: El grupo '%s' ya existe.\n", name)
				return
			}
			groupID++
		}
	}

	// Añadir el nuevo grupo al archivo claramente
	newGroup := fmt.Sprintf("%d,G,%s\n", groupID, name)

	err = AppendToFileBlock(&usersInode, indexInode, newGroup, file, tempSuperblock, buffer)
	if err != nil {
		fmt.Fprintf(buffer, "Error MKGRP: Al escribir nuevo grupo: %v\n", err)
		return
	}

	fmt.Fprintf(buffer, "Grupo '%s' creado exitosamente.\n", name)
}
````

Este comando permite crear un nuevo grupo de usuarios en el sistema de archivos simulado. El grupo se almacena en el archivo lógico `/users.txt`, ubicado en la raíz del sistema de archivos (`/`). Solo el usuario `root` tiene permisos para ejecutar esta instrucción.

---

- Solo el usuario `root` puede crear grupos.
- El nombre del grupo es sensible a mayúsculas y minúsculas.
- Si el grupo ya existe, se debe mostrar un error.
- Se agrega una nueva línea en el archivo `users.txt` con la forma:  
  `ID,G,nombreGrupo`
- El archivo `users.txt` debe estar presente en el sistema EXT2.

Tu implementación **cumple con todos estos puntos.**

---

#### Descripción de la implementación

##### Método principal: `Mkgrp(name string, buffer *bytes.Buffer)`

1. **Obtiene la partición con sesión activa**
   - Recorre las particiones montadas y verifica cuál está activa (logueada) para el usuario actual.

2. **Verifica si el usuario actual es `root`**
   - Si no lo es, devuelve un error:
     ```
     Error MKGRP: Solo el usuario 'root' puede crear grupos.
     ```

3. **Abre el archivo del disco**
   - Se abre el archivo binario que representa la partición activa.

4. **Lee el MBR y localiza la partición por su ID**
   - Verifica que la partición está montada y tiene el mismo ID asociado a la sesión.

5. **Lee el Superblock y encuentra el inodo de `/users.txt`**
   - Usa `InitSearch("/users.txt")` para encontrar el índice del inodo que contiene la lista de usuarios y grupos.

6. **Verifica si el grupo ya existe**
   - Si existe una línea con el tipo `"G"` y el mismo nombre, lanza un error:
     ```
     Error MKGRP: El grupo 'usuarios' ya existe.
     ```

7. **Calcula el nuevo ID del grupo**
   - Se basa en la cantidad de grupos actuales en `users.txt` (se incrementa por cada grupo válido encontrado).

8. **Agrega la entrada del grupo al final del archivo**
   - Usa `AppendToFileBlock(...)` para agregar el nuevo grupo al contenido del archivo.

9. **Confirma al usuario**
   - Mensaje:
     ```
     Grupo 'usuarios' creado exitosamente.
     ```

---

### Comando `RMGRP`

````go
func Rmgrp(name string, buffer *bytes.Buffer) {
	fmt.Fprint(buffer, "=============RMGRP=============\n")

	// Validar que el usuario sea root
	if Data.GetIDUsuario() != "root" {
		fmt.Fprintf(buffer, "Error RMGRP: Solo el usuario 'root' puede eliminar grupos.\n")
		return
	}

	// Buscar partición montada
	mountedPartitions := DiskManagement.GetMountedPartitions()
	var filePath string
	var partitionFound bool
	for _, Particiones := range mountedPartitions {
		for _, Particion := range Particiones {
			if Particion.ID == Data.GetIDPartition() {
				filePath = Particion.Path
				partitionFound = true
				break
			}
		}
		if partitionFound {
			break
		}
	}
	if !partitionFound {
		fmt.Fprintf(buffer, "Error RMGRP: No se encontró ninguna partición montada con el ID: %s\n", Data.GetIDPartition())
		return
	}

	// Abrir archivo del disco
	file, err := Utilities.OpenFile(filePath, buffer)
	if err != nil {
		fmt.Fprintf(buffer, "Error RMGRP: No se pudo abrir el disco: %v\n", err)
		return
	}
	defer file.Close()

	// Leer MBR y Superblock
	var TempMBR Structs.MRB
	if err := Utilities.ReadObject(file, &TempMBR, 0, buffer); err != nil {
		fmt.Fprintf(buffer, "Error RMGRP: No se pudo leer el MBR: %v\n", err)
		return
	}

	var partitionIndex int = -1
	for i, part := range TempMBR.MbrPartitions {
		if strings.Trim(string(part.ID[:]), "\x00") == Data.GetIDPartition() && part.Status[0] == '1' {
			partitionIndex = i
			break
		}
	}
	if partitionIndex == -1 {
		fmt.Fprintf(buffer, "Error RMGRP: La partición no está montada o no existe.\n")
		return
	}

	var tempSuperblock Structs.Superblock
	sbStart := int64(TempMBR.MbrPartitions[partitionIndex].Start)
	if err := Utilities.ReadObject(file, &tempSuperblock, sbStart, buffer); err != nil {
		fmt.Fprintf(buffer, "Error RMGRP: No se pudo leer el Superblock: %v\n", err)
		return
	}

	// Buscar /users.txt
	indexInode := InitSearch("/users.txt", file, tempSuperblock, buffer)
	if indexInode == -1 {
		fmt.Fprintf(buffer, "Error RMGRP: No se encontró el archivo /users.txt\n")
		return
	}

	var usersInode Structs.Inode
	inodePos := int64(tempSuperblock.SB_Inode_Start + indexInode*int32(binary.Size(Structs.Inode{})))
	if err := Utilities.ReadObject(file, &usersInode, inodePos, buffer); err != nil {
		fmt.Fprintf(buffer, "Error RMGRP: No se pudo leer el Inodo: %v\n", err)
		return
	}

	// Leer contenido actual
	data := GetInodeFileData(usersInode, file, tempSuperblock, buffer)
	lines := strings.Split(data, "\n")

	var updatedLines []string
	groupFound := false

	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		fields := strings.Split(line, ",")
		if len(fields) == 3 && fields[1] == "G" && fields[2] == name && fields[0] != "0" {
			fields[0] = "0" // marcar como eliminado
			groupFound = true
			updatedLines = append(updatedLines, strings.Join(fields, ","))
		} else {
			updatedLines = append(updatedLines, line)
		}
	}

	if !groupFound {
		fmt.Fprintf(buffer, "Error RMGRP: El grupo '%s' no existe o ya fue eliminado.\n", name)
		return
	}

	// Guardar el contenido nuevo
	newData := strings.Join(updatedLines, "\n") + "\n"
	err = OverwriteFileBlock(&usersInode, indexInode, newData, file, tempSuperblock, buffer)
	if err != nil {
		fmt.Fprintf(buffer, "Error RMGRP: Al escribir el archivo actualizado: %v\n", err)
		return
	}

	fmt.Fprintf(buffer, "Grupo '%s' eliminado exitosamente.\n", name)
}
````

Este comando permite eliminar lógicamente un grupo de usuarios en el sistema de archivos simulado. Solo puede ser ejecutado por el usuario `root`, y realiza la eliminación lógica (no física) en el archivo `users.txt`, colocando el identificador del grupo como `0`.

---

- Solo el usuario `root` puede ejecutar el comando.
- El grupo debe existir y no estar previamente eliminado.
- El grupo se elimina marcando su ID como `0` en el archivo `users.txt`.
- El archivo `users.txt` se actualiza correctamente.
- Si el grupo no existe o ya fue eliminado, muestra un mensaje de error.
- Distingue entre mayúsculas y minúsculas en el nombre del grupo.

---

#### Descripción de la implementación

##### Método principal: `Rmgrp(name string, buffer *bytes.Buffer)`

1. **Verificación de permisos**
   - Si el usuario actual no es `root`, muestra:
     ```
     Error RMGRP: Solo el usuario 'root' puede eliminar grupos.
     ```

2. **Localización de la partición activa**
   - Busca entre las particiones montadas la que tiene una sesión activa.

3. **Lectura del disco y estructuras**
   - Abre el archivo del disco.
   - Lee el `MBR` para encontrar la partición correcta.
   - Lee el `Superblock`.

4. **Ubicación del archivo `/users.txt`**
   - Busca el inodo del archivo `users.txt` usando `InitSearch("/users.txt")`.

5. **Lectura del contenido actual**
   - Obtiene todo el contenido actual del archivo `users.txt` línea por línea.

6. **Modificación lógica del grupo**
   - Busca una línea con:
     - Tipo `G`
     - Nombre del grupo coincidente
     - ID distinto de `0`
   - Si lo encuentra, cambia el ID a `0` y marca como eliminado.
   - Si no se encuentra, muestra:
     ```
     Error RMGRP: El grupo 'usuarios' no existe o ya fue eliminado.
     ```

7. **Escritura del archivo actualizado**
   - Reescribe completamente el archivo `users.txt` con las nuevas líneas usando `OverwriteFileBlock(...)`.


### Comando `MKUSR`

```` go
func Mkusr(user string, pass string, grp string, buffer *bytes.Buffer) {
	fmt.Fprint(buffer, "=============MKUSR=============\n")

	// Verificar que el usuario logueado sea 'root'
	if Data.GetIDUsuario() != "root" {
		fmt.Fprintf(buffer, "Error MKUSR: Solo 'root' puede crear usuarios.\n")
		return
	}

	// Validación longitud máxima
	if len(user) > 10 {
		fmt.Fprintf(buffer, "Error MKUSR: Nombre del usuario máximo 10 caracteres.\n")
		return
	}
	if len(pass) > 10 {
		fmt.Fprintf(buffer, "Error MKUSR: Contraseña máximo 10 caracteres.\n")
		return
	}
	if len(grp) > 10 {
		fmt.Fprintf(buffer, "Error MKUSR: Nombre del grupo máximo 10 caracteres.\n")
		return
	}

	// Buscar la partición montada
	mountedPartitions := DiskManagement.GetMountedPartitions()
	var filePath string
	var partitionFound bool
	for _, Particiones := range mountedPartitions {
		for _, Particion := range Particiones {
			if Particion.ID == Data.GetIDPartition() {
				filePath = Particion.Path
				partitionFound = true
				break
			}
		}
		if partitionFound {
			break
		}
	}

	if !partitionFound {
		fmt.Fprintf(buffer, "Error MKUSR: No se encontró partición montada.\n")
		return
	}

	// Abrir el archivo de la partición
	file, err := Utilities.OpenFile(filePath, buffer)
	if err != nil {
		fmt.Fprintf(buffer, "Error MKUSR: No se pudo abrir disco: %v\n", err)
		return
	}
	defer file.Close()

	// Leer el MBR
	var TempMBR Structs.MRB
	if err := Utilities.ReadObject(file, &TempMBR, 0, buffer); err != nil {
		fmt.Fprintf(buffer, "Error MKUSR: No se pudo leer el MBR: %v\n", err)
		return
	}

	// Buscar la partición correcta
	var partitionIndex int = -1
	for i, part := range TempMBR.MbrPartitions {
		if strings.Trim(string(part.ID[:]), "\x00") == Data.GetIDPartition() && part.Status[0] == '1' {
			partitionIndex = i
			break
		}
	}
	if partitionIndex == -1 {
		fmt.Fprintf(buffer, "Error MKUSR: Partición no montada o inexistente.\n")
		return
	}

	// Leer el Superblock
	var tempSuperblock Structs.Superblock
	sbStart := int64(TempMBR.MbrPartitions[partitionIndex].Start)
	if err := Utilities.ReadObject(file, &tempSuperblock, sbStart, buffer); err != nil {
		fmt.Fprintf(buffer, "Error MKUSR: No se pudo leer Superblock: %v\n", err)
		return
	}

	// Buscar archivo /users.txt
	indexInode := InitSearch("/users.txt", file, tempSuperblock, buffer)
	if indexInode == -1 {
		fmt.Fprintf(buffer, "Error MKUSR: No se encontró archivo /users.txt\n")
		return
	}

	var usersInode Structs.Inode
	inodePos := int64(tempSuperblock.SB_Inode_Start + indexInode*int32(binary.Size(Structs.Inode{})))
	if err := Utilities.ReadObject(file, &usersInode, inodePos, buffer); err != nil {
		fmt.Fprintf(buffer, "Error MKUSR: No se pudo leer el Inodo: %v\n", err)
		return
	}

	// Leer los datos del archivo /users.txt
	data := GetInodeFileData(usersInode, file, tempSuperblock, buffer)
	lines := strings.Split(data, "\n")

	var userID int = 1
	var grupoExiste bool = false
	var userFound bool = false
	var userLine string

	// Buscar si el usuario existe en cualquier grupo
	for _, line := range lines {
		if strings.TrimSpace(line) == "" {
			continue
		}
		fields := strings.Split(line, ",")
		if len(fields) == 5 && fields[1] == "U" {
			// Si el usuario ya existe y está eliminado (ID = 0), restaurarlo
			if fields[3] == user && fields[0] == "0" {
				// Restaurar el usuario, cambiar el ID a 1
				fields[0] = "1"
				fields[4] = pass // Actualizar la contraseña
				line = strings.Join(fields, ", ") + "\n"
				userFound = true
				userLine = line
				break
			} else if fields[3] == user {
				// Si el usuario ya existe y no está eliminado, no permitir su creación
				fmt.Fprintf(buffer, "Error MKUSR: El usuario '%s' ya existe.\n", user)
				userFound = true
				return
			}
			userID++
		} else if len(fields) == 3 && fields[1] == "G" {
			// Verificar si el grupo existe
			if fields[2] == grp {
				if fields[0] == "0" {
					fmt.Fprintf(buffer, "Error MKUSR: El grupo '%s' ya fue eliminado y no puede usarse.\n", grp)
					return
				}
				grupoExiste = true
			}
		}
	}

	if userFound {
		// El usuario fue restaurado
		lines = append(lines, userLine) // Añadir la línea restaurada
		// Escribir el archivo actualizado
		newData := strings.Join(lines, "\n")
		err := AppendToFileBlock(&usersInode, indexInode, newData, file, tempSuperblock, buffer)
		if err != nil {
			fmt.Fprintf(buffer, "Error MKUSR: Al escribir el archivo actualizado: %v\n", err)
			return
		}
		fmt.Fprintf(buffer, "Usuario '%s' restaurado exitosamente.\n", user)
		return
	}

	// Verificar si el grupo existe
	if !grupoExiste {
		// El grupo no existe
		fmt.Fprintf(buffer, "Error MKUSR: El grupo '%s' no existe.\n", grp)
		return
	}

	// Crear un nuevo usuario si no existe
	newUser := fmt.Sprintf("%d,U,%s,%s,%s\n", userID, grp, user, pass)

	// Añadir el nuevo usuario al archivo
	err = AppendToFileBlock(&usersInode, indexInode, newUser, file, tempSuperblock, buffer)
	if err != nil {
		fmt.Fprintf(buffer, "Error MKUSR: Al escribir nuevo usuario: %v\n", err)
		return
	}

	fmt.Fprintf(buffer, "Usuario '%s' creado exitosamente en el grupo '%s'.\n", user, grp)
}
````

Este comando permite crear un nuevo usuario en la partición actualmente montada. El usuario creado se agrega al archivo lógico `users.txt` y se asocia a un grupo ya existente. Solo puede ser ejecutado por el usuario `root`.

---

- Solo `root` puede ejecutar este comando.  
- Valida los parámetros obligatorios: `-user`, `-pass` y `-grp`.  
- Verifica que el grupo especificado exista y no haya sido eliminado.  
- Rechaza nombres de usuario duplicados.  
- El nombre del usuario, grupo y contraseña no puede superar 10 caracteres.  
- Si el usuario existía pero fue eliminado (ID = 0), se permite su restauración.  
- El archivo `users.txt` se actualiza correctamente, con el formato adecuado.

---

#### Descripción de la Implementación

##### Función principal: `Mkusr(user, pass, grp, buffer)`

1. **Validaciones iniciales**
   - Verifica que el usuario actual es `root`.
   - Valida longitud máxima (10 caracteres) en `user`, `pass` y `grp`.

2. **Localiza la partición activa**
   - Busca entre las particiones montadas aquella que coincida con `Data.GetIDPartition()`.

3. **Lectura del disco y estructuras**
   - Abre el disco.
   - Lee el `MBR` y `Superblock`.
   - Busca el inodo del archivo `/users.txt`.

4. **Lectura del contenido de `users.txt`**
   - Lee todas las líneas del archivo y analiza:
     - Si el usuario ya existe y está activo: muestra error.
     - Si el usuario ya existía pero fue eliminado (ID = 0): lo **restaura** con nueva contraseña.
     - Verifica que el grupo exista y no haya sido eliminado.

---

### Comando `RMUSR`

```` go
func Rmusr(user string, buffer *bytes.Buffer) {
	fmt.Fprint(buffer, "=============RMUSR=============\n")

	if Data.GetIDUsuario() != "root" {
		fmt.Fprintf(buffer, "Error RMUSR: Solo 'root' puede eliminar usuarios.\n")
		return
	}

	mountedPartitions := DiskManagement.GetMountedPartitions()
	var filePath string
	var partitionFound bool

	for _, Particiones := range mountedPartitions {
		for _, Particion := range Particiones {
			if Particion.ID == Data.GetIDPartition() {
				filePath = Particion.Path
				partitionFound = true
				break
			}
		}
		if partitionFound {
			break
		}
	}

	if !partitionFound {
		fmt.Fprintf(buffer, "Error RMUSR: No se encontró ninguna partición montada con el ID: %s\n", Data.GetIDPartition())
		return
	}

	file, err := Utilities.OpenFile(filePath, buffer)
	if err != nil {
		fmt.Fprintf(buffer, "Error RMUSR: No se pudo abrir disco: %v\n", err)
		return
	}
	defer file.Close()

	var TempMBR Structs.MRB
	if err := Utilities.ReadObject(file, &TempMBR, 0, buffer); err != nil {
		fmt.Fprintf(buffer, "Error RMUSR: No se pudo leer el MBR: %v\n", err)
		return
	}

	var partitionIndex int = -1
	for i, part := range TempMBR.MbrPartitions {
		if strings.Trim(string(part.ID[:]), "\x00") == Data.GetIDPartition() && part.Status[0] == '1' {
			partitionIndex = i
			break
		}
	}
	if partitionIndex == -1 {
		fmt.Fprintf(buffer, "Error RMUSR: La partición no está montada o no existe.\n")
		return
	}

	var tempSuperblock Structs.Superblock
	sbStart := int64(TempMBR.MbrPartitions[partitionIndex].Start)
	if err := Utilities.ReadObject(file, &tempSuperblock, sbStart, buffer); err != nil {
		fmt.Fprintf(buffer, "Error RMUSR: No se pudo leer Superblock: %v\n", err)
		return
	}

	indexInode := InitSearch("/users.txt", file, tempSuperblock, buffer)
	if indexInode == -1 {
		fmt.Fprintf(buffer, "Error RMUSR: No se encontró el archivo /users.txt\n")
		return
	}

	var usersInode Structs.Inode
	inodePos := int64(tempSuperblock.SB_Inode_Start + indexInode*int32(binary.Size(Structs.Inode{})))
	if err := Utilities.ReadObject(file, &usersInode, inodePos, buffer); err != nil {
		fmt.Fprintf(buffer, "Error RMUSR: No se pudo leer el Inodo: %v\n", err)
		return
	}

	data := GetInodeFileData(usersInode, file, tempSuperblock, buffer)
	lines := strings.Split(data, "\n")

	var updatedLines []string
	userFound := false

	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}

		fields := strings.Split(line, ",")
		if len(fields) == 5 && fields[1] == "U" && fields[3] == user && fields[0] != "0" {
			fields[0] = "0" // Marcar como eliminado
			userFound = true
			updatedLine := strings.Join(fields, ",")
			updatedLines = append(updatedLines, updatedLine)
		} else {
			updatedLines = append(updatedLines, line)
		}
	}

	if !userFound {
		fmt.Fprintf(buffer, "Error RMUSR: El usuario '%s' no existe o ya fue eliminado.\n", user)
		return
	}

	newData := strings.Join(updatedLines, "\n") + "\n"

	err = OverwriteFileBlock(&usersInode, indexInode, newData, file, tempSuperblock, buffer)
	if err != nil {
		fmt.Fprintf(buffer, "Error RMUSR: Al escribir el archivo actualizado: %v\n", err)
		return
	}

	fmt.Fprintf(buffer, "Usuario '%s' eliminado exitosamente.\n", user)
}
````

Este comando permite eliminar lógicamente un usuario en el sistema de archivos. El usuario será marcado como eliminado (ID = 0) en el archivo `users.txt`, siguiendo el esquema del sistema de administración de usuarios y grupos definido en el proyecto.


- Solo el usuario `root` puede ejecutar el comando.
- Elimina lógicamente al usuario colocando `0` como UID.
- Verifica que el usuario exista antes de eliminarlo.
- No permite eliminar un usuario que ya ha sido eliminado.
- Actualiza correctamente el archivo `users.txt`.

---

#### Descripción de la Implementación

##### Función principal: `Rmusr(user string, buffer *bytes.Buffer)`

1. **Verificación de permisos**
   - Se valida que el usuario actual (`Data.GetIDUsuario()`) sea `root`. Si no lo es, se lanza un error.

2. **Búsqueda de la partición activa**
   - Recorre las particiones montadas y verifica si el ID de la sesión actual (`Data.GetIDPartition()`) corresponde a alguna de ellas.

3. **Acceso a estructuras del sistema**
   - Abre el archivo de disco.
   - Lee el `MBR` y encuentra la partición activa.
   - Lee el `Superblock` de la partición.
   - Encuentra el índice del inodo de `/users.txt`.

4. **Lectura del archivo `users.txt`**
   - Carga el contenido del inodo asociado al archivo.
   - Recorre línea por línea buscando al usuario especificado.

5. **Eliminación lógica**
   - Si encuentra una línea con un usuario con `UID != 0` que coincide en nombre:
     - Cambia su UID a `0`.
     - Marca `userFound = true`.

6. **Actualización del archivo**
   - Si el usuario fue encontrado y modificado, se sobreescribe `users.txt` con las nuevas líneas.
   - En caso contrario, se lanza un error: `El usuario 'user1' no existe o ya fue eliminado`.

---

### Comando `CHGRP`

````go 
func Chgrp(user string, newGroup string, buffer *bytes.Buffer) {
	fmt.Fprint(buffer, "=============CHGRP=============\n")

	// Validar que solo 'root' puede ejecutar el comando
	if Data.GetIDUsuario() != "root" {
		fmt.Fprintf(buffer, "Error CHGRP: Solo el usuario 'root' puede cambiar grupos.\n")
		return
	}

	// Buscar la partición activa
	mountedPartitions := DiskManagement.GetMountedPartitions()
	var filePath string
	var partitionFound bool

	for _, Particiones := range mountedPartitions {
		for _, Particion := range Particiones {
			if Particion.ID == Data.GetIDPartition() {
				filePath = Particion.Path
				partitionFound = true
				break
			}
		}
		if partitionFound {
			break
		}
	}

	if !partitionFound {
		fmt.Fprintf(buffer, "Error CHGRP: No se encontró partición montada con ID: %s\n", Data.GetIDPartition())
		return
	}

	// Abrir el archivo
	file, err := Utilities.OpenFile(filePath, buffer)
	if err != nil {
		fmt.Fprintf(buffer, "Error CHGRP: No se pudo abrir el archivo: %v\n", err)
		return
	}
	defer file.Close()

	// Leer MBR y superblock
	var TempMBR Structs.MRB
	if err := Utilities.ReadObject(file, &TempMBR, 0, buffer); err != nil {
		fmt.Fprintf(buffer, "Error CHGRP: No se pudo leer el MBR: %v\n", err)
		return
	}

	var partitionIndex int = -1
	for i, part := range TempMBR.MbrPartitions {
		if strings.Trim(string(part.ID[:]), "\x00") == Data.GetIDPartition() && part.Status[0] == '1' {
			partitionIndex = i
			break
		}
	}
	if partitionIndex == -1 {
		fmt.Fprintf(buffer, "Error CHGRP: La partición no está montada o no existe.\n")
		return
	}

	var tempSuperblock Structs.Superblock
	sbStart := int64(TempMBR.MbrPartitions[partitionIndex].Start)
	if err := Utilities.ReadObject(file, &tempSuperblock, sbStart, buffer); err != nil {
		fmt.Fprintf(buffer, "Error CHGRP: No se pudo leer el Superblock: %v\n", err)
		return
	}

	// Buscar users.txt
	indexInode := InitSearch("/users.txt", file, tempSuperblock, buffer)
	if indexInode == -1 {
		fmt.Fprintf(buffer, "Error CHGRP: No se encontró el archivo /users.txt\n")
		return
	}

	var usersInode Structs.Inode
	inodePos := int64(tempSuperblock.SB_Inode_Start + indexInode*int32(binary.Size(Structs.Inode{})))
	if err := Utilities.ReadObject(file, &usersInode, inodePos, buffer); err != nil {
		fmt.Fprintf(buffer, "Error CHGRP: No se pudo leer el inodo de users.txt: %v\n", err)
		return
	}

	data := GetInodeFileData(usersInode, file, tempSuperblock, buffer)
	lines := strings.Split(data, "\n")

	var updatedLines []string
	var userFound, groupExists bool

	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}

		fields := strings.Split(line, ",")

		// Verificar si el grupo existe (ID != 0)
		if len(fields) == 3 && fields[1] == "G" && fields[2] == newGroup && fields[0] != "0" {
			groupExists = true
		}

		// Buscar al usuario y cambiar su grupo
		if len(fields) == 5 && fields[1] == "U" && fields[3] == user && fields[0] != "0" {
			userFound = true
			fields[2] = newGroup // Cambiar grupo
			line = strings.Join(fields, ",")
		}

		updatedLines = append(updatedLines, line)
	}

	if !groupExists {
		fmt.Fprintf(buffer, "Error CHGRP: El grupo '%s' no existe o está eliminado.\n", newGroup)
		return
	}

	if !userFound {
		fmt.Fprintf(buffer, "Error CHGRP: El usuario '%s' no existe o está eliminado.\n", user)
		return
	}

	newData := strings.Join(updatedLines, "\n") + "\n"
	err = OverwriteFileBlock(&usersInode, indexInode, newData, file, tempSuperblock, buffer)
	if err != nil {
		fmt.Fprintf(buffer, "Error CHGRP: No se pudo actualizar el archivo: %v\n", err)
		return
	}

	fmt.Fprintf(buffer, "Grupo del usuario '%s' actualizado exitosamente a '%s'.\n", user, newGroup)
}
````

El comando `CHGRP` permite cambiar el grupo al que pertenece un usuario ya existente en el sistema de archivos simulado. Esta operación modifica el contenido del archivo lógico `users.txt` ubicado en la raíz del sistema EXT2.

---

- Solo puede ejecutarlo el usuario `root`.
- Requiere obligatoriamente los parámetros `-user` y `-grp`.
- Verifica que el grupo exista y no esté eliminado (`ID ≠ 0`).
- Verifica que el usuario exista y no esté eliminado (`ID ≠ 0`).
- Actualiza el grupo del usuario directamente en el archivo `users.txt`.

---

#### Descripción de la Implementación

##### Función principal: `Chgrp(user string, newGroup string, buffer *bytes.Buffer)`

1. **Verificación de privilegios**  
   - Se valida que el usuario logueado sea `root`, de lo contrario se lanza un error.

2. **Búsqueda de la partición activa**  
   - Se recorre la lista de particiones montadas para obtener la ruta del archivo de disco correspondiente a la sesión activa (`Data.GetIDPartition()`).

3. **Apertura del archivo y lectura de estructuras**  
   - Se abre el archivo binario del disco.
   - Se leen el `MBR` y el `Superblock`.
   - Se busca el inodo asociado al archivo `/users.txt`.

4. **Lectura y procesamiento de `users.txt`**  
   - Se carga el contenido del archivo.
   - Se verifica si el grupo nuevo (`newGroup`) existe y no está eliminado.
   - Se busca el usuario y se actualiza su grupo si no está eliminado.

5. **Actualización del archivo**  
   - Si el grupo y el usuario son válidos, se reescribe `users.txt` con el nuevo contenido.

---

#### Utilities 
Este módulo contiene **funciones utilitarias esenciales** para la manipulación de archivos binarios utilizados en el sistema de archivos simulado. Centraliza operaciones comunes como lectura, escritura y apertura de archivos, de forma que otros módulos (como `FileSystem`, `User`, etc.) puedan reutilizarlas sin repetir código.

---

### Funcionalidades principales

#### Manejo de archivos binarios
- `CreateFile(name, buffer)`: 
  - Crea un archivo binario en la ruta especificada, incluyendo todos los directorios intermedios si no existen.
  - Solo crea el archivo si no existe previamente.

- `OpenFile(name, buffer)`:
  - Abre un archivo binario en modo lectura/escritura.
  - Devuelve un puntero al archivo o error en caso de fallo.

- `DeleteFile(name, buffer)`:
  - Elimina físicamente el archivo del disco si existe.
  - Utilizado principalmente en operaciones de limpieza o `rmdisk`.

---

#### Lectura y escritura binaria
- `WriteObject(file, data, position, buffer)`:
  - Escribe una estructura Go (como `Superblock`, `Inode`, etc.) en una posición específica del archivo binario.
  - Utiliza codificación *Little Endian*.

- `ReadObject(file, data, position, buffer)`:
  - Lee desde una posición específica del archivo y decodifica los datos binarios en la estructura Go correspondiente.
  - Es usada extensivamente para leer MBRs, superbloques, inodos, bloques, etc.

---

### `CreateFile(name string, buffer *bytes.Buffer) error`

```` go
func CreateFile(name string, buffer *bytes.Buffer) error {
	//Se asegura que el archivo existe
	dir := filepath.Dir(name)
	if err := os.MkdirAll(dir, os.ModePerm); err != nil {
		fmt.Fprintf(buffer, "Err CreateFile dir== %v.\n", err)
		return err
	}

	// Crear archivo
	if _, err := os.Stat(name); os.IsNotExist(err) {
		file, err := os.Create(name)
		if err != nil {
			fmt.Fprintf(buffer, "Err CreateFile create== %v.\n", err)
			return err
		}
		defer file.Close()
	}
	return nil
}
````

**Descripción:**  
Crea un archivo binario en la ruta especificada, asegurándose previamente de que el directorio padre exista.

**Pasos que realiza:**
1. Verifica si la carpeta del archivo existe, si no, la crea usando `os.MkdirAll`.
2. Si el archivo **no existe**, lo crea con `os.Create`.
3. Si el archivo ya existe, **no lo sobrescribe**.

**Uso típico:** Crear un nuevo disco virtual antes de aplicar comandos como `mkdisk`.

---

### `OpenFile(name string, buffer *bytes.Buffer) (*os.File, error)`

````go
func OpenFile(name string, buffer *bytes.Buffer) (*os.File, error) {
	file, err := os.OpenFile(name, os.O_RDWR, 0644)
	if err != nil {
		fmt.Fprintf(buffer, "Err OpenFile== %v.\n", err)
		return nil, err
	}
	return file, nil
}
````

**Descripción:**  
Abre un archivo en modo lectura y escritura (`O_RDWR`).

**Detalles:**
- Retorna un puntero al archivo abierto para permitir lectura y escritura con `Seek`, `Read`, `Write`.
- Si el archivo no existe, se retorna un error.

**Uso típico:** Usado por comandos como `mkfs`, `mount`, `cat`, `mkfile`, `mkgrp`, etc.

---

### `WriteObject(file *os.File, data interface{}, position int64, buffer *bytes.Buffer) error`

```` go
func WriteObject(file *os.File, data interface{}, position int64, buffer *bytes.Buffer) error {
	file.Seek(position, 0)
	err := binary.Write(file, binary.LittleEndian, data)
	if err != nil {
		fmt.Fprintf(buffer, "Err WriteObject== %v.\n", err)
		return err
	}
	return nil

}
````

**Descripción:**  
Escribe una estructura Go serializada (cualquier `struct`) en una posición específica dentro del archivo binario.

**Detalles:**
- Usa `binary.Write` con orden `LittleEndian`.
- Se posiciona en el archivo con `file.Seek(position, 0)`.

**Uso típico:** Guardar estructuras como `Superblock`, `Inode`, `FolderBlock`, etc.

---

### `ReadObject(file *os.File, data interface{}, position int64, buffer *bytes.Buffer) error`

```` go
func ReadObject(file *os.File, data interface{}, position int64, buffer *bytes.Buffer) error {
	file.Seek(position, 0)
	err := binary.Read(file, binary.LittleEndian, data)
	if err != nil {
		fmt.Fprintf(buffer, "Err ReadObject== %v.\n", err)
		return err
	}
	return nil
}
````

**Descripción:**  
Lee una estructura desde el archivo binario en una posición determinada.

**Detalles:**
- Posiciona el puntero del archivo con `Seek`.
- Utiliza `binary.Read` para mapear los bytes leídos a una estructura.

**Uso típico:** Leer estructuras del disco para procesarlas o modificarlas.

---

### `DeleteFile(name string, buffer *bytes.Buffer) error`

```` go
func DeleteFile(name string, buffer *bytes.Buffer) error {
	if _, err := os.Stat(name); os.IsNotExist(err) {
		fmt.Fprintf(buffer, "Err archive don't exist: %v.\n", err)
		return err
	}
	err := os.Remove(name)
	if err != nil {
		fmt.Fprintf(buffer, "Error al eliminar el archivo: %v.\n", err)
		return err
	}
	return nil
}
````

**Descripción:**  
Elimina un archivo binario del sistema de archivos del sistema operativo (no del FS simulado).

**Detalles:**
- Primero verifica si el archivo existe con `os.Stat`.
- Si existe, lo elimina con `os.Remove`.

**Uso típico:** Se utiliza en el comando `rmdisk` para borrar discos virtuales.

---

## Estructuras
Este archivo contiene las estructuras fundamentales utilizadas para representar las particiones, sistema de archivos EXT2, inodos, bloques y metainformación del disco.

### MBR (Master Boot Record)
Este se encuentra en el primer sector del disco y su objetivo principal es almacenar informacion sobre el sistema de archivos y las particiones que contiene el disco.

````go
type MRB struct {
	MbrSize      int32
	CreationDate [10]byte
	Signature    int32
	Fit          [1]byte
	MbrPartitions [4]Partition
}
````

Esta esctructura esta compuesta por un tamanio, una fecha de creacion, un signature, ajuste y particiones, las cuales pueden ser primarias, logicas o extendidas.

Cabe resaltar un detalle en las particiones del disco. Este arreglo es capaz de almacenar hasta 4 particiones, cada una de ellas representada por la estructura **`partition`**, la cual contiene informacion relevante sobre cada particion.

Su funcion principal es contener informacion sobre la estructura del disco y las particiones. Es el primer sector del disco, cuando un sistema operativo necesita acceder a una particion , primero consulta el MBR para obtener informacion sobre las particiones disponibles.

### Superblock
Es una de las estructuras más críticas en un sistema de archivos, ya que almacena información clave sobre el sistema de archivos en sí, como el número total de bloques, la cantidad de bloques libres, el tamaño del sistema de archivos, la ubicación de las tablas de inodos y bloques, y otros parámetros esenciales para la gestión y operación del sistema de archivos.

#### Estructura del Superblock

La estructura del **Superblock** contiene varios campos importantes que se describen a continuación:

| Nombre                     | Tipo de Dato | Descripción |
|----------------------------|--------------|-------------|
| `SB_FileSystem_Type`        | `int32`      | Tipo de sistema de archivos (por ejemplo, EXT2) |
| `SB_Inodes_Count`           | `int32`      | Número total de inodos disponibles en el sistema de archivos |
| `SB_Blocks_Count`           | `int32`      | Número total de bloques en el sistema de archivos |
| `SB_Free_Blocks_Count`      | `int32`      | Número de bloques libres disponibles |
| `SB_Free_Inodes_Count`      | `int32`      | Número de inodos libres disponibles |
| `SB_Mtime`                  | `[17]byte`   | Fecha y hora de la última vez que el sistema de archivos fue montado |
| `SB_Umtime`                 | `[17]byte`   | Fecha y hora de la última vez que el sistema de archivos fue desmontado |
| `SB_Mnt_Count`              | `int32`      | Número de veces que el sistema de archivos ha sido montado |
| `SB_Magic`                  | `int32`      | Valor que identifica al sistema de archivos (en EXT2, siempre es `0xEF53`) |
| `SB_Inode_Size`             | `int32`      | Tamaño de cada inodo en bytes |
| `SB_Block_Size`             | `int32`      | Tamaño de cada bloque en bytes |
| `SB_Fist_Ino`               | `int32`      | Primer inodo libre (índice en la tabla de inodos) |
| `SB_First_Blo`              | `int32`      | Primer bloque libre (índice en la tabla de bloques) |
| `SB_Bm_Inode_Start`         | `int32`      | Dirección de inicio del bitmap de inodos |
| `SB_Bm_Block_Start`         | `int32`      | Dirección de inicio del bitmap de bloques |
| `SB_Inode_Start`            | `int32`      | Dirección de inicio de la tabla de inodos |
| `SB_Block_Start`            | `int32`      | Dirección de inicio de la tabla de bloques |

#### Descripción de los Campos del Superblock

1. **SB_FileSystem_Type**: Este campo indica el tipo de sistema de archivos utilizado. En este caso, para el sistema de archivos EXT2, el valor será un número identificador del sistema de archivos.

2. **SB_Inodes_Count**: Este campo contiene el número total de inodos en el sistema de archivos. Los inodos son estructuras que almacenan metadatos de los archivos (propietario, permisos, ubicación de los bloques de datos, etc.). Este número es vital para la gestión de archivos.

3. **SB_Blocks_Count**: Indica el número total de bloques que componen el sistema de archivos. Los bloques son las unidades mínimas de almacenamiento donde se almacenan los datos. La cantidad de bloques define la capacidad de almacenamiento total del sistema de archivos.

4. **SB_Free_Blocks_Count**: Muestra la cantidad de bloques libres disponibles en el sistema de archivos. Este campo es importante para determinar cuánto espacio queda en el sistema de archivos para almacenar nuevos datos.

5. **SB_Free_Inodes_Count**: Similar al campo anterior, pero para los inodos. Muestra la cantidad de inodos libres disponibles para crear nuevos archivos o directorios.

6. **SB_Mtime y SB_Umtime**: Estos campos contienen las fechas y horas de la última vez que el sistema de archivos fue montado y desmontado, respectivamente. Son útiles para saber cuándo fue la última vez que se accedió al sistema de archivos.

7. **SB_Mnt_Count**: Este campo almacena el número de veces que el sistema de archivos ha sido montado. Puede ser útil para detectar posibles problemas con el sistema de archivos o simplemente para monitorear su uso.

8. **SB_Magic**: Este campo es un número mágico utilizado para identificar de manera única el sistema de archivos EXT2. Siempre tiene el valor `0xEF53` en un sistema de archivos EXT2.

9. **SB_Inode_Size**: El tamaño de cada inodo en bytes. El tamaño del inodo define cuánta información se puede almacenar en cada uno de los inodos. Esto incluye información como el tamaño del archivo, los permisos y las ubicaciones de los bloques de datos.

10. **SB_Block_Size**: El tamaño de cada bloque en bytes. Los bloques son unidades de almacenamiento para los datos, y su tamaño afecta la eficiencia y el rendimiento del sistema de archivos.

11. **SB_Fist_Ino**: El primer inodo libre en el sistema de archivos. Este campo es importante para gestionar la creación de nuevos inodos.

12. **SB_First_Blo**: Similar al campo anterior, pero para los bloques de datos. Indica el primer bloque libre disponible en el sistema de archivos.

13. **SB_Bm_Inode_Start y SB_Bm_Block_Start**: Estos campos indican las direcciones de inicio de los bitmaps de inodos y bloques. Los bitmaps son estructuras que indican qué inodos y bloques están en uso y cuáles están libres.

14. **SB_Inode_Start y SB_Block_Start**: Estos campos indican las direcciones de inicio de la tabla de inodos y la tabla de bloques en el sistema de archivos. La tabla de inodos almacena información sobre cada archivo, y la tabla de bloques contiene los bloques de datos reales.

#### Función del Superblock

El **Superblock** es utilizado por el sistema de archivos para mantener el control de la estructura general del sistema de archivos. Cada vez que se monta el sistema de archivos, se lee el superblock para obtener la información básica sobre el sistema de archivos y su estado. También se utiliza para encontrar las tablas de inodos y bloques, que son esenciales para la lectura y escritura de archivos.

Algunas de las tareas clave que realiza el **Superblock** son:

- Almacena información crucial sobre la cantidad de espacio libre y utilizado en el sistema de archivos.
- Almacena las direcciones de inicio de los bitmaps y las tablas de inodos y bloques.
- Mantiene las fechas de la última vez que se montó y desmontó el sistema de archivos.

### **EBR (Extended Boot Record)**

El **EBR** (Extended Boot Record) es una estructura utilizada para gestionar particiones lógicas dentro de una partición extendida. El EBR permite que una partición extendida se divida en múltiples particiones lógicas.

#### **Estructura de EBR**:

- **PartMount**: Indica si la partición está montada. El valor será 1 si está montada y 0 si no lo está.
- **PartFit**: Define el tipo de ajuste de la partición:
  - `'B'` para el mejor ajuste (Best),
  - `'F'` para el primer ajuste (First),
  - `'W'` para el peor ajuste (Worst).
- **PartStart**: Indica el byte donde comienza la partición en el disco.
- **PartSize**: El tamaño de la partición en bytes.
- **PartNext**: Apunta al próximo EBR. Si no hay más, se establece en `-1`.
- **PartName**: El nombre de la partición, con un máximo de 16 caracteres.

#### **Funcionamiento**:
- El EBR actúa como una lista enlazada que gestiona las particiones lógicas dentro de una partición extendida. Cada EBR contiene información sobre una partición lógica y apunta al siguiente EBR si existe.
- El primer EBR en una partición extendida marca el inicio de la primera partición lógica, y cada EBR subsiguiente apunta a la siguiente partición lógica dentro de la misma partición extendida.

---

### **Inode**

Un **Inode** es una estructura de datos que contiene metadatos sobre un archivo o directorio en un sistema de archivos. Cada archivo o carpeta tiene un Inode que almacena información como el propietario, el grupo, los permisos, y los bloques de datos asociados.

#### **Estructura de Inode**:

- **IN_Uid**: El ID de usuario (UID) del propietario del archivo o carpeta.
- **IN_Gid**: El ID de grupo (GID) del propietario del archivo o carpeta.
- **IN_Size**: El tamaño del archivo o carpeta en bytes.
- **IN_Atime**: La última vez que el archivo fue accedido.
- **IN_Ctime**: La fecha en que se creó el archivo o carpeta.
- **IN_Mtime**: La última vez que se modificó el archivo o carpeta.
- **IN_Block**: Un arreglo de 15 enteros que apuntan a bloques de datos. Los primeros 12 bloques son directos, el 13 es un bloque simple indirecto, el 14 es un bloque doble indirecto, y el 15 es un bloque triple indirecto.
- **IN_Type**: El tipo de archivo. 1 para archivo, 0 para carpeta.
- **IN_Perm**: Los permisos del archivo o carpeta, almacenados en una representación octal (UGO).

#### **Funcionamiento**:
- El Inode guarda información crucial sobre el archivo o carpeta, permitiendo que el sistema de archivos gestione las operaciones como lectura, escritura y modificación de archivos.
- Los Inodes son fundamentales para la gestión de archivos en sistemas de archivos basados en Unix, como EXT2.

---

### **FolderBlock**

Un **FolderBlock** es una estructura que almacena la información de los archivos y directorios dentro de un directorio. Cada bloque de carpeta contiene varias entradas, donde cada entrada representa un archivo o directorio dentro del directorio.

#### **Estructura de FolderBlock**:

- **B_Content**: Un arreglo de 4 entradas de tipo **Content**. Cada entrada contiene:
  - **B_Name**: El nombre del archivo o directorio (máximo 12 caracteres).
  - **B_Inode**: El índice del Inode asociado al archivo o directorio.

#### **Funcionamiento**:
- Cada bloque de carpeta (FolderBlock) almacena información sobre los archivos y directorios que contiene un directorio. Las entradas de **Content** dentro del bloque permiten a los Inodes apuntar a archivos y directorios dentro de un directorio.
- Los FolderBlocks se utilizan para navegar y gestionar el contenido dentro de un directorio.

---

### **FileBlock**

Un **FileBlock** es una estructura que almacena el contenido real de un archivo en un sistema de archivos. Cada bloque de archivo tiene un tamaño de 64 bytes y contiene los datos del archivo.

#### **Estructura de FileBlock**:

- **B_Content**: Un arreglo de 64 bytes que contiene los datos del archivo.

#### **Funcionamiento**:
- Los FileBlocks almacenan el contenido de los archivos. Cada bloque contiene una porción de los datos del archivo. Si un archivo es demasiado grande para caber en un solo bloque, el sistema de archivos usará múltiples bloques.
- Un archivo puede tener uno o más FileBlocks dependiendo de su tamaño, y estos bloques pueden ser directos o indirectos (a través de Inodes).


### Observaciones

- Todas las funciones reciben un `*bytes.Buffer` como parámetro, lo cual permite registrar mensajes de error o depuración para mostrar en el frontend o CLI.
- Este módulo **abstrae las operaciones de bajo nivel con el disco virtual**, garantizando una interfaz segura y reutilizable.
- Es clave para mantener la integridad del sistema, ya que asegura que las operaciones de lectura/escritura binaria se hagan de manera coherente y sin corrupción de datos.

---