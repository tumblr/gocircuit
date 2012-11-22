package srv

import (
	"circuit/exp/file"
	"circuit/use/circuit"
	"os"
	"time"
)

type App struct{}

func init() {
	circuit.RegisterFunc(App{})
}

func (App) Open(filepath string) circuit.X {
	f, err := os.Open(filepath)
	if err != nil {
		return nil
	}
	circuit.Daemonize(func() { time.Sleep(5*time.Second) })
	return circuit.Ref(file.NewFileServer(f))
}
