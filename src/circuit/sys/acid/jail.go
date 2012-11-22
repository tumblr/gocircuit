package acid

import (
	"circuit/use/circuit"
	"circuit/exp/file"
	"circuit/load/config"
	"os"
	"path"
)

func (a *Acid) JailOpen(jailFile string) (circuit.X, error) {
	println("JailOpen:", jailFile)
	f, err := os.Open(path.Join(config.Config.Install.JailDir(), jailFile))
	if err != nil {
		println("joerr", err.Error())
		return nil, err
	}
	println("aha")
	return circuit.Ref(file.NewFileServer(f)), nil
}
