package n

import (
	"bufio"
	"encoding/gob"
	"fmt"
	"io"
	"os/exec"
	"strconv"
	"strings"
	"tumblr/circuit/use/lang"
	"tumblr/circuit/sys/transport"
	"tumblr/circuit/use/n"
)

func init() {
	gob.Register(&Host{})
}

type Host struct {
	Host string
}

func NewHost(host string) lang.Host {
	return &Host{host}
}

func (h Host) String() string {
	return h.Host
}

type Config struct {
	LibPath string
	Binary  string
	JailDir string
}

func New(libpath, binary, jaildir string) *Config {
	return &Config{
		LibPath: libpath,
		Binary:  binary,
		JailDir: jaildir,
	}
}

func (c *Config) Spawn(host lang.Host, anchors ...string) (n.Process, error) {

	h := host.(*Host).Host
	cmd := exec.Command("ssh", h, "sh")

	stdin, err := cmd.StdinPipe()
	if err != nil {
		return nil, err
	}

	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return nil, err
	}

	// Start process
	if err := cmd.Start(); err != nil {
		return nil, err
	}
	defer cmd.Wait() /// Make sure that ssh does not remain zombie

	// Feed stdin to shell
	id := lang.ChooseRuntimeID()

	ss := fmt.Sprintf(
		"LD_LIBRARY_PATH=%s DYLD_LIBRARY_PATH=%s %s daemonize '%s' '%s' '%s' '%s'\n", 
		c.LibPath, c.LibPath, c.Binary, h, id, c.JailDir, host)
	stdin.Write([]byte(ss))

	// Send internal per-host anchor
	fmt.Fprintf(stdin, "/host/%s\n", host.String())

	// Send user anchors
	for _, a := range anchors {
		fmt.Fprintf(stdin, "%s\n", strings.TrimSpace(a))
	}
	stdin.Write([]byte{'\n'})
	if err = stdin.Close(); err != nil {
		return nil, err
	}

	// Read the first two lines of stdout. They should hold the Port and PID of the runtime process.
	stdoutBuffer := bufio.NewReader(stdout)

	// First line equals PID
	line, err := stdoutBuffer.ReadString('\n')
	if err != nil {
		return nil, err
	}
	line = line[:len(line)-1]
	pid, err := strconv.Atoi(line)
	if err != nil {
		return nil, err
	}

	// Second line equals port
	line, err = stdoutBuffer.ReadString('\n')
	if err != nil {
		return nil, err
	}
	line = line[:len(line)-1]
	port, err := strconv.Atoi(line)
	if err != nil {
		return nil, err
	}

	addr, err := transport.NewAddr(id, pid, fmt.Sprintf("%s:%d", h, port))
	if err != nil {
		return nil, err
	}
	// println("n.Spawn -> pid=", pid, "addr=", addr.(*transport.Addr).String())

	return &Process{
		addr:   addr.(*transport.Addr),
		/*
		console: Console{
			stdin:  stdin,
			stdout: stdout,
			stderr: stderr,
		},
		*/
	}, nil
}

func (c *Config) Kill(remote lang.Addr) error {
	return kill(remote)
}

func kill(remote lang.Addr) error {
	a := remote.(*transport.Addr)
	cmd := exec.Command("ssh", a.Addr.IP.String(), "sh")

	stdinReader, stdinWriter := io.Pipe()
	cmd.Stdin = stdinReader

	if err := cmd.Start(); err != nil {
		return err
	}

	if _, err := fmt.Fprintf(stdinWriter, "kill -KILL %d\n", a.PID); err != nil {
		return err
	}
	stdinWriter.Close()
	
	return cmd.Wait()
}
