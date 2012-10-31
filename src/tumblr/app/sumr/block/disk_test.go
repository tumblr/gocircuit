package block

import (
	"os"
	"testing"
	"tumblr/circuit/facility/fs/diskfs"
)

func TestDiskWriteRead(t *testing.T) {
	// Setup OS dir
	const dirname = "_test_DiskWriteRead"
	os.RemoveAll(dirname)
	os.MkdirAll(dirname, 0700)
	disk, err := diskfs.Mount(dirname, false)
	if err != nil {
		t.Fatalf("mount fs (%s)", err)
	}

	// Mount empty and write
	d, err := Mount(disk)	
	if err != nil {
		t.Fatalf("mount disk (%s)", err)
	}
	file := d.Master()
	if _, err = file.Read(); err != ErrEndOfLog {
		t.Errorf("expecting eof")
	}
	write(t, file)
	if err := d.Unmount(); err != nil {
		t.Errorf("unmount (%s)", err)
	}

	// Mount non-empty and verify contents
	d, err = Mount(disk)	
	if err != nil {
		t.Fatalf("mount2 disk (%s)", err)
	}
	file = d.Master()
	read(t, file)
	if err := d.Unmount(); err != nil {
		t.Errorf("unmount (%s)", err)
	}
}

func TestDiskPromote(t *testing.T) {
	// Setup OS dir
	const dirname = "_test_DiskPromote"
	os.RemoveAll(dirname)
	os.MkdirAll(dirname, 0700)
	disk, err := diskfs.Mount(dirname, false)
	if err != nil {
		t.Fatalf("mount fs (%s)", err)
	}

	// Mount write
	d, err := Mount(disk)	
	if err != nil {
		t.Fatalf("mount disk (%s)", err)
	}
	file, err := d.CreateShadow()
	if err != nil {
		t.Fatalf("create shadow (%s)", err)
	}
	write(t, file)
	if err = d.Promote(file); err != nil {
		t.Fatalf("promote (%s)", err)
	}
	file = d.Master()
	write(t, file)
	if err := d.Unmount(); err != nil {
		t.Errorf("unmount (%s)", err)
	}

	// Mount non-empty and verify contents
	d, err = Mount(disk)	
	if err != nil {
		t.Fatalf("mount2 disk (%s)", err)
	}
	file = d.Master()
	read(t, file)
	read(t, file)
	if err := d.Unmount(); err != nil {
		t.Errorf("unmount (%s)", err)
	}
}
