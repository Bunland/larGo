package require

/*
#cgo CFLAGS: -I/usr/include/webkitgtk-4.0
#cgo LDFLAGS: -ljavascriptcoregtk-4.0
#include <JavaScriptCore/JavaScript.h>
#include <stdlib.h>

extern JSValueRef RequireF(JSContextRef context, JSObjectRef function, JSObjectRef thisObject, size_t argumentCount, JSValueRef arguments[], JSValueRef* exception);
*/
import "C"
import (
	"largo/src/fs"
	"largo/src/http"
	"largo/src/modules"
	"largo/src/utils"
	"os"
	"path/filepath"
	"strings"
	"unsafe"
)

func itemExists(array []string, item string) (value bool) {
	value = false
	for _, str := range array {
		if str == item {
			value = true
		}
	}
	return
}

func fileExists(path string) (value bool) {
	value = true
	if _, err := os.Stat(path); os.IsNotExist(err) {
		value = false
	}
	return
}

// createCustomFunction crea una función JavaScript personalizada y la establece como propiedad del objeto global.
func createCustomFunction(context C.JSContextRef, globalObject C.JSObjectRef, functionName string, functionCallback C.JSObjectCallAsFunctionCallback) {
	// Crear el string de C
	functionNameC := C.CString(functionName)

	// Crear una cadena JavaScript a partir del nombre de la función en formato UTF-8.
	functionString := C.JSStringCreateWithUTF8CString(functionNameC)

	// Liberar el string de C
	C.free(unsafe.Pointer(functionNameC))

	// Crear un objeto de función JavaScript usando la cadena y la devolución de llamada de la función.
	functionObject := C.JSObjectMakeFunctionWithCallback(context, functionString, functionCallback)

	// Establecer la función recién creada como propiedad del objeto global.
	C.JSObjectSetProperty(context, globalObject, functionString, functionObject, C.kJSPropertyAttributeNone, nil)

	// Liberar la cadena de función creada con JSStringCreateWithUTF8CString para evitar fugas de memoria.
	C.JSStringRelease(functionString)
}

func SetDirnameAndFilename(context C.JSContextRef, globalObject C.JSObjectRef, dirname string, filename string) {
	dirname_str := C.CString(dirname)
	dirname_js := C.JSStringCreateWithUTF8CString(dirname_str)
	C.free(unsafe.Pointer(dirname_str))
	dirnameString := C.CString("__dirname")
	dirnameStringJS := C.JSStringCreateWithUTF8CString(dirnameString)
	dirnameValue := C.CString(dirname)
	dirnameValueJS := C.JSStringCreateWithUTF8CString(dirnameValue)
	dirnameValueJSString := C.JSValueMakeString(context, dirnameValueJS)
	C.JSObjectSetProperty(context, globalObject, dirnameStringJS, dirnameValueJSString, C.kJSPropertyAttributeNone, nil)
	C.free(unsafe.Pointer(dirnameString))
	C.free(unsafe.Pointer(dirnameValue))
	C.JSStringRelease(dirnameStringJS)
	C.JSStringRelease(dirnameValueJS)
	C.JSStringRelease(dirname_js)
	filename_str := C.CString(filename)
	filename_js := C.JSStringCreateWithUTF8CString(filename_str)
	C.free(unsafe.Pointer(filename_str))
	filenameString := C.CString("__filename")
	filenameStringJS := C.JSStringCreateWithUTF8CString(filenameString)
	filenameValue := C.CString(filename)
	filenameValueJS := C.JSStringCreateWithUTF8CString(filenameValue)
	filenameValueJSString := C.JSValueMakeString(context, filenameValueJS)
	C.JSObjectSetProperty(context, globalObject, filenameStringJS, filenameValueJSString, C.kJSPropertyAttributeNone, nil)
	C.free(unsafe.Pointer(filenameString))
	C.free(unsafe.Pointer(filenameValue))
	C.JSStringRelease(filenameStringJS)
	C.JSStringRelease(filenameValueJS)
	C.JSStringRelease(filename_js)
}

