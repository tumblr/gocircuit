package main

import (
	"bytes"
	"fmt"
	"os"
	"text/template"
	"tumblr/circuit/kit/posix"
	"tumblr/circuit/kit/sched/limiter"
	"tumblr/circuit/load/config"
)

const limitParallelTasks = 20

func Install(i *config.InstallConfig, b *config.BuildConfig, hosts []string) {
	l := limiter.New(limitParallelTasks)
	for _, host_ := range hosts {
		host := host_
		l.Go(func() {
			fmt.Printf("Installing on %s\n", host)
			if err := installHost(i, b, host); err != nil {
				fmt.Fprintf(os.Stderr, "Issue on %s: %s\n", host, err)
			}
		})
	}
	l.Wait()
}

const installShSrc = `mkdir -p {{.BinDir}} && mkdir -p {{.JailDir}} && mkdir -p {{.VarDir}}`

func installHost(i *config.InstallConfig, b *config.BuildConfig, host string) error {

	// Prepare shell script
	t := template.New("_")
	template.Must(t.Parse(installShSrc))
	var w bytes.Buffer
	if err := t.Execute(&w, &struct{BinDir, JailDir, VarDir string}{
		BinDir:  i.BinDir(),
		JailDir: i.JailDir(),
		VarDir:  i.VarDir(),
	}); err != nil {
		return err
	}
	install_sh := string(w.Bytes())

	// Execute remotely
	if _, _, err := posix.RemoteShell(host, install_sh); err != nil {
		return err
	}
	if err := posix.UploadDir(host, b.ShipDir, i.BinDir()); err != nil {
		return err
	}
	return nil
}
