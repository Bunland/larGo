package console

/*
#cgo CFLAGS: -I/usr/include/webkitgtk-4.0
#cgo LDFLAGS: -ljavascriptcoregtk-4.0
#include <stdlib.h>
#include <JavaScriptCore/JavaScript.h>

// Declarar las funciones LogF, TimeF, TimeEndF para que C reconozca su existencia.
extern JSValueRef LogF(JSContextRef context, JSObjectRef function, JSObjectRef thisObject, size_t argumentCount, JSValueRef arguments[], JSValueRef* exception);
extern JSValueRef TimeF(JSContextRef context, JSObjectRef function, JSObjectRef thisObject, size_t argumentCount, JSValueRef arguments[], JSValueRef* exception);
extern JSValueRef TimeEndF(JSContextRef context, JSObjectRef function, JSObjectRef thisObject, size_t argumentCount, JSValueRef arguments[], JSValueRef* exception);
*/
import "C"
import (
	"fmt"
	"time"
	"unsafe"
)

type WatcherStruct struct {
	start time.Time
	label string
}

var watcher = make(map[string]*WatcherStruct)

// Log es la implementación de console.log de JavaScript
//
//export LogF
func LogF(context C.JSContextRef, function C.JSObjectRef, thisObject C.JSObjectRef, argumentCount C.size_t, arguments *C.JSValueRef, exception *C.JSValueRef) C.JSValueRef {
	argumentSlice := (*[1 << 30]C.JSValueRef)(unsafe.Pointer(arguments))[:argumentCount:argumentCount]
	if len(argumentSlice) <= 0 {
		return C.JSValueMakeUndefined(context)
	}
	for i := 0; i < int(argumentCount); i += 1 {
		str := C.JSValueToStringCopy(context, argumentSlice[i], nil)

		bufferSize := C.JSStringGetMaximumUTF8CStringSize(str)

		buffer := C.malloc(bufferSize)
		C.JSStringGetUTF8CString(str, (*C.char)(buffer), bufferSize)

		fmt.Printf("%s ", C.GoString((*C.char)(buffer)))

		C.free(unsafe.Pointer(buffer))

		C.JSStringRelease(str)
	}
	fmt.Print("\n")

	return C.JSValueMakeUndefined(context)
}

// Time es la implementación de console.time de JavaScript
//
//export TimeF
func TimeF(context C.JSContextRef, function C.JSObjectRef, thisObject C.JSObjectRef, argumentCount C.size_t, arguments *C.JSValueRef, exception *C.JSValueRef) C.JSValueRef {
	argumentSlice := (*[1 << 30]C.JSValueRef)(unsafe.Pointer(arguments))[:argumentCount:argumentCount]
	var label string
	if len(argumentSlice) <= 0 {
		label = "default"
	} else {
		str := C.JSValueToStringCopy(context, argumentSlice[0], nil)
		bufferSize := C.JSStringGetMaximumUTF8CStringSize(str)
		buffer := C.malloc(bufferSize)
		C.JSStringGetUTF8CString(str, (*C.char)(buffer), bufferSize)
		label = C.GoString((*C.char)(buffer))
		C.free(unsafe.Pointer(buffer))
	}
	if _, ok := watcher[label]; ok {
		return C.JSValueMakeUndefined(context)
	}
	watcher[label] = &WatcherStruct{
		start: time.Now(),
		label: label,
	}
	return C.JSValueMakeUndefined(context)
}

// TimeEnd es la implementación de console.timeEnd de JavaScript
//
//export TimeEndF
func TimeEndF(context C.JSContextRef, function C.JSObjectRef, thisObject C.JSObjectRef, argumentCount C.size_t, arguments *C.JSValueRef, exception *C.JSValueRef) C.JSValueRef {
	argumentSlice := (*[1 << 30]C.JSValueRef)(unsafe.Pointer(arguments))[:argumentCount:argumentCount]
	var label string
	if len(argumentSlice) <= 0 {
		label = "default"
	} else {
		str := C.JSValueToStringCopy(context, argumentSlice[0], nil)
		bufferSize := C.JSStringGetMaximumUTF8CStringSize(str)
		buffer := C.malloc(bufferSize)
		C.JSStringGetUTF8CString(str, (*C.char)(buffer), bufferSize)
		label = C.GoString((*C.char)(buffer))
		C.free(unsafe.Pointer(buffer))
	}
	if _, ok := watcher[label]; !ok {
		return C.JSValueMakeUndefined(context)
	}
	fmt.Printf("%s: %v\n", label, time.Since(watcher[label].start))
	delete(watcher, label)
	return C.JSValueMakeUndefined(context)
}

// Log devuelve la función de callback de C para la función Log en JavaScript.
func Log() C.JSObjectCallAsFunctionCallback {
	return C.JSObjectCallAsFunctionCallback(C.LogF)
}

// Time devuelve la función de callback de C para la función Time en JavaScript.
func Time() C.JSObjectCallAsFunctionCallback {
	return C.JSObjectCallAsFunctionCallback(C.TimeF)
}

// TimeEnd devuelve la función de callback de C para la función TimeEnd en JavaScript.
func TimeEnd() C.JSObjectCallAsFunctionCallback {
	return C.JSObjectCallAsFunctionCallback(C.TimeEndF)
}
