package http

/*
#cgo CFLAGS: -I/usr/include/webkitgtk-4.0
#cgo LDFLAGS: -ljavascriptcoregtk-4.0
#include <JavaScriptCore/JavaScript.h>
#include <stdlib.h>

extern JSValueRef GetF(JSContextRef context, JSObjectRef function, JSObjectRef thisObject, size_t argumentCount, JSValueRef arguments[], JSValueRef* exception);
extern JSValueRef ServeF(JSContextRef context, JSObjectRef function, JSObjectRef thisObject, size_t argumentCount, JSValueRef arguments[], JSValueRef* exception);
*/
import "C"
import (
	"fmt"
	"log"
	"net/http"
	"unsafe"
)

//export GetF
func GetF(context C.JSContextRef, function C.JSObjectRef, thisObject C.JSObjectRef, argumentCount C.size_t, arguments *C.JSValueRef, exception *C.JSValueRef) C.JSValueRef {
	argumentSlice := (*[1 << 30]C.JSValueRef)(unsafe.Pointer(arguments))[:argumentCount:argumentCount]
	if len(argumentSlice) <= 0 {
		return C.JSValueMakeUndefined(context)
	}
	routeStr := C.JSValueToStringCopy(context, argumentSlice[0], exception)
	bufferSize := C.JSStringGetMaximumUTF8CStringSize(routeStr)
	buffer := C.malloc(bufferSize)
	C.JSStringGetUTF8CString(routeStr, (*C.char)(buffer), bufferSize)
	routeGoStr := C.GoString((*C.char)(buffer))
	C.free(unsafe.Pointer(buffer))
	functionObject := C.JSValueToObject(context, argumentSlice[1], exception)
	if C.JSValueIsUndefined(context, functionObject) {
		log.Fatal("There isn´t an object in the function \"get\"")
		return C.JSValueMakeUndefined(context)
	}
	http.HandleFunc(routeGoStr, func(w http.ResponseWriter, r *http.Request) {
		response := C.JSObjectCallAsFunction(context, functionObject, thisObject, 0, nil, exception)
		if !C.JSValueIsString(context, response) {
			fmt.Fprintf(w, "Hi!")
			return
		}
		responseStr := C.JSValueToStringCopy(context, response, exception)
		responseBufferSize := C.JSStringGetMaximumUTF8CStringSize(routeStr)
		responseBuffer := C.malloc(responseBufferSize)
		C.JSStringGetUTF8CString(responseStr, (*C.char)(responseBuffer), responseBufferSize)
		responseValue := C.GoString((*C.char)(responseBuffer))
		C.free(unsafe.Pointer(responseBuffer))
		fmt.Fprintf(w, responseValue)
	})
	return C.JSValueMakeUndefined(context)
}

//export ServeF
func ServeF(context C.JSContextRef, function C.JSObjectRef, thisObject C.JSObjectRef, argumentCount C.size_t, arguments *C.JSValueRef, exception *C.JSValueRef) (finalValue C.JSValueRef) {
	argumentSlice := (*[1 << 30]C.JSValueRef)(unsafe.Pointer(arguments))[:argumentCount:argumentCount]
	if len(argumentSlice) <= 0 {
		return C.JSValueMakeUndefined(context)
	}
	objectRef := C.JSValueToObject(context, argumentSlice[0], exception)
	if C.JSValueIsUndefined(context, objectRef) {
		fmt.Println("No hay un objeto para la función serve()")
		return C.JSValueMakeUndefined(context)
	}
	cPortProperty := C.CString("port")
	cPortR := C.JSStringCreateWithUTF8CString(cPortProperty)
	portPropertyValue := C.JSObjectGetProperty(context, objectRef, cPortR, exception)
	C.free(unsafe.Pointer(cPortProperty))
	C.JSStringRelease(cPortR)
	portPropertyValueStringDefinition := C.JSValueToStringCopy(context, portPropertyValue, nil)
	portPropertyValueBufferSize := C.JSStringGetMaximumUTF8CStringSize(portPropertyValueStringDefinition)
	portPropertyValueBuffer := C.malloc(portPropertyValueBufferSize)
	C.JSStringGetUTF8CString(portPropertyValueStringDefinition, (*C.char)(portPropertyValueBuffer), portPropertyValueBufferSize)
	portValue := C.GoString((*C.char)(portPropertyValueBuffer))
	C.free(unsafe.Pointer(portPropertyValueBuffer))
	http.ListenAndServe(fmt.Sprintf(":%s", portValue), nil)
	finalValue = C.JSValueMakeUndefined(context)
	return
}

// Get devuelve la función callback de C para la función get en JavaScript.
func Get() C.JSObjectCallAsFunctionCallback { return C.JSObjectCallAsFunctionCallback(C.GetF) }

// Serve devuelve la función callback de C para la función serve en JavaScript.
func Serve() C.JSObjectCallAsFunctionCallback { return C.JSObjectCallAsFunctionCallback(C.ServeF) }
