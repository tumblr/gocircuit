// mkcfg prints out an empty sumr shard servers configuration
package main

import (
	"circuit/app/sumr/server"
	"encoding/json"
	"fmt"
	"time"
)

func main() {
	raw, _ := json.MarshalIndent(
		&server.Config{
			Workers: []*server.WorkerConfig{
				&server.WorkerConfig{
					Host:     "host1.datacenter.net",
					DiskPath: "/tmp/sumr",
					Forget:   time.Hour,
				},
				&server.WorkerConfig{
					Host:     "host2.datacenter.net",
					DiskPath: "/tmp/sumr",
					Forget:   time.Hour,
				},
			},
		},
		"", "\t",
	)
	fmt.Printf("%s\n", raw)
}
