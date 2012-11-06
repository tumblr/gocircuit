package config

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
)

// Role determines the context within which this executable was invoked
const (
	Main = iota
	Daemonizer
	Worker
)
var Role int

func parseRole(s string) int {
	switch s {
	case "":
		return Main
	case "daemonizer":
		return Daemonizer
	case "worker":
		return Worker
	}
	fmt.Fprintf(os.Stderr, "Unrecognized execution role '%s'\n", s)
	os.Exit(1)
	panic("unr")
}

// CIRCUIT_ROLE names the environment variable that determines the role of this invokation
const RoleEnv = "CIRCUIT_ROLE"

// init determines in what context we are being run and reads the configurations accordingly
func init() {
	Config = &WorkerConfig{}
	Role = parseRole(os.Getenv(RoleEnv))
	switch Role {
	case Main:
		readAsMain()
	case Daemonizer:
		readAsDaemonizerOrWorker()
	case Worker:
		readAsDaemonizerOrWorker()
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
}
