package zdurablefs

import (
	"path"
	"circuit/use/circuit"
	"circuit/kit/zookeeper"
)

var (
	ErrClosed  = circuit.NewError("closed")
)

// FS implements a durable file system on top of Zookeeper
type FS struct {
	conn     *zookeeper.Conn
	zroot    string
}

func New(conn *zookeeper.Conn, zroot string) *FS {
	return &FS{conn: conn, zroot: zroot}
}

func (fs *FS) Remove(fpath string) error {
	return fs.conn.Delete(path.Join(fs.zroot, fpath), -1)
}
