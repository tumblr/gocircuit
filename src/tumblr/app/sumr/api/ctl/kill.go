package ctl

import (
	"log"
	"strconv"
	"tumblr/circuit/use/anchorfs"
	"tumblr/circuit/use/n"
)

func Kill(c *Config) error {
	dir, e := anchorfs.OpenDir(c.Anchor)
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
			if e = n.Kill(f.Owner()); e != nil {
				log.Printf("Problem killing %s (%s)", f.Owner(), e)
			}
			break
		}
	}
	return nil
}
