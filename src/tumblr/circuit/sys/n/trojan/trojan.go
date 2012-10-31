package trojan

import (
	"bufio"
	"fmt"
	"io"
	"math/rand"
	"os"
	"os/exec"
	"path"
	"strconv"
	"strings"
	"time"
	"tumblr/circuit/kit/debug"
	"tumblr/circuit/kit/lockfile"
	"tumblr/circuit/sys/lang"
	"tumblr/circuit/sys/transport"
	"tumblr/circuit/use/anchorfs"
	"tumblr/circuit/use/circuit"
)

func usage() {
	fmt.Fprintf(os.Stderr, "usage: %s [daemonize|run] Addr RuntimeID JailDir Host\n", os.Args[0])
	os.Exit(1)
}

type NewTransportFunc func(id circuit.RuntimeID, addr, host string) circuit.Transport

// Main is the 'func main' of a circuit binary that can spawn a circuit runtime
// and report back its address.
//
func Main(newTransport NewTransportFunc) {
	debug.InstallCtrlCPanic()
	rand.Seed(time.Now().UnixNano())
	if len(os.Args) != 6 {
		usage()
	}
	switch os.Args[1] {
	case "run":
		run(newTransport, os.Args[2], os.Args[3], os.Args[4], os.Args[5])
	case "daemonize":
		daemonize(os.Args[2], os.Args[3], os.Args[4], os.Args[5])
	}
	usage()
}

func pie(err interface{}) {
	if err != nil {
		panic(err)
	}
}

func pie2(underscore interface{}, err interface{}) {
	pie(err)
}

func piefwd(stdout, stderr *os.File, err interface{}) {
	if err == nil {
		return
	}
	fmt.Fprintf(os.Stderr, "PANIC\n")
	os.Stderr.WriteString("Standard output:\n")
	stdout.Seek(0, 0)
	io.Copy(os.Stderr, stdout)
	os.Stderr.WriteString("Standard error:\n")
	stderr.Seek(0, 0)
	io.Copy(os.Stderr, stderr)
	os.Stderr.WriteString("Daemonizer error:\n")
	panic(err)
}


// dbg is for debugging purposes.
// It provides an alternate way of logging.
func dbg(n, s string) {
	cmd := exec.Command("sh")
	stdin, err := cmd.StdinPipe()
	if err != nil {
		panic("huh")
	}
	cmd.Start()
	defer cmd.Wait()
	fmt.Fprintf(stdin, "echo '%s' >> /Users/petar/tmp/%s\n", s, n)
	stdin.Close()
}


// Run starts the runtime server in this process and blocks.
func run(newTransport NewTransportFunc, addr, id, magic, host string) {
	// Avoid manual execution
	if magic != "cleaver" {
		panic("missing magic")
	}

	// Read anchors from stdin
	r := bufio.NewReader(os.Stdin)
	var anchor []string
	for {
		line, err := r.ReadString('\n')
		pie(err)
		line = strings.TrimSpace(line)
		if line == "" {
			break
		}
		anchor = append(anchor, line)
	}

	// Create lock file
	pie2(lockfile.Create("lock"))

	// Start runtime
	id_, err := circuit.ParseRuntimeID(id)
	pie(err)
	t := newTransport(id_, addr, host)
	circuit.Bind(lang.New(t))

	// Create anchors
	for _, a := range anchor {
		pie(anchorfs.CreateFile(a, t.Addr()))
	}

	// Send back PID and port
	backpipe := os.NewFile(3, "circuit•backpipe")
	pie2(backpipe.WriteString(strconv.Itoa(os.Getpid()) + "\n"))
	pie2(backpipe.WriteString(strconv.Itoa(t.(*transport.Transport).Port()) + "\n"))

	pie(backpipe.Close())

	// Hang forever
	<-(chan int)(nil)
}

func daemonize(addr, id, jaildir, host string) {

	// Make jail directory
	id_ := circuit.ParseOrHashRuntimeID(id)
	jail := path.Join(jaildir, id_.String())
	pie(os.MkdirAll(jail, 0700))

	// Prepare exec
	cmd := exec.Command(os.Args[0], "run", addr, id_.String(), "cleaver", host)
	cmd.Dir = jail

	// Out-of-band pipe for reading child PID and port
	bpr, bpw, err := os.Pipe()
	pie(err)
	cmd.ExtraFiles = []*os.File{bpw}

	// stdin 
	// Relay stdin of daemonizer to stdin of child runtime process
	cmd.Stdin = os.Stdin
	defer os.Stdin.Close()

	// stdout
	stdout, err := os.Create(path.Join(jail, "out"))
	if err != nil {
		panic(err)
	}
	defer stdout.Close()
	cmd.Stdout = stdout

	// stderr
	stderr, err := os.Create(path.Join(jail, "err"))
	if err != nil {
		panic(err)
	}
	defer stderr.Close()
	cmd.Stderr = stderr

	// start
	pie(cmd.Start())
	go func() {
		cmd.Wait()
		piefwd(stdout, stderr, bpw.Close())
	}()
	
	// Read the first two lines of stdout. They should hold the Port and PID of the runtime process.
	back := bufio.NewReader(bpr)

	// Read PID
	line, err := back.ReadString('\n')
	piefwd(stdout, stderr, err)

	pid, err := strconv.Atoi(strings.TrimSpace(line))
	piefwd(stdout, stderr, err)

	// Read port
	line, err = back.ReadString('\n')
	piefwd(stdout, stderr, err)
	port, err := strconv.Atoi(strings.TrimSpace(line))
	piefwd(stdout, stderr, err)

	// Close the pipe
	piefwd(stdout, stderr, bpr.Close())

	if cmd.Process.Pid != pid {
		piefwd(stdout, stderr, "pid mismatch")
	}

	fmt.Printf("%d\n%d\n", pid, port)
	// Sync is not supported on os.Stdout, at least on OSX
	// os.Stdout.Sync()

	// dbg("d", "daemonize succeeded!")
	os.Exit(0)
}
