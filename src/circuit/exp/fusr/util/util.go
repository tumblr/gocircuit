package util

import (
	"sync"
	"tumblr/database/levigo"
)


type Server struct {
	slk          sync.Mutex
	cache        *levigo.Cache
	DB           *levigo.DB
	ReadAndCache *levigo.ReadOptions
	WriteSync    *levigo.WriteOptions
	WriteNoSync  *levigo.WriteOptions
}

func (srv *Server) Init(dbDir string, cacheSize int) error {
	var err error
	opts := levigo.NewOptions()
	srv.cache = levigo.NewLRUCache(cacheSize)
	opts.SetCache(srv.cache)
	opts.SetCreateIfMissing(true)

	if srv.DB, err = levigo.Open(dbDir, opts); err != nil {
		srv.cache.Close()
		return err
	}

	srv.ReadAndCache = levigo.NewReadOptions()
	srv.ReadAndCache.SetFillCache(true)

	srv.WriteSync = levigo.NewWriteOptions()
	srv.WriteSync.SetSync(true)

	srv.WriteNoSync = levigo.NewWriteOptions()
	srv.WriteSync.SetSync(false)
	
	return nil
}

func (srv *Server) Close() error {
	srv.slk.Lock()
	defer srv.slk.Unlock()
	if srv.cache != nil {
		srv.cache.Close()
		srv.cache = nil
	}
	srv.ReadAndCache.Close()
	srv.ReadAndCache = nil
	srv.WriteSync.Close()
	srv.WriteSync = nil
	srv.WriteNoSync.Close()
	srv.WriteNoSync = nil
	return nil
}

