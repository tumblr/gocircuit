package n

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"os/exec"
	"strconv"
	"circuit/kit/posix"
	"circuit/sys/transport"
	"circuit/use/circuit"
	"circuit/use/n"
	"circuit/load/config"
)

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

func (c *Config) Spawn(host circuit.Host, anchors ...string) (n.Process, error) {

	h := host.(*n.Host).Host
	cmd := exec.Command("ssh", h, "sh")

	stdin, err := cmd.StdinPipe()
	if err != nil {
		return nil, err
	}

	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return nil, err
	}

	stderr, err := cmd.StderrPipe()
	if err != nil {
		return nil, err
	}
	//posix.ForwardStderrBatch(stderr)
	id := circuit.ChooseRuntimeID()
	posix.ForwardStderr(fmt.Sprintf("|%s:stderr> ", id), stderr)

	// Start process
	if err := cmd.Start(); err != nil {
		return nil, err
	}
	defer cmd.Wait() /// Make sure that ssh does not remain zombie

	// Feed shell script to execute circuit binary
	var sh string
	if c.LibPath == "" {
		sh = fmt.Sprintf("%s=%s %s\n", config.RoleEnv, config.Daemonizer, c.Binary)
	} else {
		sh = fmt.Sprintf(
			"LD_LIBRARY_PATH=%s DYLD_LIBRARY_PATH=%s %s=%s %s\n", 
			c.LibPath, c.LibPath, config.RoleEnv, config.Daemonizer, c.Binary)
	}
	stdin.Write([]byte(sh))

	// Write worker configuration to stdin of running worker process
	wc := &config.WorkerConfig{
		Spark: &config.SparkConfig{
			ID:       id,
			BindAddr: "",
			Host:     h,
			Anchor:   append(anchors, fmt.Sprintf("/host/%s", host.String())),
		},
		Zookeeper: config.Config.Zookeeper,
		Install:   config.Config.Install,
	}
	if err := json.NewEncoder(stdin).Encode(wc); err != nil {
		return nil, err
	}

	// Close stdin
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
	//println("n.Spawn -> pid=", pid, "addr=", addr.(*transport.Addr).String())

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

func (c *Config) Kill(remote circuit.Addr) error {
	return kill(remote)
}

func kill(remote circuit.Addr) error {
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
