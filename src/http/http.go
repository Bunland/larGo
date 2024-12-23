package http

/*
#cgo CFLAGS: -I/usr/include/webkitgtk-4.0
#cgo LDFLAGS: -ljavascriptcoregtk-4.0
#include <JavaScriptCore/JavaScript.h>
#include <stdlib.h>

extern JSValueRef GetF(JSContextRef context, JSObjectRef function, JSObjectRef thisObject, size_t argumentCount, JSValueRef arguments[], JSValueRef* exception);
extern JSValueRef PostF(JSContextRef context, JSObjectRef function, JSObjectRef thisObject, size_t argumentCount, JSValueRef arguments[], JSValueRef* exception);
extern JSValueRef ServeF(JSContextRef context, JSObjectRef function, JSObjectRef thisObject, size_t argumentCount, JSValueRef arguments[], JSValueRef* exception);
extern JSValueRef FetchF(JSContextRef context, JSObjectRef function, JSObjectRef thisObject, size_t argumentCount, JSValueRef arguments[], JSValueRef* exception);
*/
import "C"
import (
	"fmt"
	"io"
	"net/http"
	"unsafe"
)

//export GetF
func GetF(context C.JSContextRef, function C.JSObjectRef, thisObject C.JSObjectRef, argumentCount C.size_t, arguments *C.JSValueRef, exception *C.JSValueRef) C.JSValueRef {
	// Validar argumentos
	argumentSlice := (*[1 << 30]C.JSValueRef)(unsafe.Pointer(arguments))[:argumentCount:argumentCount]
	if len(argumentSlice) <= 0 {
		return C.JSValueMakeUndefined(context)
	}

	// Convertir la ruta
	routeStr := C.JSValueToStringCopy(context, argumentSlice[0], exception)
	if routeStr == nil {
		return C.JSValueMakeUndefined(context)
	}
	defer C.JSStringRelease(routeStr)

	bufferSize := C.JSStringGetMaximumUTF8CStringSize(routeStr)
	buffer := C.malloc(bufferSize)
	if buffer == nil {
		return C.JSValueMakeUndefined(context)
	}
	defer C.free(unsafe.Pointer(buffer))

	C.JSStringGetUTF8CString(routeStr, (*C.char)(buffer), bufferSize)
	routeGoStr := C.GoString((*C.char)(buffer))

	// Proteger el contexto y la función
	functionObject := C.JSValueToObject(context, argumentSlice[1], exception)
	if C.JSValueIsUndefined(context, functionObject) {
		return C.JSValueMakeUndefined(context)
	}

	// Proteger las referencias JavaScript
	C.JSValueProtect(context, C.JSValueRef(functionObject))
	defer C.JSValueUnprotect(context, C.JSValueRef(functionObject))

	http.HandleFunc(routeGoStr, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "GET" {
			http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
			return
		}

		response := C.JSObjectCallAsFunction(context, functionObject, thisObject, 0, nil, exception)
		if response == nil || !C.JSValueIsString(context, response) {
			fmt.Fprint(w, "")
			return
		}

		responseStr := C.JSValueToStringCopy(context, response, exception)
		if responseStr == nil {
			fmt.Fprint(w, "")
			return
		}
		defer C.JSStringRelease(responseStr)

		responseBufferSize := C.JSStringGetMaximumUTF8CStringSize(responseStr)
		responseBuffer := C.malloc(responseBufferSize)
		if responseBuffer == nil {
			fmt.Fprint(w, "")
			return
		}
		defer C.free(unsafe.Pointer(responseBuffer))

		C.JSStringGetUTF8CString(responseStr, (*C.char)(responseBuffer), responseBufferSize)
		responseValue := C.GoString((*C.char)(responseBuffer))
		fmt.Fprint(w, responseValue)
	})

	return C.JSValueMakeUndefined(context)
}

//export PostF
func PostF(context C.JSContextRef, function C.JSObjectRef, thisObject C.JSObjectRef, argumentCount C.size_t, arguments *C.JSValueRef, exception *C.JSValueRef) C.JSValueRef {
	// Validar argumentos
	argumentSlice := (*[1 << 30]C.JSValueRef)(unsafe.Pointer(arguments))[:argumentCount:argumentCount]
	if len(argumentSlice) <= 0 {
		return C.JSValueMakeUndefined(context)
	}

	// Convertir la ruta
	routeStr := C.JSValueToStringCopy(context, argumentSlice[0], exception)
	if routeStr == nil {
		return C.JSValueMakeUndefined(context)
	}
	defer C.JSStringRelease(routeStr)

	bufferSize := C.JSStringGetMaximumUTF8CStringSize(routeStr)
	buffer := C.malloc(bufferSize)
	if buffer == nil {
		return C.JSValueMakeUndefined(context)
	}
	defer C.free(unsafe.Pointer(buffer))

	C.JSStringGetUTF8CString(routeStr, (*C.char)(buffer), bufferSize)
	routeGoStr := C.GoString((*C.char)(buffer))

	// Proteger el contexto y la función
	functionObject := C.JSValueToObject(context, argumentSlice[1], exception)
	if C.JSValueIsUndefined(context, functionObject) {
		return C.JSValueMakeUndefined(context)
	}

	// Proteger las referencias JavaScript
	C.JSValueProtect(context, C.JSValueRef(functionObject))
	defer C.JSValueUnprotect(context, C.JSValueRef(functionObject))

	http.HandleFunc(routeGoStr, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "POST" {
			http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
			return
		}

		response := C.JSObjectCallAsFunction(context, functionObject, thisObject, 0, nil, exception)
		if response == nil || !C.JSValueIsString(context, response) {
			fmt.Fprint(w, "")
			return
		}

		responseStr := C.JSValueToStringCopy(context, response, exception)
		if responseStr == nil {
			fmt.Fprint(w, "")
			return
		}
		defer C.JSStringRelease(responseStr)

		responseBufferSize := C.JSStringGetMaximumUTF8CStringSize(responseStr)
		responseBuffer := C.malloc(responseBufferSize)
		if responseBuffer == nil {
			fmt.Fprint(w, "")
			return
		}
		defer C.free(unsafe.Pointer(responseBuffer))

		C.JSStringGetUTF8CString(responseStr, (*C.char)(responseBuffer), responseBufferSize)
		responseValue := C.GoString((*C.char)(responseBuffer))
		fmt.Fprint(w, responseValue)
	})

	return C.JSValueMakeUndefined(context)
}

