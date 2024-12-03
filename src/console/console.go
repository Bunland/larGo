package console

/*
#cgo CFLAGS: -I/usr/include/webkitgtk-4.0
#cgo LDFLAGS: -ljavascriptcoregtk-4.0
#include <JavaScriptCore/JavaScript.h>
#include <stdlib.h>

// Declarar las funciones LogF, TimeF, TimeEndF para que C reconozca su existencia.
extern JSValueRef LogF(JSContextRef context, JSObjectRef function, JSObjectRef thisObject, size_t argumentCount, JSValueRef arguments[], JSValueRef* exception);
extern JSValueRef WarnF(JSContextRef context, JSObjectRef function, JSObjectRef thisObject, size_t argumentCount, JSValueRef arguments[], JSValueRef* exception);
extern JSValueRef ErrorF(JSContextRef context, JSObjectRef function, JSObjectRef thisObject, size_t argumentCount, JSValueRef arguments[], JSValueRef* exception);
extern JSValueRef AssertF(JSContextRef context, JSObjectRef function, JSObjectRef thisObject, size_t argumentCount, JSValueRef arguments[], JSValueRef* exception);
extern JSValueRef TimeF(JSContextRef context, JSObjectRef function, JSObjectRef thisObject, size_t argumentCount, JSValueRef arguments[], JSValueRef* exception);
extern JSValueRef TimeEndF(JSContextRef context, JSObjectRef function, JSObjectRef thisObject, size_t argumentCount, JSValueRef arguments[], JSValueRef* exception);
extern JSValueRef ClearF(JSContextRef context, JSObjectRef function, JSObjectRef thisObject, size_t argumentCount, JSValueRef arguments[], JSValueRef* exception);
extern JSValueRef ColorF(JSContextRef context, JSObjectRef function, JSObjectRef thisObject, size_t argumentCount, JSValueRef arguments[], JSValueRef* exception);
extern JSValueRef PromptF(JSContextRef context, JSObjectRef function, JSObjectRef thisObject, size_t argumentCount, JSValueRef arguments[], JSValueRef* exception);
extern JSValueRef CountF(JSContextRef context, JSObjectRef function, JSObjectRef thisObject, size_t argumentCount, JSValueRef arguments[], JSValueRef* exception);
extern JSValueRef CountResetF(JSContextRef context, JSObjectRef function, JSObjectRef thisObject, size_t argumentCount, JSValueRef arguments[], JSValueRef* exception);
extern JSValueRef TimeLogF(JSContextRef context, JSObjectRef function, JSObjectRef thisObject, size_t argumentCount, JSValueRef arguments[], JSValueRef* exception);
*/
import "C"
import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"time"
	"unsafe"

	"github.com/fatih/color"
)

type WatcherStruct struct {
	start time.Time
	label string
}

type CounterStruct struct {
	count int
}

var watcher = make(map[string]*WatcherStruct)
var counter = make(map[string]*CounterStruct)

func IsConstructor(context C.JSContextRef, value C.JSValueRef, stringToCheck string) (isInstance bool) {
	constructorString := C.CString(stringToCheck)
	constructorJSString := C.JSStringCreateWithUTF8CString(constructorString)
	C.free(unsafe.Pointer(constructorString))
	constructorValue := C.JSObjectGetProperty(context, C.JSContextGetGlobalObject(context), constructorJSString, nil)
	C.JSStringRelease(constructorJSString)
	isInstance = bool(C.JSValueIsInstanceOfConstructor(context, value, constructorValue, nil))
	return
}

