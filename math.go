package main

/*
#cgo CFLAGS: -I/usr/include/webkitgtk-4.0
#cgo LDFLAGS: -ljavascriptcoregtk-4.0
#include <JavaScriptCore/JavaScript.h>

// Declarar la función Add para que C la reconozca
extern JSValueRef Add(JSContextRef context, JSObjectRef function, JSObjectRef thisObject, size_t argumentCount, JSValueRef arguments[], JSValueRef* exception);

*/
import "C"
import (
	"fmt"
	"unsafe"
)

//export Add
func Add(context C.JSContextRef, function C.JSObjectRef, thisObject C.JSObjectRef, argumentCount C.size_t, arguments *C.JSValueRef, exception *C.JSValueRef) C.JSValueRef {
	if argumentCount < 2 || arguments == nil {
		fmt.Println("La función requiere 2 argumentos.")
		return C.JSValueMakeUndefined(context)
	}

	// Convertir la rebanada de argumentos a una rebanada/slice de Go
	argumentSlice := (*[1 << 30]C.JSValueRef)(unsafe.Pointer(arguments))[:argumentCount:argumentCount]

	numa := int(C.JSValueToNumber(context, argumentSlice[0], nil))
	numb := int(C.JSValueToNumber(context, argumentSlice[1], nil))

	sum := numa + numb

	return C.JSValueMakeNumber(context, C.double(sum))
}

func main() {
	// Crear un contexto JavaScript
	context := C.JSGlobalContextCreate(nil)

	// Crear un objeto JavaScript
	globalObject := C.JSContextGetGlobalObject(context)

	// Definir el nombre de la función en Go y convertirlo a una cadena de JavaScript
	functionName := "add"
	functionString := C.JSStringCreateWithUTF8CString(C.CString(functionName))
	defer C.JSStringRelease(functionString)

	// Crear la función en JavaScript utilizando el nombre y el callback de la función en C
	functionObject := C.JSObjectMakeFunctionWithCallback(context, functionString, (*[0]byte)(unsafe.Pointer(C.Add)))

	// Establecer la función recién creada como una propiedad del objeto global
	C.JSObjectSetProperty(context, globalObject, functionString, functionObject, C.kJSPropertyAttributeNone, nil)

	// Llamar a la función desde JavaScript
	scriptSuma := "add(5, 7);"
	scriptSumaJS := C.JSStringCreateWithUTF8CString(C.CString(scriptSuma))
	defer C.JSStringRelease(scriptSumaJS)

	// Evaluar el script de la suma
	result := C.JSEvaluateScript(context, scriptSumaJS, globalObject, nil, 1, nil)

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
