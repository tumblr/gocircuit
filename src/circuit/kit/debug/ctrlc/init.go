// Package ctrlc has the side effect of installing a Ctrl-C signal handler that throws a panic
package ctrlc

import "circuit/kit/debug"

func init() {
	debug.InstallCtrlCPanic()
}
