package boot

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
)

// BuildConfig holds configuration parameters for the automated circuit app build system
type BuildConfig struct {
	Repo             string  // Repo to fetch
	Pkg              string  // Package to build
	Host             string  // Build host
	Tool             string  // Build tool path
	Jail             string  // Build jail path
	RebuildGo        bool    // Rebuild Go even if no new version available
	ShipDir          string  // Local directory where built runtime and libraries will be delivered
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
