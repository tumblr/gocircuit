package boot

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
)

// BuildConfig holds configuration parameters for the automated circuit app build system
type BuildConfig struct {
	Binary           string  // Desired name for circuit runtime binary
	Jail             string  // Build jail path on build host

	AppRepo          string  // App repo URL
	AppSrc           string  // App GOPATH relative to app repo; or empty string if app repo meant to be cloned inside a GOPATH

	Pkg              string  // User side-effect package to include in the circuit runtime build
	Show             bool
	RebuildGo        bool    // Rebuild Go even if a newer version is not available

	ZookeeperInclude string  // Path to Zookeeper include files on build host
	ZookeeperLib     string  // Path to Zookeeper library files on build host

	CircuitRepo      string
	CircuitSrc       string

	Host             string  // Host where build takes place
	Tool             string  // Build tool path on build host
	ShipDir          string  // Local directory where built runtime binary and dynamic libraries will be delivered
}

var Build *BuildConfig

func parseBuildConfig() {
	bfile := os.Getenv("CIR_BUILD")
	Build = &BuildConfig{}
	data, err := ioutil.ReadFile(bfile)
	if err != nil {
		Build = nil
		fmt.Fprintf(os.Stderr, "Not using a build config file (%s)", err)
		return
	}
	if err = json.Unmarshal(data, Build); err != nil {
		fmt.Fprintf(os.Stderr, "Problem parsing build file (%s)", err)
		os.Exit(1)
	}
}
