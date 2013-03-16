// Package config provides access the circuit configuration of this worker process
package config

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	// _ "circuit/kit/debug/ctrlc"
)

// Role determines the context within which this executable was invoked
const (
	Main       = "main"
	Daemonizer = "daemonizer"
	Worker     = "worker"
)
var Role string

// CIRCUIT_ROLE names the environment variable that determines the role of this invokation
const RoleEnv = "CIRCUIT_ROLE"

// init determines in what context we are being run and reads the configurations accordingly
func init() {
	Config = &WorkerConfig{}
	Role = os.Getenv(RoleEnv)
	if Role == "" {
		Role = Main
	}
	switch Role {
	case Main:
		readAsMain()
	case Daemonizer:
		readAsDaemonizerOrWorker()
	case Worker:
		readAsDaemonizerOrWorker()
	default:
		fmt.Fprintf(os.Stderr, "Circuit role '%s' not recognized\n", Role)
		os.Exit(1)
	}
	if Config.Spark == nil {
		Config.Spark = DefaultSpark
	}
}

func readAsMain() {
	// If CIR is set, it points to a single file that contains all three configuration structures in JSON format.
	cir := os.Getenv("CIR")
	if cir == "" {
		// Otherwise, each one is parsed independently
		parseZookeeper()
		parseInstall()
		parseBuild()
		// Spark is nil when executing as main
		return
	}
	file, err := os.Open(cir)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Problem opening all-in-one config file (%s)", err)
		os.Exit(1)
	}
	defer file.Close()
	parseBag(file)
}

func readAsDaemonizerOrWorker() {
	parseBag(os.Stdin)
}

// WorkerConfig captures the configuration parameters of all sub-systems
// Depending on context of execution, some will be nil.
// Zookeeper and Install should always be non-nil.
type WorkerConfig struct {
	Spark     *SparkConfig
	Zookeeper *ZookeeperConfig
	Install   *InstallConfig
	Build     *BuildConfig
}
var Config *WorkerConfig

func parseBag(r io.Reader) {
	Config = &WorkerConfig{}
	if err := json.NewDecoder(r).Decode(Config); err != nil {
		fmt.Fprintf(os.Stderr, "Problem parsing config (%s)", err)
		os.Exit(1)
	}
	if Config.Install == nil {
		Config.Install = &InstallConfig{}
	}
}
