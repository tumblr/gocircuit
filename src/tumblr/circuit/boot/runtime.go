package boot

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
	"tumblr/circuit/use/lang"
	un "tumblr/circuit/use/n"
	
	// _ "tumblr/TUMBLR/app" // Registers all apps used by the environment
)


func init() {
	parseZookeeperConfig()
	parseInstallConfig()
	parseBuildConfig()

	// Seed random number generator
	rand.Seed(time.Now().UnixNano())
	
	// Connect to Zookeeper for anchor file system
	aconn := zanchorfs.Dial(Zookeeper.Workers)
	anchorfs.Bind(zanchorfs.New(aconn, Zookeeper.AnchorDir()))

	// Connect to Zookeeper for durable file system
	dconn := zdurablefs.Dial(Zookeeper.Workers)
	durablefs.Bind(zdurablefs.New(dconn, Zookeeper.DurableDir()))

	// Connect to Zookeeper for issue file system
	iconn := zissuefs.Dial(Zookeeper.Workers)
	issuefs.Bind(zissuefs.New(iconn, Zookeeper.IssueDir()))

	// Create network
	un.Bind(sn.New(Install.LibPath, Install.Binary, Install.JailDir()))
}

func NewHost(h string) lang.Host {
	return sn.NewHost(h)
}

func NewTransport(id lang.RuntimeID, addr, host string) lang.Transport {
	return transport.New(id, addr, host)
}
