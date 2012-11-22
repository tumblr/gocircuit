package acid

import (
	"circuit/use/circuit"
	"circuit/exp/file"
	"circuit/load/config"
	"os"
	"path"
)

func (a *Acid) JailOpen(jailFile string) (circuit.X, error) {
	f, err := os.Open(path.Join(config.Config.Install.JailDir(), jailFile))
	if err != nil {
		return nil, err
	}
	return circuit.Ref(file.NewFileServer(f)), nil
}
