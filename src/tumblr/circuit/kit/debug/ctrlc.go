package debug

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"
)

// InstallTimeout panics the current process in ns time
func InstallTimeoutPanic(ns int64) {
	go func() {
		k := int(ns / 1e9)
		for i := 0; i < k; i++ {
			time.Sleep(time.Second)
			fmt.Fprintf(os.Stderr, "•%d/%d•\n", i, k)
		}
		//time.Sleep(time.Duration(ns))
		panic("process timeout")
	}()
}

// InstallCtrlCPanic installs a Ctrl-C signal handler that panics
func InstallCtrlCPanic() {
	go func() {
		//defer SavePanicTrace()
		ch := make(chan os.Signal)
		signal.Notify(ch, os.Interrupt, os.Kill)
		for s := range ch {
			log.Printf("ctrl-c/kill interruption: %s\n", s)
			panic("ctrl-c/kill")
		}
	}()
}

func SavePanicTrace() {
	r := recover()
	if r == nil {
		return
	}
	// Redirect stderr
	file, err := os.Create("panic")
	if err != nil {
		panic("dumper (no file) " + r.(fmt.Stringer).String())
	}
	syscall.Dup2(int(file.Fd()), int(os.Stderr.Fd()))
	// TRY: defer func() { file.Close() }()
	panic("dumper " + r.(string))
}
