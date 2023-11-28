package math

/*
#cgo CFLAGS: -I/usr/include/webkitgtk-4.0
#cgo LDFLAGS: -ljavascriptcoregtk-4.0
#include <JavaScriptCore/JavaScript.h>

// Declarar la función AddFoo para que C la reconozca
extern JSValueRef AddFoo(JSContextRef context, JSObjectRef function, JSObjectRef thisObject, size_t argumentCount, JSValueRef arguments[], JSValueRef* exception);
extern JSValueRef MultF(JSContextRef context, JSObjectRef function, JSObjectRef thisObject, size_t argumentCount, JSValueRef arguments[], JSValueRef* exception);
*/
import "C"
import (
	"unsafe"
)

// AddFoo es la implementación de la función Add en C que será llamada desde JavaScript.
//
//export AddFoo
func AddFoo(context C.JSContextRef, function C.JSObjectRef, thisObject C.JSObjectRef, argumentCount C.size_t, arguments *C.JSValueRef, exception *C.JSValueRef) C.JSValueRef {
	// Verificar si se proporcionan al menos 2 argumentos.
	if int(argumentCount) < 2 || arguments == nil {
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

// MultF es la implementación de la función Mult en C que será llamada desde JavaScript.
//
//export MultF
func MultF(context C.JSContextRef, function C.JSObjectRef, thisObject C.JSObjectRef, argumentCount C.size_t, arguments *C.JSValueRef, exception *C.JSValueRef) C.JSValueRef {
	if int(argumentCount) < 2 || arguments == nil {
		return C.JSValueMakeUndefined(context)
	}

	argumentSlice := (*[1 << 30]C.JSValueRef)(unsafe.Pointer(arguments))[:argumentCount:argumentCount]

	numa := int(C.JSValueToNumber(context, argumentSlice[0], nil))
	numb := int(C.JSValueToNumber(context, argumentSlice[1], nil))

	return C.JSValueMakeNumber(context, C.double(numa*numb))
}

// Add devuelve la función de callback de C para la función Add en JavaScript.
func Add() C.JSObjectCallAsFunctionCallback {
	return C.JSObjectCallAsFunctionCallback(C.AddFoo)
}

// Mult devuelve la función de callback de C para la función Mult en JavaScript.
func Mult() C.JSObjectCallAsFunctionCallback {
	return C.JSObjectCallAsFunctionCallback(C.MultF)
}
