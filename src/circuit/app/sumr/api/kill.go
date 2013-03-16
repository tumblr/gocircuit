package api

import (
	"circuit/use/anchorfs"
	"circuit/use/worker"
	"log"
	"strconv"
)

// Kill kills any live workers belonging to an API started with config
func Kill(config *Config) error {
	dir, e := anchorfs.OpenDir(config.Anchor)
	if e != nil {
		return e
	}
	dirs, e := dir.Dirs()
	if e != nil {
		return e
	}
	for _, d := range dirs {
		_, e := strconv.Atoi(d)
		if e != nil {
			continue
		}
		wdir, e := dir.OpenDir(d)
		if e != nil {
			return e
		}
		_, files, e := wdir.Files()
		if e != nil {
			return e
		}
		for _, f := range files {
			if e = worker.Kill(f.Owner()); e != nil {
				log.Printf("Problem killing %s (%s)", f.Owner(), e)
			}
			break
		}
	}
	return nil
}
