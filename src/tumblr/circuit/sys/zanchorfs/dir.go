package zanchorfs

import (
	"bytes"
	"encoding/gob"
	"log"
	"path"
	"sync"
	"time"
	"tumblr/circuit/use/lang"
	"tumblr/circuit/use/anchorfs"
	"tumblr/circuit/kit/zookeeper"
	"tumblr/circuit/kit/zookeeper/zutil"
)

/*
	TODO: When a directory changes on Zookeeper, Dir will fetch the entire
	list of children, instea of just the difference. This can be
	inefficient with large directories.
 */

// Dir is responsible for keeping a fresh list of the files in a given Zookeeper directory
type Dir struct {
	fs        *FS
	anchor     string

	sync.Mutex
	stat      *zookeeper.Stat
	watch     *zutil.Watch
	files     map[lang.RuntimeID]*File
	dirs      map[string]struct{}
}

func makeDir(fs *FS, anchor string) (*Dir, error) {
	dir := &Dir{
		fs:     fs,
		anchor: anchor,
	}
	dir.watch = zutil.InstallWatch(fs.zookeeper, dir.zdir())
	// The semantics of AnchorFS pretend that all directories always exist,
	// which is not the case in Zookeeper. To make this work, we create the 
	// directory on access.
	if err := zutil.CreateRecursive(dir.fs.zookeeper, dir.zdir(), zutil.PermitAll); err != nil {
		return nil, err
	}
	return dir, nil
}

func (dir *Dir) zdir() string {
	return path.Join(dir.fs.root, dir.anchor)
}

func (dir *Dir) Name() string {
	return dir.anchor
}

// Files returns the current view of the files in this directory
func (dir *Dir) Files() (rev int64, files map[lang.RuntimeID]anchorfs.File, err error) {
	if err = dir.sync(); err != nil {
		return 0, nil, err
	}
	dir.Lock()
	defer dir.Unlock()
	return dir.rev(), copyFiles(dir.files), nil
}

// dir must be locked before calling rev.
func (dir *Dir) rev() int64 {
	if dir.stat == nil {
		return 0
	}
	return int64(dir.stat.CVersion())
}

func (dir *Dir) Change(sinceRev int64) (rev int64, files map[lang.RuntimeID]anchorfs.File, err error) {
	if err = dir.change(sinceRev, 0); err != nil {
		return 0, nil, err
	}
	dir.Lock()
	defer dir.Unlock()
	return dir.rev(), copyFiles(dir.files), nil
}

func (dir *Dir) ChangeExpire(sinceRev int64, expire time.Duration) (rev int64, files map[lang.RuntimeID]anchorfs.File, err error) {
	if err = dir.change(sinceRev, expire); err != nil {
		return 0, nil, err
	}
	dir.Lock()
	defer dir.Unlock()
	return dir.rev(), copyFiles(dir.files), nil
}

func copyFiles(files map[lang.RuntimeID]*File) map[lang.RuntimeID]anchorfs.File {
	copied := make(map[lang.RuntimeID]anchorfs.File)
	for id, f := range files {
		copied[id] = f
	}
	return copied
}

func (dir *Dir) Dirs() (dirs []string, err error) {
	dirs, _, err = dir.syncDirs()
	return
}

func (dir *Dir) syncDirs() (dirs []string, stat *zookeeper.Stat, err error) {
	if err = dir.sync(); err != nil {
		return nil, nil, err
	}

	dir.Lock()
	dirs = make([]string, 0, len(dir.dirs))
	for anchor, _ := range dir.dirs {
		dirs = append(dirs, anchor)
	}
	stat = dir.stat
	dir.Unlock()

	return dirs, stat, nil
}

// OpenFile returns the worker view of the worker with the specified ID
func (dir *Dir) OpenFile(id lang.RuntimeID) (anchorfs.File, error) {
	if err := dir.sync(); err != nil {
		return nil, err
	}
	dir.Lock()
	file, present := dir.files[id]
	dir.Unlock()
	if !present {
		return nil, anchorfs.ErrNotFound
	}
	return file, nil
}

