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
	"largo/src/fs"
	"largo/src/math"
	"largo/src/require"
	"largo/src/utils"
	"os"
	"os/signal"
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
	lar_str := C.CString("Lar")
	lar_js := C.JSStringCreateWithUTF8CString(lar_str)
	C.free(unsafe.Pointer(lar_str))
	larGlobalObject := C.JSObjectMake(context, nil, nil)
	C.JSObjectSetProperty(context, globalObject, lar_js, larGlobalObject, C.kJSPropertyAttributeNone, nil)
	fs_str := C.CString("fs")
	fs_js := C.JSStringCreateWithUTF8CString(fs_str)
	C.free(unsafe.Pointer(fs_str))
	fsGlobalObject := C.JSObjectMake(context, nil, nil)
	C.JSObjectSetProperty(context, globalObject, fs_js, fsGlobalObject, C.kJSPropertyAttributeNone, nil)
	createCustomFunction(context, globalObject, "Add", C.JSObjectCallAsFunctionCallback(math.Add()))
	createCustomFunction(context, globalObject, "Mult", C.JSObjectCallAsFunctionCallback(math.Mult()))
	createCustomFunction(context, globalObject, "Div", C.JSObjectCallAsFunctionCallback(math.Div()))
	createCustomFunction(context, globalObject, "require", C.JSObjectCallAsFunctionCallback(require.Require()))
	createCustomFunction(context, globalObject, "print", C.JSObjectCallAsFunctionCallback(console.Log()))
	createCustomFunction(context, globalObject, "prompt", C.JSObjectCallAsFunctionCallback(console.Prompt()))
	createCustomFunction(context, consoleGlobalObject, "log", C.JSObjectCallAsFunctionCallback(console.Log()))
	createCustomFunction(context, consoleGlobalObject, "error", C.JSObjectCallAsFunctionCallback(console.Error()))
	createCustomFunction(context, consoleGlobalObject, "time", C.JSObjectCallAsFunctionCallback(console.Time()))
	createCustomFunction(context, consoleGlobalObject, "timeEnd", C.JSObjectCallAsFunctionCallback(console.TimeEnd()))
	createCustomFunction(context, consoleGlobalObject, "clear", C.JSObjectCallAsFunctionCallback(console.Clear()))
	createCustomFunction(context, larGlobalObject, "color", C.JSObjectCallAsFunctionCallback(console.Color()))
	createCustomFunction(context, fsGlobalObject, "readFileSync", C.JSObjectCallAsFunctionCallback(fs.ReadFileSync()))
	createCustomFunction(context, fsGlobalObject, "writeFileSync", C.JSObjectCallAsFunctionCallback(fs.WriteFileSync()))
	C.JSStringRelease(console_js)
	C.JSStringRelease(lar_js)
	C.JSStringRelease(fs_js)
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
