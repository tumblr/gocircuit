package main

import (
	"fmt"
	"os"
	"tumblr/util/posix"
	"tumblr/circuit/load"
	"tumblr/circuit/kit/sched/limiter"
)

const limitParallelTasks = 20

func Install(c *boot.InstallConfig, shipDir string, hosts []string) {
	l := limiter.New(limitParallelTasks)
	for _, host_ := range hosts {
		host := host_
		l.Go(func() {
			fmt.Printf("Installing on %s\n", host)
			if err := installHost(c, shipDir, host); err != nil {
				fmt.Fprintf(os.Stderr, "Issue on %s: %s\n", host, err)
			}
		})
	}
	l.Wait()
}

const install_sh_src = `
mkdir -p {{.BinDir}} && mkdir -p {{.JailDir}} && mkdir -p {{.VarDir}}
`

func installHost(c *boot.InstallConfig, shipDir string, host string) error {
	install_sh := posix.MustParseAndExecute(install_sh_src, &struct{BinDir, JailDir, VarDir string}{
		BinDir:  c.BinDir(),
		JailDir: c.JailDir(),
		VarDir:  c.VarDir(),
	})
	if _, _, err := posix.RunRemoteShell(host, install_sh); err != nil {
		return err
	}
	if err := posix.UploadDir(host, shipDir, c.BinDir()); err != nil {
		return err
	}
	return nil
}
