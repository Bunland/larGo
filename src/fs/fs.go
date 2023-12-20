package fs

/*
#cgo CFLAGS: -I/usr/include/webkitgtk-4.0
#cgo LDFLAGS: -ljavascriptcoregtk-4.0
#include <stdlib.h>
#include <JavaScriptCore/JavaScript.h>

extern JSValueRef ReadFileF(JSContextRef context, JSObjectRef function, JSObjectRef thisObject, size_t argumentCount, JSValueRef arguments[], JSValueRef* exception);
*/
import "C"
import (
	"os"
	"unsafe"
)

// TODO: Hacer más métodos y el objeto de fs.readFile()

// ReadFileF hace la función fs.readFile() de JavaScript
//
//export ReadFileF
func ReadFileF(context C.JSContextRef, function C.JSObjectRef, thisObject C.JSObjectRef, argumentCount C.size_t, arguments *C.JSValueRef, exception *C.JSValueRef) C.JSValueRef {
	if int(argumentCount) < 1 {
		return C.JSValueMakeUndefined(context)
	}
	argumentSlice := (*[1 << 30]C.JSValueRef)(unsafe.Pointer(arguments))[:argumentCount:argumentCount]
	str := C.JSValueToStringCopy(context, argumentSlice[0], nil)
	bufferSize := C.JSStringGetMaximumUTF8CStringSize(str)

	buffer := C.malloc(bufferSize)
	C.JSStringGetUTF8CString(str, (*C.char)(buffer), bufferSize)
	file, err := os.ReadFile(C.GoString((*C.char)(buffer)))
	C.free(unsafe.Pointer(buffer))
	if err != nil {
		return C.JSValueMakeUndefined(context)
	}
	c_string := C.CString(string(file))
	file_c_string := C.JSStringCreateWithUTF8CString(c_string)
	C.free(unsafe.Pointer(c_string))
	return C.JSValueMakeString(context, file_c_string)
}

// ReadFile devuelve la función callback de JavaScript en C para la función ReadFile en JavaScript
func ReadFile() C.JSObjectCallAsFunctionCallback {
	return C.JSObjectCallAsFunctionCallback(C.ReadFileF)
}