// Require es la implementación de require de JavaScript
//
//export RequireF
func RequireF(context C.JSContextRef, function C.JSObjectRef, thisObject C.JSObjectRef, argumentCount C.size_t, arguments *C.JSValueRef, exception *C.JSValueRef) (finalValue C.JSValueRef) {
	argumentSlice := (*[1 << 30]C.JSValueRef)(unsafe.Pointer(arguments))[:argumentCount:argumentCount]
	// dirnameCString := C.CString("__dirname")
	// dirnameJS := C.JSStringCreateWithUTF8CString(dirnameCString)
	// C.free(unsafe.Pointer(dirnameCString))
	// filenameCString := C.CString("__filename")
	// filenameJS := C.JSStringCreateWithUTF8CString(filenameCString)
	// C.free(unsafe.Pointer(filenameCString))
	// dirnameStr := C.JSValueToStringCopy(context, C.JSObjectGetProperty(context, C.JSContextGetGlobalObject(context), dirnameJS, nil), exception)
	// filenameStr := C.JSValueToStringCopy(context, C.JSObjectGetProperty(context, C.JSContextGetGlobalObject(context), filenameJS, nil), exception)
	// dirnameSizeStr := C.JSStringGetMaximumUTF8CStringSize(dirnameStr)
	// filenameSizeStr := C.JSStringGetMaximumUTF8CStringSize(filenameStr)
	// dirnameOldValueBuffer := C.malloc(dirnameSizeStr)
	// filenameOldValueBuffer := C.malloc(filenameSizeStr)
	// C.JSStringGetUTF8CString(dirnameStr, (*C.char)(dirnameOldValueBuffer), dirnameSizeStr)
	// C.JSStringGetUTF8CString(filenameStr, (*C.char)(filenameOldValueBuffer), filenameSizeStr)
	// filenameOldValue := C.GoString((*C.char)(filenameOldValueBuffer))
	// dirnameOldValue := C.GoString((*C.char)(dirnameOldValueBuffer))
	if len(argumentSlice) <= 0 {
		return C.JSValueMakeUndefined(context)
	}
	moduleName := C.JSValueToStringCopy(context, argumentSlice[0], nil)
	bufferSize := C.JSStringGetMaximumUTF8CStringSize(moduleName)
	buffer := C.malloc(bufferSize)
	C.JSStringGetUTF8CString(moduleName, (*C.char)(buffer), bufferSize)
	module := C.GoString((*C.char)(buffer))
	C.free(unsafe.Pointer(buffer))

	if ok := itemExists(modules.Commands, module); ok {
		apiObject := C.JSObjectMake(context, nil, nil)
		switch module {
		case "fs", "node:fs":
			createCustomFunction(context, apiObject, "readFileSync", C.JSObjectCallAsFunctionCallback(fs.ReadFileSync()))
			createCustomFunction(context, apiObject, "writeFileSync", C.JSObjectCallAsFunctionCallback(fs.WriteFileSync()))
			createCustomFunction(context, apiObject, "readdirSync", C.JSObjectCallAsFunctionCallback(fs.ReadDirSync()))
			createCustomFunction(context, apiObject, "mkdirSync", C.JSObjectCallAsFunctionCallback(fs.MkDirSync()))
			createCustomFunction(context, apiObject, "readFile", C.JSObjectCallAsFunctionCallback(fs.ReadFile()))
		case "lar:http":
			createCustomFunction(context, apiObject, "get", C.JSObjectCallAsFunctionCallback(http.Get()))
			createCustomFunction(context, apiObject, "fetch", C.JSObjectCallAsFunctionCallback(http.Fetch()))
			createCustomFunction(context, apiObject, "serve", C.JSObjectCallAsFunctionCallback(http.Serve()))
		}
		finalValue = (C.JSValueRef)(apiObject)
	} else {
		var content string
		if strings.HasSuffix(module, ".js") || strings.HasSuffix(module, ".ts") {
			var b bool = false
			if strings.HasSuffix(module, ".js") && !fileExists(module+".ts") && !fileExists(module+".js") {
				moduleWithoutSuffix, _ := strings.CutSuffix(module, ".js")
				content = utils.ReadFile(moduleWithoutSuffix + ".ts")
				module = moduleWithoutSuffix + ".ts"
				b = true
			}
			if !b {
				content = utils.ReadFile(module)
			}
		} else if strings.HasSuffix(module, ".json") {
			content = utils.ReadFile(module)
			jsonCString := C.CString(content)
			jsonString := C.JSStringCreateWithUTF8CString(jsonCString)
			C.free(unsafe.Pointer(jsonCString))
			jsonValueWithContent := C.JSValueMakeFromJSONString(context, jsonString)
			C.JSStringRelease(jsonString)
			return jsonValueWithContent
		} else {
			if fileExists(module + ".js") {
				content = utils.ReadFile(module + ".js")
				module = module + ".js"
			} else if fileExists(module + ".ts") {
				content = utils.ReadFile(module + ".ts")
				module = module + ".ts"
			} else {
				return C.JSValueMakeUndefined(context)
			}
		}

		if content == "" {
			return C.JSValueMakeUndefined(context)
		}

		moduleC := C.CString("module")
		moduleNameR := C.JSStringCreateWithUTF8CString(moduleC)
		C.free(unsafe.Pointer(moduleC))
		moduleObject := C.JSObjectMake(context, nil, nil)
		C.JSObjectSetProperty(context, thisObject, moduleNameR, moduleObject, C.kJSPropertyAttributeNone, nil)
		contentC := C.CString(content)
		script := C.JSStringCreateWithUTF8CString(contentC)
		absolutePath, _ := filepath.Abs(module)
		dirname := absolutePath[0 : len(absolutePath)-(len(filepath.Base(module))+1)]
		SetDirnameAndFilename(context, C.JSContextGetGlobalObject(context), dirname, absolutePath)
		C.JSEvaluateScript(context, script, nil, nil, 0, exception)

		C.JSStringRelease(script)
		C.free(unsafe.Pointer(contentC))

		moduleValueScript := C.JSObjectGetProperty(context, C.JSContextGetGlobalObject(context), moduleNameR, nil)
		C.JSStringRelease(moduleNameR)
		moduleObjectScript := C.JSValueToObject(context, moduleValueScript, exception)

		exportsName := C.CString("exports")
		exportsNameJSC := C.JSStringCreateWithUTF8CString(exportsName)
		finalValue = C.JSObjectGetProperty(context, moduleObjectScript, exportsNameJSC, nil)
		C.free(unsafe.Pointer(exportsName))
		C.JSStringRelease(exportsNameJSC)
		// SetDirnameAndFilename(context, C.JSContextGetGlobalObject(context), dirnameOldValue, filenameOldValue)
	}

	return
}

// Require devuelve la función de callback de C para la función require en JavaScript.
func Require() C.JSObjectCallAsFunctionCallback {
	return C.JSObjectCallAsFunctionCallback(C.RequireF)
}
