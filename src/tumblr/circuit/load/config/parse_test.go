package config

import (
	"fmt"
	"os"
	"testing"
)

// To run this test, you must set the CIR environment variable first
func TestParse(t *testing.T) {
	fmt.Printf("CIR=%#v %#v %#v\n", Zookeeper, Install, Build)
}
