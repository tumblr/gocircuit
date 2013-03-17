// mkcfg prints out an empty sumr api configuration
package main

import (
	"circuit/app/sumr/api"
	"encoding/json"
	"fmt"
)

func main() {
	raw, _ := json.MarshalIndent(
		&api.Config{
			Anchor: "",
			Workers: []*api.WorkerConfig{
				&api.WorkerConfig{
					Host: "host3.datacenter.net",
					Port: 4000,
				},
				&api.WorkerConfig{
					Host: "host4.datacenter.net",
					Port: 4000,
				},
			},
		},
		"", "\t",
	)
	fmt.Printf("%s\n", raw)
}
