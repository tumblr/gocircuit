package trace

import (
	"net/http"
	"runtime/pprof"
	"strconv"
)

func init() {
	http.HandleFunc("/_pprof", serveRuntimeProfile)
	http.HandleFunc("/_g", serveGoroutineProfile)
}

func serveGoroutineProfile(w http.ResponseWriter, r *http.Request) {
	prof := pprof.Lookup("goroutine")
	if prof == nil {
		http.Error(w, "unknown profile name", 400)
		return
	}
	prof.WriteTo(w, 1)
}

func serveRuntimeProfile(w http.ResponseWriter, r *http.Request) {
	name := r.URL.Query().Get("n")
	debug, err := strconv.Atoi(r.URL.Query().Get("d"))
	if err != nil {
		http.Error(w, "non-integer or missing debug flag", 400)
		return
	}

	prof := pprof.Lookup(name)
	if prof == nil {
		http.Error(w, "unknown profile name", 400)
		return
	}
	prof.WriteTo(w, debug)
}
