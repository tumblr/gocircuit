package config

import "tumblr/circuit/use/circuit"

type SparkConfig struct {
	Cmd      string
	BindAddr string
	ID       circuit.RuntimeID
	Host     string
	Anchor   []string
}

var Spark *SparkConfig
