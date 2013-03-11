package kafka

import (
	"fmt"
	"testing"
)

func TestEcosystem(t *testing.T) {
	eco, err := NewEcosystem("127.0.0.1:2181")
	if err != nil {
		t.Fatalf("connect to Zookeeper (%s)", err)
	}
	brokers, err := eco.Brokers()
	if err != nil {
		t.Fatalf("get brokers (%s)", err)
	}
	for _, be := range brokers {
		fmt.Printf("%s\n", be)
	}
}
