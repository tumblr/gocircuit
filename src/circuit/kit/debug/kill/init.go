// Package kill has the side effect of installing a KILL signal handler that throws a panic
package kill

import "circuit/kit/debug"

func init() {
	debug.InstallKillPanic()
}
