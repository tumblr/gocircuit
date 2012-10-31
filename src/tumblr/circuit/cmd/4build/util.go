package main

import (
	"bytes"
	"fmt"
	"io"
	"math/rand"
	"os"
	"os/exec"
	"path"
	"strconv"
	"strings"
)

func Printf(fmt_ string, arg_ ...interface{}) {
	fmt.Printf(fmt_, arg_...)
}

func Errorf(fmt_ string, arg_ ...interface{}) {
	fmt.Fprintf(os.Stderr, fmt_, arg_...)
}

func Fatalf(fmt_ string, arg_ ...interface{}) {
	Errorf(fmt_, arg_...)
	os.Exit(1)
}

func MakeTempDir() (string, error) {
	tempRoot := os.TempDir()
	abs := path.Join(tempRoot, strconv.FormatInt(rand.Int63(), 16))
	if err := os.MkdirAll(abs, 0755); err != nil {
		return "", err
	}
	return abs, nil
}

func Exists(path string) (bool, error) {
	if _, err := os.Stat(path); err != nil {
		if os.IsNotExist(err) {
			return false, nil
		}
		return false, err
	}
	return true, nil
}

func Exec(env Env, dir, prog string, argv ...string) error {
	cmd := exec.Command(prog, argv...)
	cmd.Dir = dir
	if env != nil {
		cmd.Env = env.Environ()
	}
	combined, err := cmd.CombinedOutput()
	if *flagShow && len(combined) > 0 {
		println(string(combined))
	}
	return err
}

// IsExitError returns true if err represents a process exit error
func IsExitError(err error) bool {
	_, ok := err.(*exec.ExitError)
	return ok
}

// Env holds environment variables
type Env map[string]string

func OSEnv() Env {
	environ := os.Environ()
	r := make(Env)
	for _, ev := range environ {
		kv := strings.SplitN(ev, "=", 2)
		if len(kv) != 2 {
			continue
		}
		r[kv[0]] = kv[1]
	}
	return r
}

func (env Env) Environ() []string {
	var r []string
	for k, v := range env {
		r = append(r, k + "=" + v)
	}
	return r
}

func (env Env) Get(key string) string {
	return env[key]
}

func (env Env) Set(key, value string) {
	env[key] = value
}

func (env Env) Copy() Env {
	r := make(Env)
	for k, v := range env {
		r[k] = v
	}
	return r
}

func CopyFile(srcName, dstName string) (written int64, err error) {
	src, err := os.Open(srcName)
	if err != nil {
		return
	}
	defer src.Close()

	dst, err := os.Create(dstName)
	if err != nil {
		return
	}
	defer dst.Close()

	return io.Copy(dst, src)
}

func ShellCopyFile(src, dst string) error {
	cmd := exec.Command("sh")
	cmd.Stdin = bytes.NewBufferString(fmt.Sprintf("cp %s %s\n", src, dst))
	combined, err := cmd.CombinedOutput()
	if *flagShow {
		println(string(combined))
	}
	return err
}
