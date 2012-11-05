package config

import "tumblr/circuit/use/circuit"

type SparkConfig struct {
	ID       circuit.RuntimeID
	BindAddr string
	Host     string
	Anchor   []string
}

var Spark *SparkConfig

var DefaultSpark = &SparkConfig{
	ID:       circuit.ChooseRuntimeID(),
	BindAddr: "",
	Host:     "",
	Anchor:   []string{},
}
