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

package front

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
