package acid

import (
	"circuit/use/circuit"
	teleio "circuit/kit/tele/io"
	"circuit/load/config"
	"io"
	"os/exec"
	"path"
)

// JailTail opens a file within this worker's jail directory and prepares a
// cross-circuit pointer to the open file
func (a *Acid) JailTail(jailFile string) (circuit.X, error) {
	abs := path.Join(config.Config.Install.JailDir(), circuit.WorkerAddr().RuntimeID().String(), jailFile)
	
	cmd := exec.Command("/bin/sh", "-c", "tail -f " + abs)
	/*
	stdin, err := cmd.StdinPipe()
	if err != nil {
		return nil, circuit.FlattenError(err)
	}*/

	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return nil, circuit.FlattenError(err)
	}

	if err = cmd.Start(); err != nil {
		return nil, circuit.FlattenError(err)
	}

	/*sh := []byte("tail -f " + abs)
	if _, err = stdin.Write(sh); err != nil {
		return nil, circuit.FlattenError(err)
	}
	if err = stdin.Close(); err != nil {
		return nil, circuit.FlattenError(err)
	}*/

	return circuit.Ref(teleio.NewServer(&tailStdout{stdout})), nil
}

type tailStdout struct {
	io.ReadCloser
}

func (*tailStdout) Write([]byte) (int, error) {
	panic("write not supported")
}
