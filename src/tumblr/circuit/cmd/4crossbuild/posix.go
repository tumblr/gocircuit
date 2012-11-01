package main

import (
	"os/exec"
	"io/ioutil"
)

func Run(prog, dir string, stdin string, argv ...string) (stdout, stderr string, err error) {
	cmd := exec.Command(prog, argv...)
	cmd.Dir = dir

	stdinWriter, err := cmd.StdinPipe()
	if err != nil {
		return "", "", err
	}
	stdoutReader, err := cmd.StdoutPipe()
	if err != nil {
		return "", "", err
	}
	stderrReader, err := cmd.StderrPipe()
	if err != nil {
		return "", "", err
	}
	if err := cmd.Start(); err != nil {
		return "", "", err
	}
	// Since Run is meant for non-interactive execution, we pump all the stdin first,
	// then we (sequentially) read all of stdout and then all of stderr.
	// XXX: Is it possible to block if the program's stderr buffer fills while we are
	// consuming the stdout?
	_, err = stdinWriter.Write([]byte(stdin))
	if err != nil {
		return "", "", err
	}
	err = stdinWriter.Close()
	if err != nil {
		return "", "", err
	}
	stdoutBuf, _ := ioutil.ReadAll(stdoutReader)
	stderrBuf, _ := ioutil.ReadAll(stderrReader)

	return string(stdoutBuf), string(stderrBuf), cmd.Wait()
}

func RunShell(shellStdin string) (stdout, stderr string, err error) {
	return Run("sh", "", shellStdin)
}

func DownloadDir(remoteHost, remoteDir, sourceDir string) error {
	_, _, err := Run("rsync", "", "", "-acrv", "--rsh=ssh", remoteHost + ":" + remoteDir + "/", sourceDir + "/")
	return err
}
