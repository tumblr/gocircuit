package config

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path"
)

// InstallConfig holds configuration parameters regarding circuit installation on host machines
type InstallConfig struct {
	RootDir string  // Root directory of circuit installation on
	LibPath string  // Any additions to the library path for execution time
	Binary  string  // Desired name for the circuit runtime binary
}

func (i *InstallConfig) BinDir() string {
	return path.Join(i.RootDir, "bin")
}

func (i *InstallConfig) JailDir() string {
	return path.Join(i.RootDir, "jail")
}

func (i *InstallConfig) VarDir() string {
	return path.Join(i.RootDir, "var")
}

func (i *InstallConfig) BinaryPath() string {
	return path.Join(i.BinDir(), i.Binary)
}

func (i *InstallConfig) ClearHelperPath() string {
	return path.Join(i.BinDir(), "4clear-helper")
}

func parseInstall() {
	Config.Install = &InstallConfig{}

	// Try parsing install config from environment
	Config.Install.RootDir = os.Getenv("_CIR_IR")
	Config.Install.LibPath = os.Getenv("_CIR_IL")
	Config.Install.Binary = os.Getenv("_CIR_IB")
	if Config.Install.RootDir != "" {
		return
	}

	// Try parsing the install config from a file
	ifile := os.Getenv("CIR_INSTALL")
	if ifile == "" {
		Config.Install = nil
		return
	}
	data, err := ioutil.ReadFile(ifile)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Problem reading install config file (%s)", err)
		os.Exit(1)
	}
	if err := json.Unmarshal(data, Config.Install); err != nil {
		fmt.Fprintf(os.Stderr, "Problem parsing install config file (%s)", err)
		os.Exit(1)
	}
}
