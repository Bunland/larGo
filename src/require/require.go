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
	"largo/src/utils"
	"unsafe"
)

// Require es la implementación de require de JavaScript
//
//export RequireF
func RequireF(context C.JSContextRef, function C.JSObjectRef, thisObject C.JSObjectRef, argumentCount C.size_t, arguments *C.JSValueRef, exception *C.JSValueRef) C.JSValueRef {
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

	exportValue := C.JSObjectGetProperty(context, C.JSContextGetGlobalObject(context), exportName, nil)
	C.JSStringRelease(exportName)

	return exportValue
}

// Require devuelve la función de callback de C para la función require en JavaScript.
func Require() C.JSObjectCallAsFunctionCallback {
	return C.JSObjectCallAsFunctionCallback(C.RequireF)
}
