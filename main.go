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

	// Crear un contexto de JavaScript
	context := C.JSGlobalContextCreate(nil)
	defer C.JSGlobalContextRelease(context)

	// Crear una cadena JavaScript
	jsCode := C.JSStringCreateWithUTF8CString(C.CString("const some = 'hello sdfsdfsdf'; some"))
	defer C.JSStringRelease(jsCode)

	// Evaluar el código JavaScript en el contexto
	result := C.JSEvaluateScript(context, jsCode, nil, nil, 0, nil)
	defer C.JSValueUnprotect(context, result)

	// Verificar si ocurrió un error durante la evaluación
	if bool(C.JSValueIsUndefined(context, result)) {
		fmt.Println("Error evaluating JavaScript")
		return
	}
	jsString := C.JSValueToStringCopy(context, result, nil)
	defer C.JSStringRelease(jsString)

	// Obtener la longitud de la cadena
	length := C.JSStringGetLength(jsString)

	// Obtener la cadena en formato UTF-8
	var buffer [512]C.char
	C.JSStringGetUTF8CString(jsString, &buffer[0], length+1)

	// Convertir el puntero C a una cadena de Go
	goString := C.GoString(&buffer[0])

	// Imprimir la cadena
	fmt.Println(goString)
}
