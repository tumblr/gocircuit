package c

import (
	"fmt"
	"os"
	"sync"
)

var (
	llk    sync.Mutex
	indent int
)

func Indent() {
	llk.Lock()
	defer llk.Unlock()
	indent++
}

func Unindent() {
	llk.Lock()
	defer llk.Unlock()
	indent--
}

func Log(fmt_ string, arg_ ...interface{}) {
	for i := 0; i < indent; i++ {
		print("  ")
	}
	fmt.Fprintf(os.Stderr, fmt_, arg_...)
	println("")
}
