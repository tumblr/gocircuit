package zipfs

import (
	"archive/zip"
	"io"
	"os"
	"path"
	"strings"
	"sync"
	"time"
	fspkg "tumblr/circuit/kit/fs"
)

// FS provides a read-only file system access to a local zip file
type FS struct {
	r    *zip.ReadCloser
	root *Dir
}

func Mount(zipfile string) (fs *FS, err error) {
	fs = &FS{root: NewDir("") }
	if fs.r, err = zip.OpenReader(zipfile); err != nil {
		return nil, err
	}
	for _, f := range fs.r.File {
		fs.root.addFile(splitPath(f.Name), f)
	}
	return fs, nil
}

func splitPath(p string) []string {
	parts := strings.Split(path.Clean(p), "/")
	if len(parts) > 0 && parts[0] == "" {
		parts = parts[1:]
	}
	return parts
}

func (fs *FS) Close() error {
	return fs.r.Close()
}

func (fs *FS) Open(name string) (fspkg.File, error) {
	return fs.root.openFile(splitPath(name))
}

func (fs *FS) OpenFile(name string, flag int, perm os.FileMode) (fspkg.File, error) {
	panic("not supported")
}

func (fs *FS) Create(name string) (fspkg.File, error) {
	return nil, fspkg.ErrReadOnly
}

func (fs *FS) Remove(name string) error {
	return fspkg.ErrReadOnly
}

func (fs *FS) Rename(oldname, newname string) error {
	return fspkg.ErrReadOnly
}

func (fs *FS) Stat(name string) (os.FileInfo, error) {
	return fs.root.statFile(splitPath(name))
}

func (fs *FS) Mkdir(name string) error {
	return fspkg.ErrReadOnly
}

func (fs *FS) MkdirAll(name string) error {
	return fspkg.ErrReadOnly
}

// Dir is an open directory in a zip archive
type Dir struct {
	name  string
	lk    sync.Mutex
	dirs  map[string]*Dir
	files map[string]*zip.File
}

func NewDir(name string) *Dir {
	return &Dir{
		name:  name,
		dirs:  make(map[string]*Dir),
		files: make(map[string]*zip.File),
	}
}

func (d *Dir) addFile(parts []string, file *zip.File) error {
	d.lk.Lock()
	defer d.lk.Unlock()

	if len(parts) == 1 {
		d.files[parts[0]] = file
		return nil
	}
	sub := d.dirs[parts[0]]
	if sub == nil {
		sub = NewDir(parts[0])
		d.dirs[parts[0]] = sub
	}
	return sub.addFile(parts[1:], file)
}

func (d *Dir) openFile(parts []string) (fspkg.File, error) {
	d.lk.Lock()
	defer d.lk.Unlock()

	if len(parts) == 1 {
		file := d.files[parts[0]]
		if file == nil {
			return nil, fspkg.ErrNotFound
		}
		return OpenFile(file)
	}
	sub := d.dirs[parts[0]]
	if sub == nil {
		return nil, fspkg.ErrNotFound
	}
	return sub.openFile(parts[1:])
}

func (d *Dir) statFile(parts []string) (os.FileInfo, error) {
	d.lk.Lock()
	defer d.lk.Unlock()

	if len(parts) == 1 {
		file := d.files[parts[0]]
		if file == nil {
			return nil, fspkg.ErrNotFound
		}
		return file.FileInfo(), nil
	}
	sub := d.dirs[parts[0]]
	if sub == nil {
		return nil, fspkg.ErrNotFound
	}
	return sub.statFile(parts[1:])
}

func (d *Dir) Close() error {
	return nil
}

func (d *Dir) Stat() (os.FileInfo, error) {
	return &fspkg.FileInfo{
		XName:    d.name,
		XSize:    0,
		XMode:    0700,
		XModTime: time.Time{},
		XIsDir:   true,
	}, nil
}

func (d *Dir) Readdir(count int) ([]os.FileInfo, error) {
	d.lk.Lock()
	defer d.lk.Unlock()
	ls := make([]os.FileInfo, 0, len(d.dirs) + len(d.files))
	for _, dir := range d.dirs {
		fi, _ := dir.Stat()
		ls = append(ls, fi)
	}
	for _, f := range d.files {
		ls = append(ls, f.FileInfo())
	}
	return ls, nil
}

func (d *Dir) Read(p []byte) (int, error) {
	return 0, fspkg.ErrOp
}

func (d *Dir) Seek(offset int64, whence int) (int64, error) {
	return 0, fspkg.ErrOp
}

func (d *Dir) Truncate(size int64) error {
	return fspkg.ErrOp
}

func (d *Dir) Write(q []byte) (int, error) {
	return 0, fspkg.ErrOp
}

func (d *Dir) Sync() error {
	return fspkg.ErrReadOnly
}

// File represents an open file from the zip archive fs
type File struct {
	file *zip.File
	rc   io.ReadCloser
}

func OpenFile(file *zip.File) (fspkg.File, error) {
	rc, err := file.Open()
	if err != nil {
		return nil, err
	}
	return &File{file, rc}, nil
}

func (f *File) Close() error {
	return f.rc.Close()
}

func (f *File) Stat() (os.FileInfo, error) {
	return f.file.FileInfo(), nil
}

func (f *File) Readdir(count int) ([]os.FileInfo, error) {
	return nil, fspkg.ErrOp
}

func (f *File) Read(p []byte) (int, error) {
	return f.rc.Read(p)
}

func (f *File) Seek(offset int64, whence int) (int64, error) {
	return 0, fspkg.ErrOp
}

func (f *File) Truncate(size int64) error {
	return fspkg.ErrReadOnly
}

func (f *File) Write(q []byte) (int, error) {
	return 0, fspkg.ErrReadOnly
}

func (f *File) Sync() error {
	return fspkg.ErrReadOnly
}
