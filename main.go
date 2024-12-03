package main

/*
#cgo CFLAGS: -I/usr/include/webkitgtk-4.0
#cgo LDFLAGS: -ljavascriptcoregtk-4.0
#include <JavaScriptCore/JavaScript.h>
#include <stdlib.h>
*/
import "C"
import (
	"bufio"
	"fmt"
	"largo/src/console"
	"largo/src/http"
	"largo/src/math"
	"largo/src/modules"
	"largo/src/require"
	"largo/src/utils"
	"os"
	"os/signal"
	"path/filepath"
	"unsafe"

	"github.com/evanw/esbuild/pkg/api"

	"github.com/fatih/color"
)

// createCustomFunction crea una función JavaScript personalizada y la establece como propiedad del objeto global.
func createCustomFunction(context C.JSGlobalContextRef, globalObject C.JSObjectRef, functionName string, functionCallback C.JSObjectCallAsFunctionCallback) {
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

// Apis define las API disponibles en JavaScript.
func Apis(context C.JSGlobalContextRef, globalObject C.JSObjectRef) {
	console_str := C.CString("console")
	console_js := C.JSStringCreateWithUTF8CString(console_str)
	C.free(unsafe.Pointer(console_str))
	consoleGlobalObject := C.JSObjectMake(context, nil, nil)
	C.JSObjectSetProperty(context, globalObject, console_js, consoleGlobalObject, C.kJSPropertyAttributeNone, nil)
	C.JSStringRelease(console_js)
	lar_str := C.CString("Lar")
	lar_js := C.JSStringCreateWithUTF8CString(lar_str)
	C.free(unsafe.Pointer(lar_str))
	larGlobalObject := C.JSObjectMake(context, nil, nil)
	versionString := C.CString("version")
	versionStringJS := C.JSStringCreateWithUTF8CString(versionString)
	versionValue := C.CString("0.0.1")
	versionValueJS := C.JSStringCreateWithUTF8CString(versionValue)
	versionValueJSString := C.JSValueMakeString(context, versionValueJS)
	C.JSObjectSetProperty(context, larGlobalObject, versionStringJS, versionValueJSString, C.kJSPropertyAttributeNone, nil)
	C.free(unsafe.Pointer(versionString))
	C.free(unsafe.Pointer(versionValue))
	C.JSStringRelease(versionStringJS)
	C.JSStringRelease(versionValueJS)
	C.JSObjectSetProperty(context, globalObject, lar_js, larGlobalObject, C.kJSPropertyAttributeNone, nil)
	C.JSStringRelease(lar_js)
	process_str := C.CString("process")
	process_js := C.JSStringCreateWithUTF8CString(process_str)
	C.free(unsafe.Pointer(process_str))
	processGlobalObject := C.JSObjectMake(context, nil, nil)
	processVersionString := C.CString("version")
	processVersionStringJS := C.JSStringCreateWithUTF8CString(processVersionString)
	processVersionValue := C.CString("20.11.0")
	processVersionValueJS := C.JSStringCreateWithUTF8CString(processVersionValue)
	processVersionValueJSString := C.JSValueMakeString(context, processVersionValueJS)
	C.JSObjectSetProperty(context, processGlobalObject, processVersionStringJS, processVersionValueJSString, C.kJSPropertyAttributeNone, nil)
	C.free(unsafe.Pointer(processVersionString))
	C.free(unsafe.Pointer(processVersionValue))
	C.JSStringRelease(processVersionStringJS)
	C.JSStringRelease(processVersionValueJS)
	C.JSObjectSetProperty(context, globalObject, process_js, processGlobalObject, C.kJSPropertyAttributeNone, nil)
	C.JSStringRelease(process_js)
	createCustomFunction(context, globalObject, "Add", C.JSObjectCallAsFunctionCallback(math.Add()))
	createCustomFunction(context, globalObject, "Mult", C.JSObjectCallAsFunctionCallback(math.Mult()))
	createCustomFunction(context, globalObject, "Div", C.JSObjectCallAsFunctionCallback(math.Div()))
	createCustomFunction(context, globalObject, "require", C.JSObjectCallAsFunctionCallback(require.Require()))
	createCustomFunction(context, globalObject, "print", C.JSObjectCallAsFunctionCallback(console.Log()))
	createCustomFunction(context, globalObject, "prompt", C.JSObjectCallAsFunctionCallback(console.Prompt()))
	createCustomFunction(context, globalObject, "serve", C.JSObjectCallAsFunctionCallback(http.Serve()))
	createCustomFunction(context, globalObject, "get", C.JSObjectCallAsFunctionCallback(http.Get()))
	createCustomFunction(context, consoleGlobalObject, "log", C.JSObjectCallAsFunctionCallback(console.Log()))
	createCustomFunction(context, consoleGlobalObject, "warn", C.JSObjectCallAsFunctionCallback(console.Warn()))
	createCustomFunction(context, consoleGlobalObject, "error", C.JSObjectCallAsFunctionCallback(console.Error()))
	createCustomFunction(context, consoleGlobalObject, "assert", C.JSObjectCallAsFunctionCallback(console.Assert()))
	createCustomFunction(context, consoleGlobalObject, "count", C.JSObjectCallAsFunctionCallback(console.Count()))
	createCustomFunction(context, consoleGlobalObject, "countReset", C.JSObjectCallAsFunctionCallback(console.CountReset()))
	createCustomFunction(context, consoleGlobalObject, "timeLog", C.JSObjectCallAsFunctionCallback(console.TimeLog()))
	createCustomFunction(context, consoleGlobalObject, "time", C.JSObjectCallAsFunctionCallback(console.Time()))
	createCustomFunction(context, consoleGlobalObject, "timeEnd", C.JSObjectCallAsFunctionCallback(console.TimeEnd()))
	createCustomFunction(context, consoleGlobalObject, "clear", C.JSObjectCallAsFunctionCallback(console.Clear()))
	createCustomFunction(context, larGlobalObject, "color", C.JSObjectCallAsFunctionCallback(console.Color()))
	modules.Register("fs", "node:fs")
}

func SetDirnameAndFilename(context C.JSGlobalContextRef, globalObject C.JSObjectRef, dirname string, filename string) {
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

func main() {
	// Crear un contexto JavaScript global.
	context := C.JSGlobalContextCreate(nil)
	globalObject := C.JSContextGetGlobalObject(context)

	// Configurar las API en el objeto global.
	Apis(context, globalObject)

	// Verificar si hay argumentos de línea de comandos y si se proporciona el comando "run".
	switch os.Args[1] {
	case "run":
		if len(os.Args) < 2 {
			color.New(color.BgRed).Println("No JavaScript/TypeScript file provided to execute")
			os.Exit(1)
		}
		jsFileName := os.Args[2]
		absolutePath, _ := filepath.Abs(jsFileName)
		dirname := absolutePath[0 : len(absolutePath)-(len(filepath.Base(jsFileName))+1)]
		SetDirnameAndFilename(context, globalObject, dirname, absolutePath)

		// Leer el contenido del archivo JavaScript.
		fileContent := utils.ReadFile(jsFileName)

		// Crear una cadena JavaScript a partir del contenido del archivo.
		fileContentC := C.CString(fileContent)
		jsCode := C.JSStringCreateWithUTF8CString(fileContentC)
		C.free(unsafe.Pointer(fileContentC))
		defer C.JSStringRelease(jsCode)

		// Evaluar el script JavaScript.
		result := C.JSEvaluateScript(context, jsCode, globalObject, nil, 1, nil)

		// Convertir el resultado a una cadena de Go.
		resultStringJS := C.JSValueToStringCopy(context, result, nil)
		defer C.JSStringRelease(resultStringJS)

		// Obtener el tamaño máximo necesario para la cadena UTF-8.
		bufferSize := C.JSStringGetMaximumUTF8CStringSize(resultStringJS)
		resultCString := make([]C.char, bufferSize)
		C.JSStringGetUTF8CString(resultStringJS, &resultCString[0], bufferSize)

		// Imprimir el resultado.
		fmt.Printf("%s\n", C.GoString(&resultCString[0]))
	case "repl":
		// Iniciar el modo REPL.
		for {
			fmt.Print("> ")
			reader := bufio.NewReader(os.Stdin)
			content, err := reader.ReadString('\n')
			if err != nil {
				color.New(color.BgRed).Println("Error reading input")
				fmt.Println(err)
				os.Exit(1)
			}
			if content == ".exit" {
				fmt.Print("\n")
				os.Exit(0)
			}
			result_ts := api.Transform(string("_ ="+content), api.TransformOptions{
				Loader: api.LoaderTS,
				TsconfigRaw: `{
					"experimentalDecorators": true,
					"emitDecoratorMetadata": true,
					"allowJs": true,
				}`,
				Format: api.FormatCommonJS,
			})
			if len(result_ts.Errors) != 0 {
				os.Exit(1)
			}
			content = string(result_ts.Code)

			contentC := C.CString(content)
			jsCode := C.JSStringCreateWithUTF8CString(contentC)
			C.free(unsafe.Pointer(contentC))

			result := C.JSEvaluateScript(context, jsCode, globalObject, nil, 1, nil)
			C.JSStringRelease(jsCode)
			resultStringJS := C.JSValueToStringCopy(context, result, nil)

			bufferSize := C.JSStringGetMaximumUTF8CStringSize(resultStringJS)
			resultCString := make([]C.char, bufferSize)
			C.JSStringGetUTF8CString(resultStringJS, &resultCString[0], bufferSize)
			C.JSStringRelease(resultStringJS)
			fmt.Printf("%s\n", C.GoString(&resultCString[0]))

			c := make(chan os.Signal, 1)
			signal.Notify(c, os.Interrupt)
			go func() {
				for range c {
					fmt.Print("\n")
					os.Exit(0)
				}
			}()
		}
	default:
		color.New(color.BgRed).Println("No command provided")
		os.Exit(1)
	}

	// Liberar el contexto JavaScript global.
	C.JSGlobalContextRelease(context)
}
