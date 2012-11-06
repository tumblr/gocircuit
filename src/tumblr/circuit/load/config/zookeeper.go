package config

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"strings"
)

// ZookeeperConfig holds configuration parameters regarding the zookeeper cluster for the circuit
type ZookeeperConfig struct {
	Workers []string  // Access points for the Zookeeper cluster
	RootDir string    // Root directory for this circuit instance within Zookeeper
}

func (z *ZookeeperConfig) Zookeepers() string {
	var w bytes.Buffer
	for i, u := range z.Workers {
		w.WriteString(u)
		if i + 1 < len(z.Workers) {
			w.WriteByte(',')
		}
	}
	return string(w.Bytes())
}

func (z *ZookeeperConfig) AnchorDir() string {
	return path.Join(z.RootDir, "anchor")
}

func (z *ZookeeperConfig) IssueDir() string {
	return path.Join(z.RootDir, "issue")
}

func (z *ZookeeperConfig) DurableDir() string {
	return path.Join(z.RootDir, "durable")
}

func parseZookeeper() {
	Config.Zookeeper = &ZookeeperConfig{}

	// Try parsing Zookeeper config out of environment variables
	zw := os.Getenv("_CIR_ZW")
	if zw != "" {
		Config.Zookeeper.Workers = strings.Split(zw, ",")
		Config.Zookeeper.RootDir = os.Getenv("_CIR_ZR")
		if Config.Zookeeper.RootDir == "" {
			fmt.Fprintf(os.Stderr, "No Zookeeper root directory in $_CIR_ZR")
			Config.Zookeeper = nil
		}
		return
	}

	// Otherwise, parse Zookeeper config out of a file
	ifile := os.Getenv("CIR_ZOOKEEPER")
	data, err := ioutil.ReadFile(ifile)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Problem reading install file (%s)", err)
		os.Exit(1)
	}
	if err := json.Unmarshal(data, Config.Install); err != nil {
		fmt.Fprintf(os.Stderr, "Problem parsing install file (%s)", err)
		os.Exit(1)
	}
}
