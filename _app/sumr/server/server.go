package server

import (
	"os"
	"time"
	"tumblr/app/sumr"
	"tumblr/app/sumr/block"
	"tumblr/circuit/kit/sched/limiter"
	"tumblr/circuit/kit/fs/diskfs"
)

type Server struct {
	block      *block.Block
	lmtr       *limiter.Limiter
}

func New(diskPath string, forgetAfter time.Duration) (*Server, error) {
	s := &Server{}

	os.MkdirAll(diskPath, 0700)

	// Mount disk
	disk, err := diskfs.Mount(diskPath, false)
	if err != nil {
		return nil, err
	}
	// Make db block
	if s.block, err = block.NewBlock(disk, forgetAfter); err != nil {
		return nil, err
	}
	// Prepare incoming call rate-limiter
	s.lmtr = limiter.New(10)
	return s, nil
}

func (s *Server) Add(updateTime time.Time, key sumr.Key, value float64) float64 {
	s.lmtr.Open()
	defer s.lmtr.Close()
	return s.block.Add(updateTime, key, value)
}

func (s *Server) Sum(key sumr.Key) float64 {
	s.lmtr.Open()
	defer s.lmtr.Close()
	return s.block.Sum(key)
}

func (s *Server) Stat() *block.Stat {
	s.lmtr.Open()
	defer s.lmtr.Close()
	return s.block.Stat()
}
