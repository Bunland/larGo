package fs

/*
#cgo CFLAGS: -I/usr/include/webkitgtk-4.0
#cgo LDFLAGS: -ljavascriptcoregtk-4.0
#include <stdlib.h>
#include <JavaScriptCore/JavaScript.h>

extern JSValueRef ReadFileSyncF(JSContextRef context, JSObjectRef function, JSObjectRef thisObject, size_t argumentCount, JSValueRef arguments[], JSValueRef* exception);
*/
import "C"
import (
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"os"
	"unsafe"
)

// ReadFileSyncF hace la función fs.readFileSync() de JavaScript
//
//export ReadFileSyncF
func ReadFileSyncF(context C.JSContextRef, function C.JSObjectRef, thisObject C.JSObjectRef, argumentCount C.size_t, arguments *C.JSValueRef, exception *C.JSValueRef) (finalValue C.JSValueRef) {
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
	if C.JSValueIsUndefined(context, argumentSlice[1]) {
		fmt.Println("Funciona el condicional")
		return C.JSValueMakeUndefined(context)
	}
	obj := C.JSValueToObject(context, argumentSlice[1], exception)
	if C.JSValueIsUndefined(context, obj) {
		fmt.Println("Funciona el condicional")
		return C.JSValueMakeUndefined(context)
	}
	propertyObjectC := C.CString("encoding")
	propertyObjectJS := C.JSStringCreateWithUTF8CString(propertyObjectC)
	value := C.JSObjectGetProperty(context, obj, propertyObjectJS, exception)
	C.free(unsafe.Pointer(propertyObjectC))
	C.JSStringRelease(propertyObjectJS)
	if C.JSValueIsUndefined(context, value) {
		fmt.Println("Funciona el condicional")
		return C.JSValueMakeUndefined(context)
	}
	encodingPropertyString := C.JSValueToStringCopy(context, value, exception)
	bufferSizeString := C.JSStringGetMaximumUTF8CStringSize(encodingPropertyString)
	bufferString := C.malloc(bufferSizeString)
	C.JSStringGetUTF8CString(encodingPropertyString, (*C.char)(bufferString), bufferSizeString)
	encoding := C.GoString((*C.char)(bufferString))
	C.free(unsafe.Pointer(bufferString))
	finalValue = C.JSValueMakeUndefined(context)
	switch encoding {
	case "base64":
		cString := C.CString(base64.StdEncoding.EncodeToString(file))
		finalValue = C.JSValueMakeString(context, C.JSStringCreateWithUTF8CString(cString))
		C.free(unsafe.Pointer(cString))
	case "utf8", "utf-8":
		cString := C.CString(string(file))
		fileCString := C.JSStringCreateWithUTF8CString(cString)
		C.free(unsafe.Pointer(cString))
		finalValue = C.JSValueMakeString(context, fileCString)
	case "hex":
		cString := C.CString(hex.EncodeToString(file))
		finalValue = C.JSValueMakeString(context, C.JSStringCreateWithUTF8CString(cString))
		C.free(unsafe.Pointer(cString))
	}
	return
}

// ReadFileSync devuelve la función callback de JavaScript en C para la función ReadFileSync en JavaScript
func ReadFileSync() C.JSObjectCallAsFunctionCallback {
	return C.JSObjectCallAsFunctionCallback(C.ReadFileSyncF)
}
