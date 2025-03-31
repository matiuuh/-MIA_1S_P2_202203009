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

### Observaciones

- Todas las funciones reciben un `*bytes.Buffer` como parámetro, lo cual permite registrar mensajes de error o depuración para mostrar en el frontend o CLI.
- Este módulo **abstrae las operaciones de bajo nivel con el disco virtual**, garantizando una interfaz segura y reutilizable.
- Es clave para mantener la integridad del sistema, ya que asegura que las operaciones de lectura/escritura binaria se hagan de manera coherente y sin corrupción de datos.

---