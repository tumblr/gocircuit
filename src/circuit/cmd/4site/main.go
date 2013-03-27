// Copyright 2013 Tumblr, Inc.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

// 4site is a static web server used to serve the documentation in /doc
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
