package lockfile

import (
	"testing"
)

func TestLockFile(t *testing.T) {
	const name = "/tmp/test.lock"
	lock, err := Create(name)
	if err != nil {
		t.Fatalf("create lock (%s)", err)
	}

	if _, err := Create(name); err == nil {
		t.Errorf("re-create lock should not succceed", err)
	}

	if err = lock.Release(); err != nil {
		t.Fatalf("release lock (%s)", err)
	}
}
