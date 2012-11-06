package load

import (
	"fmt"
	"math/rand"
	"os"
	"strconv"
	"time"

	"tumblr/circuit/kit/lockfile"
	
	"tumblr/circuit/sys/lang"
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
	//"tumblr/circuit/use/n/hijack"
	
	"tumblr/circuit/load/config" // Side-effect of reading in configurations
)


func init() {

	// Seed random number generator
	rand.Seed(time.Now().UnixNano())	

	switch config.Role {
	case config.Main:
		start(false, config.Zookeeper, config.Install, config.Spark)
	case config.Worker:
		start(true, config.Zookeeper, config.Install, config.Spark)
	case config.Daemonizer:
		???
	}
}

func start(worker bool, z *config.ZookeeperConfig, i *config.InstallConfig, s *config.SparkConfig) {
	// If this is a worker, create a lock file in its working directory
	if worker {
		if _, err := lockfile.Create("lock"); err != nil {
			fmt.Fprintf(os.Stderr, "Worker cannot obtain lock (%s)\n", err)
			os.Exit(1)
		}
	}

	// Connect to Zookeeper for anchor file system
	aconn := zanchorfs.Dial(z.Workers)
	anchorfs.Bind(zanchorfs.New(aconn, z.AnchorDir()))

	// Connect to Zookeeper for durable file system
	dconn := zdurablefs.Dial(z.Workers)
	durablefs.Bind(zdurablefs.New(dconn, z.DurableDir()))

	// Connect to Zookeeper for issue file system
	iconn := zissuefs.Dial(z.Workers)
	issuefs.Bind(zissuefs.New(iconn, z.IssueDir()))

	// Initialize the networking module
	un.Bind(sn.New(i.LibPath, i.Binary, i.JailDir()))

	// Initialize transport module
	t := transport.New(s.ID, s.BindAddr, s.Host)

	// Initialize language runtime
	circuit.Bind(lang.New(t))

	// Create anchors
	for _, a := range s.Anchor {
		if err := anchorfs.CreateFile(a, t.Addr()); err != nil {
			fmt.Fprintf(os.Stderr, "Problem creating anchor '%s' (%s)\n", a, err)
			os.Exit(1)
		}
	}

	if worker {
		// A worker sends back its PID and runtime port to its invoker (the daemonizer)
		backpipe := os.NewFile(3, "backpipe")
		if _, err := backpipe.WriteString(strconv.Itoa(os.Getpid()) + "\n"); err != nil {
			panic(err)
		}
		if _, err := backpipe.WriteString(strconv.Itoa(t.Port()) + "\n"); err != nil {
			panic(err)
		}
		if err := backpipe.Close(); err != nil {
			panic(err)
		}

		// Hang forever
		<-(chan int)(nil)
	}
}

func NewHost(h string) circuit.Host {
	return sn.NewHost(h)
}

func NewTransport(id circuit.RuntimeID, addr, host string) circuit.Transport {
	return transport.New(id, addr, host)
}