//export FetchF
func FetchF(context C.JSContextRef, function C.JSObjectRef, thisObject C.JSObjectRef, argumentCount C.size_t, arguments *C.JSValueRef, exception *C.JSValueRef) C.JSValueRef {
	argumentSlice := (*[1 << 30]C.JSValueRef)(unsafe.Pointer(arguments))[:argumentCount:argumentCount]
	if len(argumentSlice) <= 0 {
		return C.JSValueMakeUndefined(context)
	}

	urlStr := C.JSValueToStringCopy(context, argumentSlice[0], exception)
	if urlStr == nil {
		return C.JSValueMakeUndefined(context)
	}
	defer C.JSStringRelease(urlStr)

	bufferSize := C.JSStringGetMaximumUTF8CStringSize(urlStr)
	buffer := C.malloc(bufferSize)
	if buffer == nil {
		return C.JSValueMakeUndefined(context)
	}
	defer C.free(unsafe.Pointer(buffer))

	C.JSStringGetUTF8CString(urlStr, (*C.char)(buffer), bufferSize)
	urlGoStr := C.GoString((*C.char)(buffer))

	response, err := http.Get(urlGoStr)
	if err != nil {
		return C.JSValueMakeUndefined(context)
	}
	defer response.Body.Close()

	body, err := io.ReadAll(response.Body)
	if err != nil {
		return C.JSValueMakeUndefined(context)
	}

	responseStr := string(body)
	cResponseStr := C.CString(responseStr)
	if cResponseStr == nil {
		return C.JSValueMakeUndefined(context)
	}
	defer C.free(unsafe.Pointer(cResponseStr))

	responseJS := C.JSStringCreateWithUTF8CString(cResponseStr)
	if responseJS == nil {
		return C.JSValueMakeUndefined(context)
	}
	defer C.JSStringRelease(responseJS)

	return C.JSValueMakeString(context, responseJS)
}

//export ServeF
func ServeF(context C.JSContextRef, function C.JSObjectRef, thisObject C.JSObjectRef, argumentCount C.size_t, arguments *C.JSValueRef, exception *C.JSValueRef) C.JSValueRef {
	argumentSlice := (*[1 << 30]C.JSValueRef)(unsafe.Pointer(arguments))[:argumentCount:argumentCount]
	if len(argumentSlice) <= 0 {
		return C.JSValueMakeUndefined(context)
	}

	objectRef := C.JSValueToObject(context, argumentSlice[0], exception)
	if C.JSValueIsUndefined(context, objectRef) {
		return C.JSValueMakeUndefined(context)
	}

	cPortProperty := C.CString("port")
	if cPortProperty == nil {
		return C.JSValueMakeUndefined(context)
	}
	defer C.free(unsafe.Pointer(cPortProperty))

	cPortR := C.JSStringCreateWithUTF8CString(cPortProperty)
	if cPortR == nil {
		return C.JSValueMakeUndefined(context)
	}
	defer C.JSStringRelease(cPortR)

	portPropertyValue := C.JSObjectGetProperty(context, objectRef, cPortR, exception)
	if C.JSValueIsUndefined(context, portPropertyValue) {
		return C.JSValueMakeUndefined(context)
	}

	port := int(C.JSValueToNumber(context, portPropertyValue, exception))
	if err := http.ListenAndServe(fmt.Sprintf(":%d", port), nil); err != nil {
		return C.JSValueMakeUndefined(context)
	}

	return C.JSValueMakeUndefined(context)
}

// Get devuelve la función callback de C para la función get en JavaScript.
func Get() C.JSObjectCallAsFunctionCallback { return C.JSObjectCallAsFunctionCallback(C.GetF) }

// Serve devuelve la función callback de C para la función serve en JavaScript.
func Serve() C.JSObjectCallAsFunctionCallback { return C.JSObjectCallAsFunctionCallback(C.ServeF) }

// Fetch devuelve la función callback de C para la función fetch en JavaScript.
func Fetch() C.JSObjectCallAsFunctionCallback { return C.JSObjectCallAsFunctionCallback(C.FetchF) }

// Post devuelve la función callback de C para la función post en JavaScript.
func Post() C.JSObjectCallAsFunctionCallback { return C.JSObjectCallAsFunctionCallback(C.PostF) }
