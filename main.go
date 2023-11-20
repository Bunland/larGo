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
	"largo/src/math"
	"largo/src/utils"
	"os"
)

// createCustomFunction crea una función JavaScript personalizada y la establece como propiedad del objeto global.
func createCustomFunction(context C.JSGlobalContextRef, globalObject C.JSObjectRef, functionName string, functionCallback C.JSObjectCallAsFunctionCallback) {
	// Crear una cadena JavaScript a partir del nombre de la función en formato UTF-8.
	functionString := C.JSStringCreateWithUTF8CString(C.CString(functionName))

	// Crear un objeto de función JavaScript usando la cadena y la devolución de llamada de la función.
	functionObject := C.JSObjectMakeFunctionWithCallback(context, functionString, functionCallback)

	// Establecer la función recién creada como propiedad del objeto global.
	C.JSObjectSetProperty(context, globalObject, functionString, functionObject, C.kJSPropertyAttributeNone, nil)

	// Liberar la cadena de función creada con JSStringCreateWithUTF8CString para evitar fugas de memoria.
	C.JSStringRelease(functionString)
}

// Apis define las API disponibles en JavaScript.
func Apis(context C.JSGlobalContextRef, globalObject C.JSObjectRef) {
	createCustomFunction(context, globalObject, "Add", C.JSObjectCallAsFunctionCallback(math.Add()))
}

func main() {
	// Crear un contexto JavaScript global.
	context := C.JSGlobalContextCreate(nil)
	globalObject := C.JSContextGetGlobalObject(context)

	// Configurar las API en el objeto global.
	Apis(context, globalObject)

	// Verificar si hay argumentos de línea de comandos y si se proporciona el comando "run".
	if len(os.Args) > 2 && os.Args[1] == "run" {
		jsFileName := os.Args[2]

		// Leer el contenido del archivo JavaScript.
		fileContent := utils.ReadFile(jsFileName)

		// Crear una cadena JavaScript a partir del contenido del archivo.
		jsCode := C.JSStringCreateWithUTF8CString(C.CString(fileContent))
		defer C.JSStringRelease(jsCode)

		// Evaluar el script JavaScript.
		result := C.JSEvaluateScript(context, jsCode, globalObject, nil, 1, nil)

		// Convertir el resultado a una cadena de Go.
		resultStringJS := C.JSValueToStringCopy(context, result, nil)
		defer C.JSStringRelease(resultStringJS)

		// Obtener el tamaño máximo necesario para la cadena UTF-8.
		bufferSize := C.JSStringGetMaximumUTF8CStringSize(resultStringJS)
		resultCString := make([]C.char, bufferSize)
		C.JSStringGetUTF8CString(resultStringJS, &resultCString[0], bufferSize)

		// Imprimir el resultado.
		fmt.Printf("%s\n", C.GoString(&resultCString[0]))
	}

	// Liberar el contexto JavaScript global.
	C.JSGlobalContextRelease(context)
}
