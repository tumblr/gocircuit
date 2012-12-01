package fmt

import (
	"os"
	"testing"
)

func TestDeep(t *testing.T) {
	s := []interface{}{"a", 2, "c"}
	Deep(os.Stdout, s)
}
