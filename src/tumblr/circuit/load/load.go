package load

import (
	"math/rand"
	"time"
	
	sn "tumblr/circuit/sys/n"
	"tumblr/circuit/sys/transport"
	"tumblr/circuit/sys/zanchorfs"
	"tumblr/circuit/sys/zdurablefs"
	"tumblr/circuit/sys/zissuefs"

	"tumblr/circuit/use/anchorfs"
	"tumblr/circuit/use/durablefs"
	"tumblr/circuit/use/issuefs"
	"tumblr/circuit/use/circuit"
	un "tumblr/circuit/use/n"
	"tumblr/circuit/use/n/hijack"
	
	"tumblr/circuit/load/config"
)


func init() {

	// Seed random number generator
	rand.Seed(time.Now().UnixNano())	

	if config.Worker {
		hijack.Main(NewTransport)
		panic("never reach")
	}

	// Connect to Zookeeper for anchor file system
	aconn := zanchorfs.Dial(config.Zookeeper.Workers)
	anchorfs.Bind(zanchorfs.New(aconn, config.Zookeeper.AnchorDir()))

	// Connect to Zookeeper for durable file system
	dconn := zdurablefs.Dial(config.Zookeeper.Workers)
	durablefs.Bind(zdurablefs.New(dconn, config.Zookeeper.DurableDir()))

	// Connect to Zookeeper for issue file system
	iconn := zissuefs.Dial(config.Zookeeper.Workers)
	issuefs.Bind(zissuefs.New(iconn, config.Zookeeper.IssueDir()))

	// Create network
	un.Bind(sn.New(config.Install.LibPath, config.Install.Binary, config.Install.JailDir()))
}

func NewHost(h string) circuit.Host {
	return sn.NewHost(h)
}

func NewTransport(id circuit.RuntimeID, addr, host string) circuit.Transport {
	return transport.New(id, addr, host)
}
