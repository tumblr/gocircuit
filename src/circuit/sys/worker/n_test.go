package worker

import (
	"testing"
	//"circuit/kit/debug"
)

func TestKick(t *testing.T) {
	//debug.InstallCtrlCPanic()

	// Start runtime on localhost
	remote, err := SpawnRemote("localhost", 51222, "", "/Users/petar/platform/go/src/tumblr/circuit/sys/rex/srv/srv")
	if err != nil {
		t.Fatalf("spawn (%s)", err)
	}
	defer Kill(remote.Addr)

	// Make local runtime
	r := SpawnLocal()

	// Verify operational
	_, err = r.TryDial(remote.Addr)
	if err != nil {
		t.Errorf("dial (%s)", err)
	}
}
