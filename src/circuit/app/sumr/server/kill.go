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
