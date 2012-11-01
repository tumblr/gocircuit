package config

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
)

type allInOne struct {
	Zookeeper *ZookeeperConfig
	Install   *InstallConfig
	Build     *BuildConfig
}

func init() {
	cc := os.Getenv("CIR")
	if cc != "" {
		data, err := ioutil.ReadFile(cc)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Problem reading all-in-one config file (%s)", err)
			os.Exit(1)
		}
		all := &allInOne{}
		if err := json.Unmarshal(data, all); err != nil {
			fmt.Fprintf(os.Stderr, "Problem parsing all-in-one config file (%s)", err)
			os.Exit(1)
		}
		Zookeeper = all.Zookeeper
		Install = all.Install
		Build = all.Build
		return
	}

	parseZookeeper()
	parseInstall()
	parseBuild()
}
