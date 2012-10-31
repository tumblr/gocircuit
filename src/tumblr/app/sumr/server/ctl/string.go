package ctl

import (
	"bytes"
	"fmt"
)

func (s *Checkpoint) String() string {
	var w bytes.Buffer
	for i, shc := range s.Config.Workers {
		srvstr := "•"
		key    := "•"
		if shs := s.Workers[i]; shs != nil {
			srvstr = shs.Server.String()
			key = shs.Key.String()
		}
		fmt.Fprintf(&w, "KEY=%s SERVER=%s HOST=%s DISK=%s FORGET=%s\n", key, srvstr, shc.Host, shc.DiskPath, shc.Forget)
	}
	return string(w.Bytes())
}
