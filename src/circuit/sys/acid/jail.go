package acid

import (
	"circuit/use/circuit"
	"circuit/kit/tele/file"
	"circuit/load/config"
	"os"
	"path"
)

// JailOpen opens a file within this worker's jail directory and prepares a
// cross-circuit pointer to the open file
func (a *Acid) JailOpen(jailFile string) (circuit.X, error) {
	f, err := os.Open(path.Join(config.Config.Install.JailDir(), circuit.WorkerAddr().RuntimeID().String(), jailFile))
	if err != nil {
		return nil, err
	}
	return circuit.Ref(file.NewFileServer(f)), nil
}
