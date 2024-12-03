package fs

/*
#cgo CFLAGS: -I/usr/include/webkitgtk-4.0
#cgo LDFLAGS: -ljavascriptcoregtk-4.0
#include <JavaScriptCore/JavaScript.h>
#include <stdlib.h>

extern JSValueRef ReadFileF(JSContextRef context, JSObjectRef function, JSObjectRef thisObject, size_t argumentCount, JSValueRef arguments[], JSValueRef* exception);
extern JSValueRef ReadFileSyncF(JSContextRef context, JSObjectRef function, JSObjectRef thisObject, size_t argumentCount, JSValueRef arguments[], JSValueRef* exception);
extern JSValueRef WriteFileSyncF(JSContextRef context, JSObjectRef function, JSObjectRef thisObject, size_t argumentCount, JSValueRef arguments[], JSValueRef* exception);
extern JSValueRef ReadDirSyncF(JSContextRef context, JSObjectRef function, JSObjectRef thisObject, size_t argumentCount, JSValueRef arguments[], JSValueRef* exception);
extern JSValueRef MkDirSyncF(JSContextRef context, JSObjectRef function, JSObjectRef thisObject, size_t argumentCount, JSValueRef arguments[], JSValueRef* exception);
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
		jsString := C.JSStringCreateWithUTF8CString(cString)
		finalValue = C.JSValueMakeString(context, jsString)
		C.JSStringRelease(jsString)
		C.free(unsafe.Pointer(cString))
	case "utf8", "utf-8":
		cString := C.CString(string(file))
		fileCString := C.JSStringCreateWithUTF8CString(cString)
		C.free(unsafe.Pointer(cString))
		finalValue = C.JSValueMakeString(context, fileCString)
		C.JSStringRelease(fileCString)
	case "hex":
		cString := C.CString(hex.EncodeToString(file))
		jsString := C.JSStringCreateWithUTF8CString(cString)
		finalValue = C.JSValueMakeString(context, jsString)
		C.free(unsafe.Pointer(cString))
		C.JSStringRelease(jsString)
	}
	return
}

// WriteFileSyncF hace la función fs.writeFileSync() de JavaScript.
//
//export WriteFileSyncF
func WriteFileSyncF(context C.JSContextRef, function C.JSObjectRef, thisObject C.JSObjectRef, argumentCount C.size_t, arguments *C.JSValueRef, exception *C.JSValueRef) C.JSValueRef {
	if int(argumentCount) < 2 {
		return C.JSValueMakeUndefined(context)
	}
	argumentSlice := (*[1 << 30]C.JSValueRef)(unsafe.Pointer(arguments))[:argumentCount:argumentCount]
	str := C.JSValueToStringCopy(context, argumentSlice[0], nil)
	bufferSize := C.JSStringGetMaximumUTF8CStringSize(str)
	buffer := C.malloc(bufferSize)
	C.JSStringGetUTF8CString(str, (*C.char)(buffer), bufferSize)
	fileName := C.GoString((*C.char)(buffer))

	contentStr := C.JSValueToStringCopy(context, argumentSlice[1], nil)
	bufferSizeContent := C.JSStringGetMaximumUTF8CStringSize(contentStr)
	bufferContent := C.malloc(bufferSizeContent)
	C.JSStringGetUTF8CString(contentStr, (*C.char)(bufferContent), bufferSizeContent)
	content := C.GoString((*C.char)(bufferContent))

	err := os.WriteFile(fileName, []byte(content), 06444)
	if err != nil {
		return C.JSValueMakeNull(context)
	}

	return C.JSValueMakeUndefined(context)
}

// ReadDirSyncF representa la función fs.readdirSync() de JavaScript.
//
//export ReadDirSyncF
func ReadDirSyncF(context C.JSContextRef, function C.JSObjectRef, thisObject C.JSObjectRef, argumentCount C.size_t, arguments *C.JSValueRef, exception *C.JSValueRef) C.JSValueRef {
	if int(argumentCount) < 1 {
		return C.JSValueMakeUndefined(context)
	}

	argumentSlice := (*[1 << 30]C.JSValueRef)(unsafe.Pointer(arguments))[:argumentCount:argumentCount]
	str := C.JSValueToStringCopy(context, argumentSlice[0], nil)
	bufferSize := C.JSStringGetMaximumUTF8CStringSize(str)
	buffer := C.malloc(bufferSize)
	C.JSStringGetUTF8CString(str, (*C.char)(buffer), bufferSize)
	dirName := C.GoString((*C.char)(buffer))

	files, err := os.ReadDir(dirName)
	C.free(unsafe.Pointer(buffer))
	if err != nil {
		return C.JSValueMakeUndefined(context)
	}
	newFiles := make([]C.JSValueRef, len(files))
	for index, item := range files {
		itemCString := C.CString(item.Name())
		itemJSOpaqueString := C.JSStringCreateWithUTF8CString(itemCString)
		C.free(unsafe.Pointer(itemCString))
		itemValueString := C.JSValueMakeString(context, itemJSOpaqueString)
		C.JSStringRelease(itemJSOpaqueString)
		newFiles[index] = itemValueString
	}
	objectArray := C.JSObjectMakeArray(context, C.ulong(len(newFiles)), &newFiles[0], exception)
	return (C.JSValueRef)(objectArray)
}

// MkDirSyncF representa la función fs.mkdirSync() de JavaScript.
//
//export MkDirSyncF
func MkDirSyncF(context C.JSContextRef, function C.JSObjectRef, thisObject C.JSObjectRef, argumentCount C.size_t, arguments *C.JSValueRef, exception *C.JSValueRef) C.JSValueRef {
	if int(argumentCount) < 1 {
		return C.JSValueMakeUndefined(context)
	}

	argumentSlice := (*[1 << 30]C.JSValueRef)(unsafe.Pointer(arguments))[:argumentCount:argumentCount]
	str := C.JSValueToStringCopy(context, argumentSlice[0], nil)
	bufferSize := C.JSStringGetMaximumUTF8CStringSize(str)
	buffer := C.malloc(bufferSize)
	C.JSStringGetUTF8CString(str, (*C.char)(buffer), bufferSize)
	dirName := C.GoString((*C.char)(buffer))

	err := os.Mkdir(dirName, 0755)
	C.free(unsafe.Pointer(buffer))
	if err != nil {
		return C.JSValueMakeUndefined(context)
	}
	return C.JSValueMakeUndefined(context)
}

// ReadFileF representa la función fs.readFile() de JavaScript.
//
//export ReadFileF
func ReadFileF(context C.JSContextRef, function C.JSObjectRef, thisObject C.JSObjectRef, argumentCount C.size_t, arguments *C.JSValueRef, exception *C.JSValueRef) (finalValue C.JSValueRef) {
	if int(argumentCount) < 3 {
		return C.JSValueMakeUndefined(context)
	}
	argumentSlice := (*[1 << 30]C.JSValueRef)(unsafe.Pointer(arguments))[:argumentCount:argumentCount]
	str := C.JSValueToStringCopy(context, argumentSlice[0], nil)
	bufferSize := C.JSStringGetMaximumUTF8CStringSize(str)
	buffer := C.malloc(bufferSize)
	C.JSStringGetUTF8CString(str, (*C.char)(buffer), bufferSize)
	fileName := C.GoString((*C.char)(buffer))
	C.free(unsafe.Pointer(buffer))
	strEncoding := C.JSValueToStringCopy(context, argumentSlice[1], nil)
	bufferSizeEncoding := C.JSStringGetMaximumUTF8CStringSize(strEncoding)
	bufferEncoding := C.malloc(bufferSizeEncoding)
	C.JSStringGetUTF8CString(strEncoding, (*C.char)(bufferEncoding), bufferSizeEncoding)
	encoding := C.GoString((*C.char)(bufferEncoding))
	C.free(unsafe.Pointer(bufferEncoding))
	finalValue = C.JSValueMakeUndefined(context)
	functionObject := C.JSValueToObject(context, argumentSlice[2], exception)
	switch encoding {
	case "utf8", "utf-8":
		file, err := os.ReadFile(fileName)
		if err != nil {
			errorCString := C.CString(err.Error())
			errorJSString := C.JSStringCreateWithUTF8CString(errorCString)
			C.free(unsafe.Pointer(errorCString))
			errorJSStringValue := C.JSValueMakeString(context, errorJSString)
			C.JSStringRelease(errorJSString)
			nullData := C.JSValueMakeNull(context)
			arguments := []C.JSValueRef{errorJSStringValue, nullData}
			C.JSObjectCallAsFunction(context, functionObject, thisObject, 2, &arguments[0], exception)
			return
		}
		fileStringC := C.CString(string(file))
		fileStringJS := C.JSStringCreateWithUTF8CString(fileStringC)
		C.free(unsafe.Pointer(fileStringC))
		fileStringValue := C.JSValueMakeString(context, fileStringJS)
		C.JSStringRelease(fileStringJS)
		nullError := C.JSValueMakeNull(context)
		arguments := []C.JSValueRef{nullError, fileStringValue}
		C.JSObjectCallAsFunction(context, functionObject, thisObject, 2, &arguments[0], exception)
	case "base64":
		file, err := os.ReadFile(fileName)
		if err != nil {
			errorCString := C.CString(err.Error())
			errorJSString := C.JSStringCreateWithUTF8CString(errorCString)
			C.free(unsafe.Pointer(errorCString))
			errorJSStringValue := C.JSValueMakeString(context, errorJSString)
			C.JSStringRelease(errorJSString)
			nullData := C.JSValueMakeNull(context)
			arguments := []C.JSValueRef{errorJSStringValue, nullData}
			C.JSObjectCallAsFunction(context, functionObject, thisObject, 2, &arguments[0], exception)
			return
		}
		fileStringC := C.CString(base64.StdEncoding.EncodeToString(file))
		fileStringJS := C.JSStringCreateWithUTF8CString(fileStringC)
		C.free(unsafe.Pointer(fileStringC))
		fileStringValue := C.JSValueMakeString(context, fileStringJS)
		C.JSStringRelease(fileStringJS)
		nullError := C.JSValueMakeNull(context)
		arguments := []C.JSValueRef{nullError, fileStringValue}
		C.JSObjectCallAsFunction(context, functionObject, thisObject, 2, &arguments[0], exception)
	case "hex":
		file, err := os.ReadFile(fileName)
		if err != nil {
			errorCString := C.CString(err.Error())
			errorJSString := C.JSStringCreateWithUTF8CString(errorCString)
			C.free(unsafe.Pointer(errorCString))
			errorJSStringValue := C.JSValueMakeString(context, errorJSString)
			C.JSStringRelease(errorJSString)
			nullData := C.JSValueMakeNull(context)
			arguments := []C.JSValueRef{errorJSStringValue, nullData}
			C.JSObjectCallAsFunction(context, functionObject, thisObject, 2, &arguments[0], exception)
			return
		}
		fileStringC := C.CString(hex.EncodeToString(file))
		fileStringJS := C.JSStringCreateWithUTF8CString(fileStringC)
		C.free(unsafe.Pointer(fileStringC))
		fileStringValue := C.JSValueMakeString(context, fileStringJS)
		C.JSStringRelease(fileStringJS)
		nullError := C.JSValueMakeNull(context)
		arguments := []C.JSValueRef{nullError, fileStringValue}
		C.JSObjectCallAsFunction(context, functionObject, thisObject, 2, &arguments[0], exception)
	}
	return
}

// ReadFile devuelve la función callback de JavaScript en C para la función fs.readFile() en JavaScript.
func ReadFile() C.JSObjectCallAsFunctionCallback {
	return C.JSObjectCallAsFunctionCallback(C.ReadFileF)
}

// ReadFileSync devuelve la función callback de JavaScript en C para la función readFileSync en JavaScript.
func ReadFileSync() C.JSObjectCallAsFunctionCallback {
	return C.JSObjectCallAsFunctionCallback(C.ReadFileSyncF)
}

// WriteFileSync devuelve la función callback de JavaScript en C para la función fs.writeFileSync() en JavaScript.
func WriteFileSync() C.JSObjectCallAsFunctionCallback {
	return C.JSObjectCallAsFunctionCallback(C.WriteFileSyncF)
}

// ReadDirSync devuelve la función callback de JavaScript en C para la función fs.readdirSync() en JavaScript.
func ReadDirSync() C.JSObjectCallAsFunctionCallback {
	return C.JSObjectCallAsFunctionCallback(C.ReadDirSyncF)
}

// MkDirSync devuelve la función callback de JavaScript en C para la función fs.mkdirSync() en JavaScript.
func MkDirSync() C.JSObjectCallAsFunctionCallback {
	return C.JSObjectCallAsFunctionCallback(C.MkDirSyncF)
}
