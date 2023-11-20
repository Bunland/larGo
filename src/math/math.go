package math

/*
#cgo CFLAGS: -I/usr/include/webkitgtk-4.0
#cgo LDFLAGS: -ljavascriptcoregtk-4.0
#include <JavaScriptCore/JavaScript.h>

// Declarar la función AddFoo para que C la reconozca
extern JSValueRef AddFoo(JSContextRef context, JSObjectRef function, JSObjectRef thisObject, size_t argumentCount, JSValueRef arguments[], JSValueRef* exception);
*/
import "C"
import (
	"fmt"
	"unsafe"
)

// AddFoo es la implementación de la función Add en C que será llamada desde JavaScript.
//export AddFoo
func AddFoo(context C.JSContextRef, function C.JSObjectRef, thisObject C.JSObjectRef, argumentCount C.size_t, arguments *C.JSValueRef, exception *C.JSValueRef) C.JSValueRef {
	// Verificar si se proporcionan al menos 2 argumentos.
	if argumentCount < 2 || arguments == nil {
		fmt.Println("La función requiere 2 argumentos.")
		return C.JSValueMakeUndefined(context)
	}

	// Convertir la rebanada de argumentos a una rebanada/slice de Go.
	argumentSlice := (*[1 << 30]C.JSValueRef)(unsafe.Pointer(arguments))[:argumentCount:argumentCount]

	// Convertir los valores JavaScript a números enteros de Go.
	numa := int(C.JSValueToNumber(context, argumentSlice[0], nil))
	numb := int(C.JSValueToNumber(context, argumentSlice[1], nil))

	// Calcular la suma.
	sum := numa + numb

	// Devolver el resultado como un valor JavaScript.
	return C.JSValueMakeNumber(context, C.double(sum))
}

// Add devuelve la función de callback de C para la función Add en JavaScript.
func Add() C.JSObjectCallAsFunctionCallback {
	return C.JSObjectCallAsFunctionCallback(C.AddFoo)
}
