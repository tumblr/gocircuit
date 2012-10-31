package boot

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
)

// BuildConfig holds configuration parameters for the automated circuit app build system
type BuildConfig struct {
	Repo             string  // User repo to fetch
	GoPath           string  // User GOPATH relative to user repo; or empty string if user repo meant to be cloned at root of GOPATH
	Pkg              string  // User side-effect package to include in the circuit runtime build
	Host             string  // Host where build takes place
	Tool             string  // Build tool path on build host
	Jail             string  // Build jail path on build host
	RebuildGo        bool    // Rebuild Go even if a newer version is not available
	ShipDir          string  // Local directory where built runtime binary and dynamic libraries will be delivered
	ZookeeperInclude string  // Path to Zookeeper include files on build host
	ZookeeperLib     string  // Path to Zookeeper library files on build host
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
