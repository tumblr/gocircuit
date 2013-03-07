// Package firehose implements a connection to the Tumblr Firehose
package firehose

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"net"
	"net/http"
	"net/http/httputil"
	"net/textproto"
	"net/url"
	"os"
	"strings"
	"time"
)

// Example HTTP request to the Firehose:
//
//	GET /?applicationId=1&offset=oldest&clientId=87 HTTP/1.1
//	Authorization: Basic Ym1hdGhlbnk6Zm9vYmFyYmF6YnV6
//	User-Agent: curl/7.21.4 (universal-apple-darwin11.0) libcurl/7.21.4 OpenSSL/0.9.8r zlib/1.2.5
//	Host: localhost:8000
//	Accept: */*

type Request struct {
	HostPort       string
	Username       string
	Password       string
	ApplicationID  string
	ClientID       string
	Offset         string
}

func MakeRequest(freq *Request) *http.Request {
	req, err := http.NewRequest("GET", "/", nil)
	if err != nil {
		panic("make firehose request")
	}
	args := url.Values{}
	args.Set("applicationId", freq.ApplicationID)
	args.Set("offset", freq.Offset)
	args.Set("clientId", freq.ClientID)
	req.URL = &url.URL{
		Scheme:   "http",
		Host:     freq.HostPort,
		Path:     "/",
		RawQuery: args.Encode(),
	}
	req.Host = freq.HostPort
	req.SetBasicAuth(freq.Username, freq.Password)
	return req
}

// Conn is a connection to the Tumblr Firehose with a given application, offset and client parameters
type Conn struct {
	resp *http.Response
	r    *textproto.Reader
}

// Dial connects to the firehose and returns a connection object capable of reading Firehose events iteratively
func Dial(freq *Request) (*Conn, error) {
	client := &http.Client{
		Transport: &transport{},
	}
	resp, err := client.Do(MakeRequest(freq))
	if err != nil {
		return nil, err
	}
	return &Conn{
		resp: resp,
		r:    textproto.NewReader(bufio.NewReader(resp.Body)),
	}, nil
}

// Read reads the next Firehose event into the supplied value
func (conn *Conn) ReadInterface(v interface{}) error {
	line, err := conn.r.ReadLine()
	if err != nil {
		return err
	}
	if err = json.Unmarshal([]byte(line), v); err != nil {
		fmt.Fprintf(os.Stderr, "firehose non-json response:\n= %s\n", line)
		return err
	}
	return nil
}

// Read reads the next Firehose event and returns the results parsed into an Event structure
func (conn *Conn) Read() (*Event, error) {
	m := make(map[string]interface{})
	if err := conn.ReadInterface(&m); err != nil {
		return nil, err
	}
	return ParseEvent(m)
}

// Read reads the next Firehose event in raw, non-decoded string form
func (conn *Conn) ReadRaw() (string, error) {
	return conn.r.ReadLine()
}

// Close closes the connection to the Firehose
func (conn *Conn) Close() error {
	return conn.resp.Body.Close()
}

// transport is a special http.RoundTripper designed for the streaming nature of the Firehose
type transport struct {}

func (t *transport) RoundTrip(req *http.Request) (resp *http.Response, err error) {
	conn, err := net.DialTimeout("tcp", canonicalAddr(req.URL), 2 * time.Second)
	if err != nil {
		return nil, err
	}
	cc := httputil.NewClientConn(newTimeoutConn(conn, 3 * time.Second), nil)
	resp, err = cc.Do(req)
	if err != nil {
		return nil, err
	}
	resp.Body = &disconnectOnBodyClose{ resp.Body, cc }
	return resp, nil
}

type disconnectOnBodyClose struct {
	io.ReadCloser
	clientConn *httputil.ClientConn
}

func (d *disconnectOnBodyClose) Close() error {
	err := d.clientConn.Close()
	d.ReadCloser.Close()
	return err
}

// canonicalAddr returns url.Host but always with a ":port" suffix
func canonicalAddr(url *url.URL) string {
	addr := url.Host
	if !hasPort(addr) {
		return addr + ":80"
	}
	return addr
}

// Given a string of the form "host", "host:port", or "[ipv6::address]:port",
// return true if the string includes a port.
func hasPort(s string) bool { return strings.LastIndex(s, ":") > strings.LastIndex(s, "]") }

type timeoutConn struct {
	net.Conn
	timeout time.Duration
}

func newTimeoutConn(conn net.Conn, timeout time.Duration) net.Conn {
	return &timeoutConn{
		Conn:    conn,
		timeout: timeout,
	}
}

func (t *timeoutConn) Read(b []byte) (n int, err error) {
	t.Conn.SetReadDeadline(time.Now().Add(t.timeout))
	return t.Conn.Read(b)
}

func (t *timeoutConn) Write(b []byte) (n int, err error) {
	t.Conn.SetWriteDeadline(time.Now().Add(t.timeout))
	return t.Conn.Write(b)
}
