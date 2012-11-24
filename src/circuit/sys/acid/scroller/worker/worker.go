package worker

import (
	"time"
	"circuit/use/circuit"
)

type App struct{}

func (App) Main() {
	circuit.Daemonize(func() {
		for i := 0; ; i++ {
			println(i)
			time.Sleep(time.Second)
		}
	})
}

func init() { 
	circuit.RegisterFunc(App{}) 
}