// Log es la implementación de console.log de JavaScript
//
//export LogF
func LogF(context C.JSContextRef, function C.JSObjectRef, thisObject C.JSObjectRef, argumentCount C.size_t, arguments *C.JSValueRef, exception *C.JSValueRef) C.JSValueRef {
	if int(argumentCount) <= 0 {
		return C.JSValueMakeUndefined(context)
	}
	argumentSlice := (*[1 << 30]C.JSValueRef)(unsafe.Pointer(arguments))[:argumentCount:argumentCount]
	var text string = ""
	for i := 0; i < int(argumentCount); i += 1 {
		if C.JSValueIsString(context, argumentSlice[i]) {
			str := C.JSValueToStringCopy(context, argumentSlice[i], nil)

			bufferSize := C.JSStringGetMaximumUTF8CStringSize(str)

			buffer := C.malloc(bufferSize)
			C.JSStringGetUTF8CString(str, (*C.char)(buffer), bufferSize)

			text += fmt.Sprintf("%s ", C.GoString((*C.char)(buffer)))

			C.free(unsafe.Pointer(buffer))

			C.JSStringRelease(str)
		} else if C.JSValueIsObject(context, argumentSlice[i]) && !C.JSValueIsNull(context, argumentSlice[i]) && !C.JSValueIsUndefined(context, argumentSlice[i]) {
			if IsConstructor(context, argumentSlice[i], "RegExp") || IsConstructor(context, argumentSlice[i], "Error") {
				str := C.JSValueToStringCopy(context, argumentSlice[i], nil)

				bufferSize := C.JSStringGetMaximumUTF8CStringSize(str)

				buffer := C.malloc(bufferSize)
				C.JSStringGetUTF8CString(str, (*C.char)(buffer), bufferSize)

				text += fmt.Sprintf("%s ", C.GoString((*C.char)(buffer)))

				C.free(unsafe.Pointer(buffer))

				C.JSStringRelease(str)
			}
			json := C.JSValueCreateJSONString(context, argumentSlice[i], 0, exception)

			bufferSize := C.JSStringGetMaximumUTF8CStringSize(json)
			buffer := C.malloc(bufferSize)
			C.JSStringGetUTF8CString(json, (*C.char)(buffer), bufferSize)

			text += fmt.Sprintf("%s, ", C.GoString((*C.char)(buffer)))

			C.free(unsafe.Pointer(buffer))

			C.JSStringRelease(json)
		} else if C.JSValueIsNumber(context, argumentSlice[i]) {
			number := C.JSValueToNumber(context, argumentSlice[i], exception)
			text += fmt.Sprintf("%s ", strconv.FormatFloat(float64(number), 'f', -1, 64))
		} else if C.JSValueIsBoolean(context, argumentSlice[i]) {
			boolean := C.JSValueToBoolean(context, argumentSlice[i])
			if boolean {
				text += "true "
			} else {
				text += "false "
			}
		} else if C.JSValueIsNull(context, argumentSlice[i]) {
			text += "null "
		} else if C.JSValueIsUndefined(context, argumentSlice[i]) {
			text += "undefined "
		}
	}
	fmt.Println(strings.TrimRight(text, " "))

	return C.JSValueMakeUndefined(context)
}

// Warn es la implementación de console.warn de JavaScript
//
//export WarnF
func WarnF(context C.JSContextRef, function C.JSObjectRef, thisObject C.JSObjectRef, argumentCount C.size_t, arguments *C.JSValueRef, exception *C.JSValueRef) C.JSValueRef {
	if int(argumentCount) <= 0 {
		return C.JSValueMakeUndefined(context)
	}
	argumentSlice := (*[1 << 30]C.JSValueRef)(unsafe.Pointer(arguments))[:argumentCount:argumentCount]
	for i := 0; i < int(argumentCount); i += 1 {
		str := C.JSValueToStringCopy(context, argumentSlice[i], nil)

		bufferSize := C.JSStringGetMaximumUTF8CStringSize(str)

		buffer := C.malloc(bufferSize)
		C.JSStringGetUTF8CString(str, (*C.char)(buffer), bufferSize)

		color.New(color.FgYellow).Printf("%s ", C.GoString((*C.char)(buffer)))

		C.free(unsafe.Pointer(buffer))

		C.JSStringRelease(str)
	}
	fmt.Print("\n")

	return C.JSValueMakeUndefined(context)
}

// ErrorF es la implementación de console.error de JavaScript
//
//export ErrorF
func ErrorF(context C.JSContextRef, function C.JSObjectRef, thisObject C.JSObjectRef, argumentCount C.size_t, arguments *C.JSValueRef, exception *C.JSValueRef) C.JSValueRef {
	if int(argumentCount) <= 0 {
		return C.JSValueMakeUndefined(context)
	}
	argumentSlice := (*[1 << 30]C.JSValueRef)(unsafe.Pointer(arguments))[:argumentCount:argumentCount]
	for i := 0; i < int(argumentCount); i += 1 {
		str := C.JSValueToStringCopy(context, argumentSlice[i], nil)

		bufferSize := C.JSStringGetMaximumUTF8CStringSize(str)

		buffer := C.malloc(bufferSize)
		C.JSStringGetUTF8CString(str, (*C.char)(buffer), bufferSize)

		color.New(color.FgRed).Printf("%s ", C.GoString((*C.char)(buffer)))

		C.free(unsafe.Pointer(buffer))

		C.JSStringRelease(str)
	}
	fmt.Print("\n")

	return C.JSValueMakeUndefined(context)
}

