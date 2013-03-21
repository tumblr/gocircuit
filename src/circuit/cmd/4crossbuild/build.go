package main

import (
	"bytes"
	"circuit/kit/posix"
	"circuit/load/config"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"strings"
	"text/template"
)

const build_sh_src = `{{.Tool}} ` +
	`'-binary={{.Binary}}' '-jail={{.Jail}}' ` +
	`'-app={{.AppRepo}}' '-appsrc={{.AppSrc}}' ` +
	`'-pkg={{.Pkg}}' '-show={{.Show}}' '-rebuildgo={{.RebuildGo}}' ` +
	`'-zinclude={{.ZookeeperInclude}}' '-zlib={{.ZookeeperLib}}' ` +
	`'-cir={{.CircuitRepo}}' '-cirsrc={{.CircuitSrc}}' '-prefixpath={{.PrefixPath}}' `

func Build(cfg *config.BuildConfig) error {
	// Prepare sh script
	t := template.New("_")
	template.Must(t.Parse(build_sh_src))
	var w bytes.Buffer
	if err := t.Execute(&w, cfg); err != nil {
		panic("parse cross-build script")
	}
	build_sh := string(w.Bytes())

	if cfg.Show {
		println(build_sh)
	}

	// Execute remotely
	cmd := exec.Command("ssh", cfg.Host, "sh")
	cmd.Stdin = bytes.NewBufferString(build_sh)

	// Capture stdout and stderr
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return err
	}
	stderr, err := cmd.StderrPipe()
	if err != nil {
		return err
	}
	prefix := fmt.Sprintf("%s/4build|", cfg.Host)
	posix.ForwardStderr(prefix, stderr)

	if err = cmd.Start(); err != nil {
		return err
	}

	// Read result (remote directory of built bundle) from stdout
	result, _ := ioutil.ReadAll(stdout)
	if err = cmd.Wait(); err != nil {
		return err
	}

	// Fetch the built shipping bundle
	if err = os.MkdirAll(cfg.ShipDir, 0700); err != nil {
		return err
	}

	// Make ship directory if not present
	if err := os.MkdirAll(cfg.ShipDir, 0755); err != nil {
		return err
	}

	// Clean the ship directory
	if _, _, err = posix.Shell(`rm -f ` + cfg.ShipDir + `/*`); err != nil {
		return err
	}

	// Cleanup remote dir of built files
	r := strings.TrimSpace(string(result))
	if r == "" {
		return errors.New("empty shipping source directory")
	}

	// Download files
	println("Downloading from", r)
	if err = posix.DownloadDir(cfg.Host, r, cfg.ShipDir); err != nil {
		return err
	}
	println("Download successful.")
	return nil
}
