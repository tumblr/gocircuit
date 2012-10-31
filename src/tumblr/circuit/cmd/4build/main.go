// 4crossbuild builds the circuit on a remote host (usually with a different OS
// and/or architecture) and ships the result locally.
package main

import (
	"flag"
	"path"
	"os"
	"strings"
)

// XXX: Use var for GOPATH, since GOPATH needs to hold both circuit and user roots

var (
	flagJail             = flag.String("jail", path.Join(os.Getenv("HOME"), "_circuit/_build"), "Build jail directory")
	// If flagRepo is empty, we expect that the repo subdirectory is
	// already in place inside the build jail.  This case is useful for
	// cases when the origin repo is local and the user can simply put a
	// link to it in the build jail
	flagCircuit          = flag.String("circuit", "git@github.com:tumblr/gocircuit.git", "Circuit repository")
	flagRepo             = flag.String("repo", "", "User repository")
	flagGoPath           = flag.String("gopath", "", "GOPATH relative to repository root")
	flagPkg              = flag.String("pkg", "", "Package to build")
	flagShow             = flag.Bool("show", false, "Show output of executed commands")
	flagBuildGo          = flag.Bool("buildgo", false, "Fetch and rebuild the Go compiler")
	flagZookeeperInclude = flag.String("zinclude", path.Join(os.Getenv("HOME"), "local/include/c-client-src") , "Zookeeper C headers directory")
	flagZookeeperLib     = flag.String("zlib", path.Join(os.Getenv("HOME"), "local/lib") , "Zookeeper libraries directory")
)

/*
	Build jail layout:
		_circuit/_build/go	= Go compiler cloned here
		_circuit/_build/src	= GOPATH, user repos are cloned in here
*/

var BuildPkgs []string

func main() {
	flag.Parse()
	BuildPkgs = []string{*flagPkg}

	var err error
	if err = os.MkdirAll(*flagJail, 0700); err != nil {
		Fatalf("Problem creating build jail (%s)\n", err)
	}
	osenv := OSEnv()
	env := make(Env)
	env.Set("PATH", osenv.Get("PATH"))

	Errorf("Building Go compiler\n")
	if env, err = BuildGoCompiler(env, *flagBuildGo); err != nil {
		Fatalf("Error building Go compiler (%s)\n", err)
	}
	Errorf("Updating user repository\n")
	if env, err = FetchRepo(env, *flagRepo); err != nil {
		Fatalf("Error fetching repository (%s)\n", err)
	}
	Errorf("Building binaries\n")
	if err = BuildCircuit(env); err != nil {
		Fatalf("Error building circuit (%s)\n", err)
	}
	Errorf("Shipping install package\n")
	bundleDir, err := ShipCircuit(env)
	if err != nil {
		Fatalf("Error shipping package (%s)\n", err)
	}
	Errorf("Build successful!\n")

	// Print temporary directory containing bundle
	Printf("%s\n", bundleDir)
}

func ShipCircuit(env Env) (string, error) {
	tmpdir, err := MakeTempDir()
	if err != nil {
		return "", err
	}

	// Copy binaries over to shipping directory
	for _, pkg := range BuildPkgs {
		println("Packaging", pkg)
		pkgpath := pkgPath(env, pkg)
		_, name := path.Split(pkg)
		shipFile := path.Join(tmpdir, name)
		if _, err = CopyFile(path.Join(pkgpath, name), shipFile); err != nil {
			return "", err
		}
		if err = os.Chmod(shipFile, 0755); err != nil {
			return "", err
		}
	}

	// zookeeper lib
	println("Packaging Zookeeper libraries")
	if err = ShellCopyFile(path.Join(*flagZookeeperLib, "libzookeeper*"), tmpdir + "/"); err != nil {
		return "", err
	}

	return tmpdir, nil
}

func pkgPath(env Env, pkg string) string {
	return path.Join(env.Get("GOPATH"), "src", pkg)
}

