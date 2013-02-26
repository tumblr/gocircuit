package trace

import (
	"net/http"
	"runtime/pprof"
	"strconv"
)

func init() {
	http.HandleFunc("", serveRuntimeProfile)
}

func serveRuntimeProfile(w http.ResponseWriter, r *http.Request) {
	name := r.URL.Query().Get("n")
	debug, err := strconv.Atoi(r.URL.Query().Get("n"))
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
