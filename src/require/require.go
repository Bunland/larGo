package require

/*
#cgo CFLAGS: -I/usr/include/webkitgtk-4.0
#cgo LDFLAGS: -ljavascriptcoregtk-4.0
#include <stdlib.h>
#include <JavaScriptCore/JavaScript.h>

extern JSValueRef RequireF(JSContextRef context, JSObjectRef function, JSObjectRef thisObject, size_t argumentCount, JSValueRef arguments[], JSValueRef* exception);
*/
import "C"
import (
	"largo/src/fs"
	"largo/src/utils"
	"unsafe"
)

// createCustomFunction crea una función JavaScript personalizada y la establece como propiedad del objeto global.
func createCustomFunction(context C.JSContextRef, globalObject C.JSObjectRef, functionName string, functionCallback C.JSObjectCallAsFunctionCallback) {
	// Crear el string de C
	functionNameC := C.CString(functionName)

	// Crear una cadena JavaScript a partir del nombre de la función en formato UTF-8.
	functionString := C.JSStringCreateWithUTF8CString(functionNameC)

	// Liberar el string de C
	C.free(unsafe.Pointer(functionNameC))

	// Crear un objeto de función JavaScript usando la cadena y la devolución de llamada de la función.
	functionObject := C.JSObjectMakeFunctionWithCallback(context, functionString, functionCallback)

	// Establecer la función recién creada como propiedad del objeto global.
	C.JSObjectSetProperty(context, globalObject, functionString, functionObject, C.kJSPropertyAttributeNone, nil)

	// Liberar la cadena de función creada con JSStringCreateWithUTF8CString para evitar fugas de memoria.
	C.JSStringRelease(functionString)
}

// Require es la implementación de require de JavaScript
//
//export RequireF
func RequireF(context C.JSContextRef, function C.JSObjectRef, thisObject C.JSObjectRef, argumentCount C.size_t, arguments *C.JSValueRef, exception *C.JSValueRef) (finalValue C.JSValueRef) {
	argumentSlice := (*[1 << 30]C.JSValueRef)(unsafe.Pointer(arguments))[:argumentCount:argumentCount]
	if len(argumentSlice) <= 0 {
		return C.JSValueMakeUndefined(context)
	}
	moduleName := C.JSValueToStringCopy(context, argumentSlice[0], nil)
	bufferSize := C.JSStringGetMaximumUTF8CStringSize(moduleName)
	buffer := C.malloc(bufferSize)
	C.JSStringGetUTF8CString(moduleName, (*C.char)(buffer), bufferSize)
	module := C.GoString((*C.char)(buffer))
	C.free(unsafe.Pointer(buffer))

	switch module {
	case "fs", "node:fs":
		fsObject := C.JSObjectMake(context, nil, nil)
		createCustomFunction(context, fsObject, "readFileSync", C.JSObjectCallAsFunctionCallback(fs.ReadFileSync()))
		createCustomFunction(context, fsObject, "writeFileSync", C.JSObjectCallAsFunctionCallback(fs.WriteFileSync()))

		finalValue = (C.JSValueRef)(fsObject)
	default:
		content := utils.ReadFile(module)

		if content == "" {
			return C.JSValueMakeUndefined(context)
		}

		exportC := C.CString("exports")
		exportName := C.JSStringCreateWithUTF8CString(exportC)
		C.free(unsafe.Pointer(exportC))
		exportObject := C.JSObjectMake(context, nil, nil)
		C.JSObjectSetProperty(context, thisObject, exportName, exportObject, C.kJSPropertyAttributeNone, nil)
		contentC := C.CString(content)
		script := C.JSStringCreateWithUTF8CString(contentC)
		C.JSEvaluateScript(context, script, nil, nil, 0, exception)

		C.JSStringRelease(script)
		C.free(unsafe.Pointer(contentC))

		finalValue = C.JSObjectGetProperty(context, C.JSContextGetGlobalObject(context), exportName, nil)
		C.JSStringRelease(exportName)
	}

	return
}

// Require devuelve la función de callback de C para la función require en JavaScript.
func Require() C.JSObjectCallAsFunctionCallback {
	return C.JSObjectCallAsFunctionCallback(C.RequireF)
}