func BuildCircuit(env Env) error {
	goBinary := env.Get("GOBINARY")

	// Zookeeper
	// TODO: Add Zookeeper build step. Don't rely on a prebuilt one.
	env.Set("CGO_CFLAGS", "-I" + *flagZookeeperInclude)
	env.Set("CGO_LDFLAGS", "-L" + *flagZookeeperLib + " -lm -lpthread -lzookeeper_mt")

	for _, pkg := range BuildPkgs {
		println("+Building", pkg)
		if err := Exec(env, pkgPath(env, pkg), goBinary, "build"); err != nil {
			return err
		}
	}
	return nil
}

// E.g. git@github.com:tumblr/cirapp.git -> cirapp
func repoDir(repo string) string {
	if strings.HasSuffix(repo, ".git") {
		repo = repo[:len(repo)-len(".git")]
	}
	for i := len(repo)-1; i >= 0; i++ {
		switch repo[i] {
		case ':', '/', '@':
			repo = repo[i+1:]
			break
		}
	}
	return repo
}

/*
	_build_jail/src/cloned_user_repo/a/src/a/b/c
	[==========]					GOPATH, if *flagGoPath == ""
	[===============================|=]		GOPATH, if *flagGoPath == "/a"

	[--------------]
	               [----------------]
	[-------------------------------]		repoPath
*/
func FetchRepo(env Env, repo string) (Env, error) {

	// Make _build_jail/src
	gosrc := path.Join(*flagJail, "src")
	if err := os.MkdirAll(gosrc, 0700); err != nil {
		Fatalf("Problem creating %s (%s)\n", gosrc, err)
	}
	repoPath := path.Join(gosrc, repoDir(repo))

	// Check whether repo directory exists
	ok, err := Exists(repoPath)
	if err != nil {
		return nil, err
	}
	if !ok {
		// If not, clone the source tree
		if err = Exec(nil, gosrc, "git", "clone", repo); err != nil {
			return nil, err
		}
	} else {
		if repo != "" {
			// Pull changes
			if err = Exec(nil, repoPath, "git", "pull", "origin", "master"); err != nil {
				return nil, err
			}
		}
	}
	// Create build environment for building in this repo
	env = env.Copy()
	if *flagGoPath == "" {
		env.Set("GOPATH", *flagJail)
	} else {
		env.Set("GOPATH", path.Join(repoPath, *flagGoPath))
	}
	return env, nil
}

func BuildGoCompiler(env Env, rebuild bool) (Env, error) {
	// Check whether compiler directory exists
	ok, err := Exists(path.Join(*flagJail, "/go"))
	if err != nil {
		return nil, err
	}
	if !ok {
		// If not, fetch the source tree
		if err = Exec(nil, *flagJail, "hg", "clone", "-u", "tip", "https://code.google.com/p/go"); err != nil {
			return nil, err
		}
		// Force rebuild
		rebuild = true
	} else {
		if rebuild {
			// Pull changes
			if err = Exec(nil, path.Join(*flagJail, "/go"), "hg", "pull"); err != nil {
				return nil, err
			}
			// Update working copy
			if err = Exec(nil, path.Join(*flagJail, "/go"), "hg", "update"); err != nil {
				return nil, err
			}
		}
	}
	if rebuild {
		// Build Go compiler
		if err = Exec(env, path.Join(*flagJail, "/go/src"), path.Join(*flagJail, "/go/src/all.bash")); err != nil {
			if !IsExitError(err) {
				return nil, err
			}
		}
	}
	// Create build environment for building with this compiler
	env = env.Copy()
	goroot := path.Join(*flagJail, "/go")
	gobin := path.Join(goroot, "/bin")
	env.Set("PATH", gobin + ":" + env.Get("PATH"))
	env.Set("GOROOT", goroot)
	env.Set("GOBINARY", path.Join(gobin, "go"))

	return env, nil
}