func (dir *Dir) OpenDir(name string) (anchorfs.Dir, error) {
	return dir.fs.OpenDir(path.Join(dir.anchor, name))
}

func (dir *Dir) change(sinceRev int64, expire time.Duration) error {
	// Check whether the present data is newer than sinceRev
	dir.Lock()
	if dir.rev() > sinceRev {
		dir.Unlock()
		return nil
	}
	stat := dir.stat
	dir.Unlock()

	children, stat, err := dir.watch.ChildrenChange(stat, expire)
	if err != nil {
		if zutil.IsNoNode(err) {
			return nil
		}
		return err
	}

	return dir.fetch(children, stat)
}

// sync updates the files view from Zookeeper, if necessary
func (dir *Dir) sync() error {

	children, stat, err := dir.watch.Children()
	if err != nil {
		if zutil.IsNoNode(err) {
			return nil
		}
		return err
	}

	// If no change since last time, just return
	dir.Lock()
	if dir.stat != nil && dir.stat.CVersion() >= stat.CVersion() {
		dir.Unlock()
		return nil
	}
	dir.Unlock()

	return dir.fetch(children, stat)
}

// fetch refreshes the list of files and subdirectories and applies
// recursive subdirectory pruning as necessary
func (dir *Dir) fetch(children []string, stat *zookeeper.Stat) error {
	dirsNew, filesNew, err := fetch(dir.fs.zookeeper, dir.zdir(), children)
	if err != nil {
		return err
	}
	dir.prune(dirsNew)

	dir.Lock()
	defer dir.Unlock()

	dir.dirs = dirsNew
	dir.files = filesNew
	dir.stat = stat
	return nil
}

// fetch returns the anchor files and subdirectories rooted at zdir
func fetch(z *zookeeper.Conn, zdir string, children []string) (dirs map[string]struct{}, files map[lang.RuntimeID]*File, err error) {
	dirs    = make(map[string]struct{})
	files = make(map[lang.RuntimeID]*File)
	for _, name := range children {
		id, err := lang.ParseRuntimeID(name)
		if err != nil {
			// Node names that are not files are ok. 
			// We treat them as subdirectories.
			dirs[name] = struct{}{}
			continue
		}
		znode := path.Join(zdir, name)
		data, _, err := z.Get(znode)
		if err != nil {
			// If this is a Zookeeper connection error, we should bail out
			log.Printf("Problem getting node `%s` from Zookeeper (%s)", znode, err)
			continue
		}
		zfile := &ZFile{}
		r := bytes.NewBufferString(data)
		if err := gob.NewDecoder(r).Decode(zfile); err != nil {
			log.Printf("anchor file cannot be parsed: (%s)", err)
			continue
		}

		if zfile.Addr.RuntimeID() != id {
			log.Printf("anchor file name vs addr mismatch: %s vs %s\n", id, zfile.Addr.RuntimeID())
			continue
		}
		file := &File{owner: zfile.Addr}
		files[id] = file
	}
	return dirs, files, nil
}

// prune garbage-collects zookeeper anchor directories that have no descendant files in them
func (dir *Dir) prune(dirs map[string]struct{}) error {
	for cname, _ := range dirs {
		cdir_, err := dir.OpenDir(cname)
		if err != nil {
			if zutil.IsNoNode(err) {
				delete(dirs, cname)
				continue
			}
			return err
		}
		cdir := cdir_.(*Dir)
		_, dfiles, err := cdir.Files()
		if err != nil {
			return err
		}
		if len(dfiles) > 0 {
			continue
		}
		ddirs, dstat, err := cdir.syncDirs()
		if err != nil {
			return err
		}
		if len(ddirs) > 0 {
			continue
		}
		delete(dirs, cname)
		if err = dir.fs.zookeeper.Delete(cdir.zdir(), dstat.Version()); err != nil {
			return err
		}
	}
	return nil
}
