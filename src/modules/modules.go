package modules

/*
#cgo CFLAGS: -I/usr/include/webkitgtk-4.0
#cgo LDFLAGS: -ljavascriptcoregtk-4.0
#include <JavaScriptCore/JavaScript.h>
#include <stdlib.h>
*/
import "C"

var Commands = []string{}

func Register(name, alias string) {
	Commands = append(Commands, name, alias)
}
