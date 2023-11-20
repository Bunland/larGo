package utils

import (
	"fmt"
	"io/ioutil"
)

// ReadFile lee el contenido de un archivo y devuelve su contenido como una cadena.
func ReadFile(filename string) string {
	// Intentar leer el contenido del archivo.
	content, err := ioutil.ReadFile(filename)

	// Verificar si hay errores al abrir el archivo.
	if err != nil {
		fmt.Println("Failed to open file")
		return ""
	}

	// Devolver el contenido del archivo como una cadena.
	return string(content)
}
