package config

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
)

type allConfig struct {
	Zookeeper *ZookeeperConfig
	Install   *InstallConfig
	Build     *BuildConfig
}

func parseAll(file string) {
	data, err := ioutil.ReadFile(file)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Problem reading all-in-one config file (%s)", err)
		os.Exit(1)
	}
	all := &allConfig{}
	if err := json.Unmarshal(data, all); err != nil {
		fmt.Fprintf(os.Stderr, "Problem parsing all-in-one config file (%s)", err)
		os.Exit(1)
	}
	Zookeeper = all.Zookeeper
	Install = all.Install
	Build = all.Build	
}

// Worker equals true, if this process was invoked programmatically within the circuit to
// act as a worker runtime. Otherwise, it was invoked by the user, in which case the intention
// is to run its main method after the circuitry is initialized.
var Worker bool

// WorkerEnvSwitch names the environment variable that determines whether a process is
// invoked as a worker or otherwise.
const WorkerEnvSwitch = "CIRCUIT_RUN_AS_WORKER"

func init() {
	// Check if this is a worker or a user invokation.
	Worker = os.Getenv(WorkerEnvSwitch) != ""
	if Worker {
		return
	}

	// If CIR is set, it points to a single file that contains all three configuration structures in JSON format.
	cc := os.Getenv("CIR")
	if cc != "" {
		parseAll(cc)
		return
	}

	// Otherwise, each one is parsed independently
	parseZookeeper()
	parseInstall()
	parseBuild()
}
