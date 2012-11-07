package ctrlc

import "tumblr/circuit/kit/debug"

func init() {
	debug.InstallCtrlCPanic()
}
