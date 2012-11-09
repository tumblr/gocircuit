package block

import (
	"testing"
	"circuit/kit/fs/diskfs"
)

const testNBlobs = 100

func write(t *testing.T, file File) {
	for i := 0; i < testNBlobs; i++ {
		if n, err := file.Write(encodeUint16(uint16(i))); err != nil || n != 1 {
			t.Fatalf("write n=%d (%s)", n, err)
		}
	}
}

func read(t *testing.T, file File) {
	for i := 0; i < testNBlobs; i++ {
		blob, err := file.Read()
		if err == ErrEndOfLog {
			t.Errorf("expecting k=%d, got k=%d", testNBlobs, i)
			return
		}
		if err != nil {
			t.Fatalf("read (%s)", err)
		}
		if decodeUint16(blob) != uint16(i) {
			t.Errorf("blob value, expect %d, got %d", i, decodeUint16(blob))
		}
	}
}

func TestFile(t *testing.T) {
	disk, err := diskfs.Mount(".", false)
	if err != nil {
		t.Fatalf("mount (%s)", err)
	}

	file, err := Create(disk, "_test_log_file")
	if err != nil {
		t.Fatalf("open (%s)", err)
	}
	write(t, file)
	if err := file.Close(); err != nil {
		t.Errorf("close (%s)", err)
	}

	file, err = Open(disk, "_test_log_file")
	if err != nil {
		t.Fatalf("open2 (%s)", err)
	}
	read(t, file)
	if err := file.Close(); err != nil {
		t.Errorf("close2 (%s)", err)
	}
}
