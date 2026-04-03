package modules

/*
#cgo darwin CFLAGS: -x objective-c
#cgo darwin LDFLAGS: -framework JavaScriptCore
#include <JavaScriptCore/JavaScript.h>
#include <stdlib.h>
*/
import "C"

var Commands = []string{}

func Register(name, alias string) {
	Commands = append(Commands, name, alias)
}
