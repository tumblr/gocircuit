package zdurablefs

import (
	"path"
	"sync"
	"time"
	"tumblr/circuit/kit/zookeeper"
	"tumblr/circuit/kit/zookeeper/zutil"
	"tumblr/circuit/use/durablefs"
)

// TODO: Add directory garbage collection

type Dir struct {
	conn     *zookeeper.Conn
	zroot    string
	dpath    string

	sync.Mutex
	stat     *zookeeper.Stat
	watch    *zutil.Watch
	children map[string]struct{}
}

func (fs *FS) OpenDir(dpath string) durablefs.Dir {
	return &Dir{
		conn:  fs.conn,
		zroot: fs.zroot,
		dpath: dpath,
		watch: zutil.InstallWatch(fs.conn, path.Join(fs.zroot, dpath)),
	}
}

func (dir *Dir) Path() string {
	return dir.dpath
}

// Children returns the children in this directory.
func (dir *Dir) Children() map[string]struct{} {
	if err := dir.sync(); err != nil {
		panic(err)
	}
	dir.Lock()
	defer dir.Unlock()
	return copyChildren(dir.children)
}

// Change blocks until a change in the children of dir occurs.
func (dir *Dir) Change() map[string]struct{} {
	if err := dir.change(0); err != nil {
		panic(err)
	}
	dir.Lock()
	defer dir.Unlock()
	return copyChildren(dir.children)
}

// Expire is ..
func (dir *Dir) Expire(expire time.Duration) map[string]struct{} {
	if err := dir.change(expire); err != nil {
		panic(err)
	}
	dir.Lock()
	defer dir.Unlock()
	return copyChildren(dir.children)
}

func copyChildren(w map[string]struct{}) map[string]struct{} {
	children := make(map[string]struct{})
	for k, _ := range w {
		children[k] = struct{}{}
	}
	return children
}

// Close closes the watches on this diretory
func (dir *Dir) Close() {
	dir.Lock()
	defer dir.Unlock()

	if dir.conn == nil {
		panic(ErrClosed)
	}
	dir.conn = nil

	dir.watch.Close()
	dir.watch = nil
}

func (dir *Dir) change(expire time.Duration) error {
	dir.Lock()
	watch := dir.watch
	dir.Unlock()
	if watch == nil {
		return ErrClosed
	}

	dir.Lock()
	stat := dir.stat
	dir.Unlock()

	children, stat, err := dir.watch.ChildrenChange(stat, expire)
	if err != nil {
		if zutil.IsNoNode(err) {
			return nil
		}
		return err
	}
	dir.update(children, stat)
	return nil
}

// sync updates the files view from Zookeeper, if necessary
func (dir *Dir) sync() error {
	dir.Lock()
	watch := dir.watch
	dir.Unlock()
	if watch == nil {
		return ErrClosed
	}
	children, stat, err := dir.watch.Children()
	if err != nil {
		if zutil.IsNoNode(err) {
			// No node, no problem. We represent it as present and empty.
			return nil
		}
		return err
	}
	dir.update(children, stat)
	return nil
}

func (dir *Dir) update(children []string, stat *zookeeper.Stat) {
	// If no change since last time, just return
	dir.Lock()
	defer dir.Unlock()

	if dir.stat != nil && dir.stat.CVersion() >= stat.CVersion() {
		return
	}
	dir.stat = stat
	dir.children = make(map[string]struct{})
	for _, c := range children {
		dir.children[c] = struct{}{}
	}
}
