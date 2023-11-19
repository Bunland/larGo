package main

/*
#cgo CFLAGS: -I/usr/include/webkitgtk-4.0
#cgo LDFLAGS: -ljavascriptcoregtk-4.0
#include <JavaScriptCore/JavaScript.h>
#include <stdlib.h>
*/
import "C"

import (
	"fmt"
)

func main() {
	// Crear un contexto JavaScript
	context := C.JSGlobalContextCreate(nil)

	// Crear un objeto JavaScript
	globalObject := C.JSContextGetGlobalObject(context)

	// Crear un script de JavaScript
	script := "var mensaje = '¡Hello world from JavaScript!'; mensaje;"
	scriptJS := C.JSStringCreateWithUTF8CString(C.CString(script))
	defer C.JSStringRelease(scriptJS)

	// Evaluar el script
	result := C.JSEvaluateScript(context, scriptJS, globalObject, nil, 1, nil)

	// Convertir el resultado a una cadena de Go
	resultStringJS := C.JSValueToStringCopy(context, result, nil)
	defer C.JSStringRelease(resultStringJS)

	bufferSize := C.JSStringGetMaximumUTF8CStringSize(resultStringJS)
	resultCString := make([]C.char, bufferSize)
	C.JSStringGetUTF8CString(resultStringJS, &resultCString[0], bufferSize)

	// Imprimir el resultado
	fmt.Printf("Resultado: %s\n", C.GoString(&resultCString[0]))

	// Liberar recursos
	C.JSGlobalContextRelease(context)
}