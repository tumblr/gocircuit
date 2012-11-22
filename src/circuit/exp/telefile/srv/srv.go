package srv

import (
	"circuit/exp/file"
	"circuit/use/circuit"
	"os"
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
	return circuit.Ref(file.NewFileServer(f))
}
