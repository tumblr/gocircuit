// Package file provides ways to pass open files to across circuit runtimes
package io

import (
	"runtime"
	"circuit/use/circuit"
	"io"
)

func NewClient(x circuit.X) *Client {
	return &Client{X: x}
}

type Client struct {
	circuit.X
}

func asError(x interface{}) error {
	if x == nil {
		return nil
	}
	return x.(error)
}

func asBytes(x interface{}) []byte {
	if x == nil {
		return nil
	}
	return x.([]byte)
}

func _recover(pe *error) {
	if p := recover(); p != nil {
		*pe = circuit.NewError("server died")
	}
}

func (cli *Client) Close() (err error) {
	defer _recover(&err)

	return asError(cli.Call("Close")[0])
}

func (cli *Client) Read(p []byte) (_ int, err error) {
	defer _recover(&err)

	r := cli.Call("Read", len(p))
	q, err := asBytes(r[0]), asError(r[1])
	if len(q) > len(p) {
		panic("corrupt i/o server")
	}
	copy(p, q)
	return len(q), err
}

func (cli *Client) Write(p []byte) (_ int, err error) {
	defer _recover(&err)

	r := cli.Call("Write", p)
	return r[0].(int), asError(r[1])
}

func NewServer(f io.ReadWriteCloser) *Server {
	srv := &Server{f: f}
	runtime.SetFinalizer(srv, func(srv_ *Server) {
		srv.f.Close()
	})
	return srv
}

type Server struct {
	f io.ReadWriteCloser
}

func init() {
	circuit.RegisterValue(&Server{})
}

func (srv *Server) Close() error {
	return srv.f.Close()
}

func (srv *Server) Read(n int) ([]byte, error) {
	p := make([]byte, min(n, 1e4))
	m, err := srv.f.Read(p)
	return p[:m], err
}

func min (x, y int) int {
	if x < y {
		return x
	}
	return y
}

func (srv *Server) Write(p []byte) (int, error) {
	return srv.f.Write(p)
}
