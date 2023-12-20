package console

/*
#cgo CFLAGS: -I/usr/include/webkitgtk-4.0
#cgo LDFLAGS: -ljavascriptcoregtk-4.0
#include <stdlib.h>
#include <JavaScriptCore/JavaScript.h>

// Declarar las funciones LogF, TimeF, TimeEndF para que C reconozca su existencia.
extern JSValueRef LogF(JSContextRef context, JSObjectRef function, JSObjectRef thisObject, size_t argumentCount, JSValueRef arguments[], JSValueRef* exception);
extern JSValueRef ErrorF(JSContextRef context, JSObjectRef function, JSObjectRef thisObject, size_t argumentCount, JSValueRef arguments[], JSValueRef* exception);
extern JSValueRef TimeF(JSContextRef context, JSObjectRef function, JSObjectRef thisObject, size_t argumentCount, JSValueRef arguments[], JSValueRef* exception);
extern JSValueRef TimeEndF(JSContextRef context, JSObjectRef function, JSObjectRef thisObject, size_t argumentCount, JSValueRef arguments[], JSValueRef* exception);
extern JSValueRef ClearF(JSContextRef context, JSObjectRef function, JSObjectRef thisObject, size_t argumentCount, JSValueRef arguments[], JSValueRef* exception);
extern JSValueRef ColorF(JSContextRef context, JSObjectRef function, JSObjectRef thisObject, size_t argumentCount, JSValueRef arguments[], JSValueRef* exception);
extern JSValueRef PromptF(JSContextRef context, JSObjectRef function, JSObjectRef thisObject, size_t argumentCount, JSValueRef arguments[], JSValueRef* exception);
*/
import "C"
import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"strings"
	"time"
	"unsafe"

	"github.com/fatih/color"
)

type WatcherStruct struct {
	start time.Time
	label string
}

var watcher = make(map[string]*WatcherStruct)

// Log es la implementación de console.log de JavaScript
//
//export LogF
func LogF(context C.JSContextRef, function C.JSObjectRef, thisObject C.JSObjectRef, argumentCount C.size_t, arguments *C.JSValueRef, exception *C.JSValueRef) C.JSValueRef {
	if int(argumentCount) <= 0 {
		return C.JSValueMakeUndefined(context)
	}
	argumentSlice := (*[1 << 30]C.JSValueRef)(unsafe.Pointer(arguments))[:argumentCount:argumentCount]
	for i := 0; i < int(argumentCount); i += 1 {
		str := C.JSValueToStringCopy(context, argumentSlice[i], nil)

		bufferSize := C.JSStringGetMaximumUTF8CStringSize(str)

		buffer := C.malloc(bufferSize)
		C.JSStringGetUTF8CString(str, (*C.char)(buffer), bufferSize)

		fmt.Printf("%s ", C.GoString((*C.char)(buffer)))

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
	question := C.GoString((*C.char)(buffer))
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

// Log devuelve la función de callback de C para la función Log en JavaScript.
func Log() C.JSObjectCallAsFunctionCallback {
	return C.JSObjectCallAsFunctionCallback(C.LogF)
}

// Time devuelve la función de callback de C para la función Time en JavaScript.
func Time() C.JSObjectCallAsFunctionCallback {
	return C.JSObjectCallAsFunctionCallback(C.TimeF)
}

// TimeEnd devuelve la función de callback de C para la función TimeEnd en JavaScript.
func TimeEnd() C.JSObjectCallAsFunctionCallback {
	return C.JSObjectCallAsFunctionCallback(C.TimeEndF)
}

// Clear devuelve la función de callback de C para la función Clear en JavaScript.
func Clear() C.JSObjectCallAsFunctionCallback {
	return C.JSObjectCallAsFunctionCallback(C.ClearF)
}

// Error devuelve la función de callback de C para la función Error en JavaScript.
func Error() C.JSObjectCallAsFunctionCallback {
	return C.JSObjectCallAsFunctionCallback(C.ErrorF)
}

// Color devuelve la función de callback de C para la función Color en JavaScript.
func Color() C.JSObjectCallAsFunctionCallback {
	return C.JSObjectCallAsFunctionCallback(C.ColorF)
}

func Prompt() C.JSObjectCallAsFunctionCallback {
	return C.JSObjectCallAsFunctionCallback(C.PromptF)
}
