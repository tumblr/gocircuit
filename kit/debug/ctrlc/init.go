package ctrlc

import "circuit/kit/debug"

func init() {
	debug.InstallCtrlCPanic()
}
