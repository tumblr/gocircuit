// Package install is responsible for reading the install system's configuration
// from a file named by the CIR_INSTALL environment variable
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

var Install *InstallConfig

func parseInstall() {
	Install = &InstallConfig{}

	// Try parsing install config from environment
	Install.RootDir = os.Getenv("_CIR_IR")
	Install.LibPath = os.Getenv("_CIR_IL")
	Install.Binary = os.Getenv("_CIR_IB")
	if Install.RootDir != "" {
		return
	}

	// Try parsing the install config from a file
	ifile := os.Getenv("CIR_INSTALL")
	if ifile == "" {
		Install = nil
		return
	}
	data, err := ioutil.ReadFile(ifile)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Problem reading install config file (%s)", err)
		os.Exit(1)
	}
	if err := json.Unmarshal(data, Install); err != nil {
		fmt.Fprintf(os.Stderr, "Problem parsing install config file (%s)", err)
		os.Exit(1)
	}
}
