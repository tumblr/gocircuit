// 4site is static web server, used to serve the documentation in /doc
package main

import (
	"flag"
	"net/http"
)

var (
	flagHTTP = flag.String("http", ":2020", "Bind address for HTTP server")
	flagPath = flag.String("path", "", "Path to the /doc directory within the Go Circuit repo")
)

func main() {
	flag.Parse()
	http.Handle("/", http.FileServer(http.Dir(*flagPath)))
	http.ListenAndServe(*flagHTTP, nil)
}
