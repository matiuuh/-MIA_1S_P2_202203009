package Utilities

import (
	"encoding/binary"
	//"proyecto1/Structs"
	"fmt"
	"os"
	//"strings"
	//"time"
	//"math/rand"
	"bytes"
	"path/filepath"
)

// Funcion para crear un archivo binario
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

// Funcion para abrir un archivo binario ead/write mode
func OpenFile(name string, buffer *bytes.Buffer) (*os.File, error) {
	file, err := os.OpenFile(name, os.O_RDWR, 0644)
	if err != nil {
		fmt.Fprintf(buffer, "Err OpenFile== %v.\n", err)
		return nil, err
	}
	return file, nil
}

// Funcion para escribir un objecto en un archivo binario
func WriteObject(file *os.File, data interface{}, position int64, buffer *bytes.Buffer) error {
	file.Seek(position, 0)
	err := binary.Write(file, binary.LittleEndian, data)
	if err != nil {
		fmt.Fprintf(buffer, "Err WriteObject== %v.\n", err)
		return err
	}
	return nil

}

// Funcion para leer un objeto de un archivo binario
func ReadObject(file *os.File, data interface{}, position int64, buffer *bytes.Buffer) error {
	file.Seek(position, 0)
	err := binary.Read(file, binary.LittleEndian, data)
	if err != nil {
		fmt.Fprintf(buffer, "Err ReadObject== %v.\n", err)
		return err
	}
	return nil
}

// Funcion para eliminar un archivo
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