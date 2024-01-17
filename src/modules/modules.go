package modules

/*
#cgo CFLAGS: -I/usr/include/webkitgtk-4.1
#cgo LDFLAGS: -ljavascriptcoregtk-4.1
#include <JavaScriptCore/JavaScript.h>
#include <stdlib.h>
*/
import "C"

var Commands = []string{}

func Register(name, alias string) {
	Commands = append(Commands, name, alias)
}
