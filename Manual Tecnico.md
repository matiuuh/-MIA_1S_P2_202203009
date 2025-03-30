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
sefwefwefwe

#### DiskManagement
Dicho apartado le permite al administrador modificar los datos de cualquier usuario del sistema.

#### FileSystem
Dicho apartado le permite al administrador eliminar los datos de cualquier usuario del sistema.

#### Report
Dicha función se encarga de buscar a un usuario mediante su correo electrónico y lo muestra en la respectiva tabla.

#### Structs
En ese mismo apartado se encuentra un combobox y un botón que al seleccionar el tipo de recorrido y pulsar el botón en la tabla se mostrarán los usuarios ordenados en el respectivo recorrido seleccionado.

#### User
Dicho módulo fue programado para el uso exclusivo de los usuarios de la red social. A dicho módulo se puede acceder por medio de las credenciales creadas en el registro de usuario. Al momento de ingresar 

#### Utilities 
Método encargado de buscar a un usuario por medio de su correo, dicho método obtiene el texto de la interfaz gráfica para poder manejarlo y buscar a dicho usuario en el árbol avl.

```cpp
//**************************BUSCAR****************************
//--------------BUSCAR A UN USUARIO POR SU CORREO---------------
void InterfazPrincipal::buscarUsuarioCorreo() {
    // Obtener el correo ingresado por el usuario
    QString correoBuscar = ui->txt_correoBuscar->toPlainText();

    // Convertir el QString a std::string para usarlo con AVLUsuarios
    std::string correoBuscarStd = correoBuscar.toStdString();

    // Obtener una instancia del AVLUsuarios (suponiendo que es singleton o global)
    AVLUsuarios& avlUsuarios = AVLUsuarios::getInstance();

    // Buscar el usuario por correo
    Usuario* usuarioEncontrado = avlUsuarios.buscar(correoBuscarStd);

    // Si se encuentra el usuario, mostrar sus datos en los campos correspondientes
    if (usuarioEncontrado) {
        ui->txt_nombreEncontrado->setText(QString::fromStdString(usuarioEncontrado->getNombre()));
        ui->txt_apellidoEncontrado->setText(QString::fromStdString(usuarioEncontrado->getApellidos()));
        ui->txt_correoEncontrado->setText(QString::fromStdString(usuarioEncontrado->getCorreo()));
        ui->txt_fechaEncontrada->setText(QString::fromStdString(usuarioEncontrado->getFecha()));

        // Bloquear los campos para que no sean editables
        ui->txt_nombreEncontrado->setReadOnly(true);
        ui->txt_apellidoEncontrado->setReadOnly(true);
        ui->txt_correoEncontrado->setReadOnly(true);
        ui->txt_fechaEncontrada->setReadOnly(true);

    } else {
        // Si no se encuentra el usuario, limpiar los campos y mostrar un mensaje
        ui->txt_nombreEncontrado->clear();
        ui->txt_apellidoEncontrado->clear();
        ui->txt_correoEncontrado->clear();
        ui->txt_fechaEncontrada->clear();

        // Mostrar mensaje al usuario
        QMessageBox::warning(this, "Usuario no encontrado", "El correo ingresado no corresponde a ningún usuario.");
    }
}

```


## Consideraciones Finales

Este proyecto está diseñado para ser modular, escalable y adaptable a futuras expansiones. El uso de estructuras de datos eficientes y la gestión cuidadosa de la memoria mediante punteros inteligentes aseguran la estabilidad y el rendimiento del sistema. Cualquier modificación o expansión debe seguir los principios establecidos para mantener la consistencia y eficiencia.