// Assert es la implementación de console.assert de JavaScript
//
//export AssertF
func AssertF(context C.JSContextRef, function C.JSObjectRef, thisObject C.JSObjectRef, argumentCount C.size_t, arguments *C.JSValueRef, exception *C.JSValueRef) C.JSValueRef {
	if int(argumentCount) <= 0 {
		return C.JSValueMakeNull(context)
	}
	argumentSlice := (*[1 << 30]C.JSValueRef)(unsafe.Pointer(arguments))[:argumentCount:argumentCount]
	assertionNoFailed := C.JSValueToBoolean(context, argumentSlice[0])
	if !assertionNoFailed {
		var text string = " "
		for index := 1; index < len(argumentSlice); index += 1 {
			if C.JSValueIsObject(context, argumentSlice[index]) {
				json := C.JSValueCreateJSONString(context, argumentSlice[index], 0, exception)
				maximumSize := C.JSStringGetMaximumUTF8CStringSize(json)
				buffer := C.malloc(maximumSize)
				C.JSStringGetUTF8CString(json, (*C.char)(buffer), maximumSize)
				jsonString := C.GoString((*C.char)(buffer))
				C.free(unsafe.Pointer(buffer))
				text += jsonString + ", "
			} else if C.JSValueIsString(context, argumentSlice[index]) {
				str := C.JSValueToStringCopy(context, argumentSlice[index], nil)
				bufferSize := C.JSStringGetMaximumUTF8CStringSize(str)
				buffer := C.malloc(bufferSize)
				C.JSStringGetUTF8CString(str, (*C.char)(buffer), bufferSize)
				text += C.GoString((*C.char)(buffer)) + ", "
				C.free(unsafe.Pointer(buffer))
			} else if C.JSValueIsNumber(context, argumentSlice[index]) {
				number := C.JSValueToNumber(context, argumentSlice[index], exception)
				text += fmt.Sprintf("%f, ", number)
			} else if C.JSValueIsBoolean(context, argumentSlice[index]) {
				boolean := C.JSValueToBoolean(context, argumentSlice[index])
				if boolean {
					text += "true, "
				} else {
					text += "false, "
				}
			} else if C.JSValueIsNull(context, argumentSlice[index]) {
				text += "null, "
			} else {
				text += "undefined, "
			}
		}
		textWithoutCommaAndSpace, _ := strings.CutSuffix(text, ", ")
		if strings.TrimSpace(textWithoutCommaAndSpace) == "" {
			fmt.Println("Assertion failed")
		} else {
			fmt.Println("Assertion failed:" + textWithoutCommaAndSpace)
		}
	}
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

// Clear es la implementación de console.assert de JavaScript
//
//export ClearF
func ClearF(context C.JSContextRef, function C.JSObjectRef, thisObject C.JSObjectRef, argumentCount C.size_t, arguments *C.JSValueRef, exception *C.JSValueRef) C.JSValueRef {
	cmd := exec.Command("clear")
	cmd.Stdout = os.Stdout
	cmd.Run()
	return C.JSValueMakeUndefined(context)
}

// Color es la implementación de Lar.color de JavaScript
//
//export ColorF
func ColorF(context C.JSContextRef, function C.JSObjectRef, thisObject C.JSObjectRef, argumentCount C.size_t, arguments *C.JSValueRef, exception *C.JSValueRef) C.JSValueRef {
	if int(argumentCount) < 2 {
		return C.JSValueMakeUndefined(context)
	}
	argumentSlice := (*[1 << 30]C.JSValueRef)(unsafe.Pointer(arguments))[:argumentCount:argumentCount]
	// Convert the first argument to a string.
	str := C.JSValueToStringCopy(context, argumentSlice[0], nil)
	bufferSize := C.JSStringGetMaximumUTF8CStringSize(str)
	buffer := C.malloc(bufferSize)
	C.JSStringGetUTF8CString(str, (*C.char)(buffer), bufferSize)
	color_type := C.GoString((*C.char)(buffer))
	C.free(unsafe.Pointer(buffer))
	for i := 1; i < int(argumentCount); i += 1 {
		str_two := C.JSValueToStringCopy(context, argumentSlice[i], nil)
		bufferSize_two := C.JSStringGetMaximumUTF8CStringSize(str_two)
		buffer_two := C.malloc(bufferSize_two)
		C.JSStringGetUTF8CString(str_two, (*C.char)(buffer_two), bufferSize_two)
		color_value := C.GoString((*C.char)(buffer_two))
		switch color_type {
		case "red":
			color.New(color.FgRed).Printf("%s ", color_value)
		case "green":
			color.New(color.FgGreen).Printf("%s ", color_value)
		case "blue":
			color.New(color.FgBlue).Printf("%s ", color_value)
		case "black":
			color.New(color.FgBlack).Printf("%s ", color_value)
		case "cyan":
			color.New(color.FgCyan).Printf("%s ", color_value)
		case "yellow":
			color.New(color.FgYellow).Printf("%s ", color_value)
		case "magenta":
			color.New(color.FgMagenta).Printf("%s ", color_value)
		case "white":
			color.New(color.FgWhite).Printf("%s ", color_value)
		default:
			return C.JSValueMakeUndefined(context)
		}
		C.free(unsafe.Pointer(buffer_two))
	}
	fmt.Print("\n")
	return C.JSValueMakeUndefined(context)
}

// PromptF hace la función prompt() de JavaScript
//
//export PromptF
func PromptF(context C.JSContextRef, function C.JSObjectRef, thisObject C.JSObjectRef, argumentCount C.size_t, arguments *C.JSValueRef, exception *C.JSValueRef) C.JSValueRef {
	if int(argumentCount) < 1 {
		return C.JSValueMakeUndefined(context)
	}
	argumentSlice := (*[1 << 30]C.JSValueRef)(unsafe.Pointer(arguments))[:argumentCount:argumentCount]
	str := C.JSValueToStringCopy(context, argumentSlice[0], nil)
	bufferSize := C.JSStringGetMaximumUTF8CStringSize(str)
	buffer := C.malloc(bufferSize)
	C.JSStringGetUTF8CString(str, (*C.char)(buffer), bufferSize)
	question := C.GoString((*C.char)(buffer)) + " "
	fmt.Print(question)
	reader := bufio.NewReader(os.Stdin)
	answer, err := reader.ReadString('\n')
	if err != nil {
		return C.JSValueMakeUndefined(context)
	}
	C.free(unsafe.Pointer(buffer))
	c_string := C.CString(strings.TrimRight(answer, "\n"))
	file_c_string := C.JSStringCreateWithUTF8CString(c_string)
	C.free(unsafe.Pointer(c_string))
	return C.JSValueMakeString(context, file_c_string)
}

// CountF es la implementación de console.count() de JavaScript
//
//export CountF
func CountF(context C.JSContextRef, function C.JSObjectRef, thisObject C.JSObjectRef, argumentCount C.size_t, arguments *C.JSValueRef, exception *C.JSValueRef) C.JSValueRef {
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
	if _, ok := counter[label]; !ok {
		counter[label] = &CounterStruct{
			count: 0,
		}
	}
	counter[label].count += 1
	fmt.Printf("%s: %d\n", label, counter[label].count)
	return C.JSValueMakeUndefined(context)
}

// CountResetF es la implementación de console.countReset() de JavaScript
//
//export CountResetF
func CountResetF(context C.JSContextRef, function C.JSObjectRef, thisObject C.JSObjectRef, argumentCount C.size_t, arguments *C.JSValueRef, exception *C.JSValueRef) C.JSValueRef {
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
	if _, ok := counter[label]; !ok {
		return C.JSValueMakeUndefined(context)
	}
	counter[label].count = 0
	return C.JSValueMakeUndefined(context)
}

// TimeLogF es la implementación de console.timeLog() de JavaScript
//
//export TimeLogF
func TimeLogF(context C.JSContextRef, function C.JSObjectRef, thisObject C.JSObjectRef, argumentCount C.size_t, arguments *C.JSValueRef, exception *C.JSValueRef) (result C.JSValueRef) {
	argumentSlice := (*[1 << 30]C.JSValueRef)(unsafe.Pointer(arguments))[:argumentCount:argumentCount]
	var label string
	if argumentCount <= 0 {
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
	for i := 1; i < int(argumentCount); i += 1 {
		if C.JSValueIsString(context, argumentSlice[i]) {
			str := C.JSValueToStringCopy(context, argumentSlice[i], nil)
			bufferSize := C.JSStringGetMaximumUTF8CStringSize(str)
			buffer := C.malloc(bufferSize)
			C.JSStringGetUTF8CString(str, (*C.char)(buffer), bufferSize)
			fmt.Printf("%s ", C.GoString((*C.char)(buffer)))
			C.free(unsafe.Pointer(buffer))
		} else if C.JSValueIsObject(context, argumentSlice[i]) && !C.JSValueIsNull(context, argumentSlice[i]) && !C.JSValueIsUndefined(context, argumentSlice[i]) {
			if IsConstructor(context, argumentSlice[i], "RegExp") || IsConstructor(context, argumentSlice[i], "Error") {
				str := C.JSValueToStringCopy(context, argumentSlice[i], nil)
				bufferSize := C.JSStringGetMaximumUTF8CStringSize(str)
				buffer := C.malloc(bufferSize)
				C.JSStringGetUTF8CString(str, (*C.char)(buffer), bufferSize)
				fmt.Printf("%s ", C.GoString((*C.char)(buffer)))
				C.free(unsafe.Pointer(buffer))
			}
			json := C.JSValueCreateJSONString(context, argumentSlice[i], 0, exception)
			bufferSize := C.JSStringGetMaximumUTF8CStringSize(json)
			buffer := C.malloc(bufferSize)
			C.JSStringGetUTF8CString(json, (*C.char)(buffer), bufferSize)
			fmt.Printf("%s ", C.GoString((*C.char)(buffer)))
			C.free(unsafe.Pointer(buffer))
			C.JSStringRelease(json)
		} else if C.JSValueIsNumber(context, argumentSlice[i]) {
			number := C.JSValueToNumber(context, argumentSlice[i], exception)
			fmt.Printf("%f ", number)
		} else if C.JSValueIsBoolean(context, argumentSlice[i]) {
			boolean := C.JSValueToBoolean(context, argumentSlice[i])
			if boolean {
				fmt.Printf("true ")
			} else {
				fmt.Printf("false ")
			}
		} else if C.JSValueIsNull(context, argumentSlice[i]) {
			fmt.Printf("null ")
		} else if C.JSValueIsUndefined(context, argumentSlice[i]) {
			fmt.Printf("undefined ")
		}
	}
	return C.JSValueMakeUndefined(context)
}

// Log devuelve la función de callback de C para la función console.log() en JavaScript.
func Log() C.JSObjectCallAsFunctionCallback {
	return C.JSObjectCallAsFunctionCallback(C.LogF)
}

// Time devuelve la función de callback de C para la función console.time() en JavaScript.
func Time() C.JSObjectCallAsFunctionCallback {
	return C.JSObjectCallAsFunctionCallback(C.TimeF)
}

// TimeEnd devuelve la función de callback de C para la función console.timeEnd() en JavaScript.
func TimeEnd() C.JSObjectCallAsFunctionCallback {
	return C.JSObjectCallAsFunctionCallback(C.TimeEndF)
}

// Clear devuelve la función de callback de C para la función console.clear() en JavaScript.
func Clear() C.JSObjectCallAsFunctionCallback {
	return C.JSObjectCallAsFunctionCallback(C.ClearF)
}

// Error devuelve la función de callback de C para la función console.error() en JavaScript.
func Error() C.JSObjectCallAsFunctionCallback {
	return C.JSObjectCallAsFunctionCallback(C.ErrorF)
}

// Color devuelve la función de callback de C para la función Lar.color() en JavaScript.
func Color() C.JSObjectCallAsFunctionCallback {
	return C.JSObjectCallAsFunctionCallback(C.ColorF)
}

// Prompt devuelve la función callback de C para la función prompt() de JavaScript.
func Prompt() C.JSObjectCallAsFunctionCallback {
	return C.JSObjectCallAsFunctionCallback(C.PromptF)
}

// Assert devuelve la función callback de C para la función console.assert() de JavaScript.
func Assert() C.JSObjectCallAsFunctionCallback {
	return C.JSObjectCallAsFunctionCallback(C.AssertF)
}

// Warn devuelve la función callback de C para la función console.warn() de JavaScript.
func Warn() C.JSObjectCallAsFunctionCallback {
	return C.JSObjectCallAsFunctionCallback(C.WarnF)
}

// Count devuelve la función callback de C para la función console.count() de JavaScript.
func Count() C.JSObjectCallAsFunctionCallback {
	return C.JSObjectCallAsFunctionCallback(C.CountF)
}

// CountReset devuelve la función callback de C para la función console.countReset() de JavaScript.
func CountReset() C.JSObjectCallAsFunctionCallback {
	return C.JSObjectCallAsFunctionCallback(C.CountResetF)
}

// TimeLog devuelve la función callback de C para la función console.timeLog() de JavaScript.
func TimeLog() C.JSObjectCallAsFunctionCallback {
	return C.JSObjectCallAsFunctionCallback(C.TimeLogF)
}
