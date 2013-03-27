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

package server

import (
	"circuit/use/circuit"
	"circuit/use/durablefs"
	"circuit/use/worker"
	"log"
)

func remove(durableFile string) error {
	return durablefs.Remove(durableFile)
}

// Kill reads a sumr service checkpoint from the durableFile in the durable file system and kills the entire service
func Kill(durableFile string) error {
	chk, err := ReadCheckpoint(durableFile)
	if err != nil {
		return circuit.NewError("Problem reading checkpoint (%s)", err)
	}
	if err = killCheckpoint(chk); err != nil {
		return circuit.NewError("Problems killing shards (%s)", err)
	}
	return remove(durableFile)
}

func killCheckpoint(chk *Checkpoint) error {
	var err error
	for i, wrkr := range chk.Workers {
		println("s", i)
		if e := worker.Kill(wrkr.Addr); err != nil {
			log.Printf("Problem killing server %s (%s)", wrkr.Addr, e)
			err = e
		}
	}
	return err
}
