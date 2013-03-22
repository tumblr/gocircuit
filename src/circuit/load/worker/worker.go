// Importing package worker has the side effect of turning your program into a circuit worker executable
package worker

import (
	_ "circuit/load"
	_ "circuit/kit/debug/kill"
)

func init() {
	// After package load installs and activates all circuit-related logic,
	// this function blocks forever, never allowing execution of main.
	<-(chan struct{})(nil)
}
