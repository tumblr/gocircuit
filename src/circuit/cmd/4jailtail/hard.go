package main

import (
	"circuit/load/config"
	"circuit/use/circuit"
	"io"
	"os"
	"os/exec"
	"path"
)

func tailViaSSH(addr circuit.Addr, jailpath string) {

	abs := path.Join(config.Config.Install.JailDir(), jailpath)

	cmd := exec.Command("ssh", addr.Host().String(), "tail -f " + abs)

	stdout, err := cmd.StdoutPipe()
	if err != nil {
		println("Pipe problem:", err.Error())
		os.Exit(1)
	}

	if err = cmd.Start(); err != nil {
		println("Exec problem:", err.Error())
		os.Exit(1)
	}

	io.Copy(os.Stdout, stdout)
}
