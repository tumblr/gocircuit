// 4crossbuild builds the circuit on a remote host (usually with a different OS
// and/or architecture) and ships the result locally.
package main

import (
	"bytes"
	"flag"
	"io/ioutil"
	"path"
	"os"
	"strings"
	"text/template"
)

var (
	flagBinary           = flag.String("binary", "4r", "Preferred name for the resulting runtime binary")
	flagJail             = flag.String("jail", path.Join(os.Getenv("HOME"), "_circuit/build"), "Build jail directory")

	flagAppRepo          = flag.String("app", "", "App repository")
	flagAppPath          = flag.String("appsrc", "", "GOPATH relative to app repository")

	flagPkg              = flag.String("pkg", "", "Package to import for side-effects in circuit runtime binary")

	flagShow             = flag.Bool("show", false, "Show output of underlying build commands")
	flagRebuildGo        = flag.Bool("rebuildgo", false, "Force fetch and rebuild of the Go compiler")

	flagZookeeperInclude = flag.String("zinclude", path.Join(os.Getenv("HOME"), "local/include/c-client-src") , "Zookeeper C headers directory")
	flagZookeeperLib     = flag.String("zlib", path.Join(os.Getenv("HOME"), "local/lib") , "Zookeeper libraries directory")

	flagCircuitRepo      = flag.String("cir", "git@github.com:tumblr/gocircuit.git", "Circuit repository")
	flagCircuitPath      = flag.String("cirsrc", ".", "GOPATH relative to circuit repository")
)

/*
	Build jail layout:
		_build/go
		_build/app/src
		_build/circuit/src
*/

var x struct {
	env       Env
	jail      string
	appPkgs   []string
	binary    string
	zinclude  string
	zlib      string
	goRoot    string
	goBin     string
	goCmd     string
	goPath    map[string]string
}

func main() {
	flag.Parse()

	// Initialize build environment
	x.binary = *flagBinary
	x.env = OSEnv()
	x.jail = *flagJail
	x.appPkgs = []string{*flagPkg}
	x.zinclude = *flagZookeeperInclude
	x.zlib = *flagZookeeperLib
	x.goPath = make(map[string]string)

	// Make jail if not present
	var err error
	if err = os.MkdirAll(x.jail, 0700); err != nil {
		Fatalf("Problem creating build jail (%s)\n", err)
	}

	Errorf("Building Go compiler\n")
	buildGoCompiler(*flagRebuildGo)

	Errorf("Updating circuit repository\n")
	fetchRepo("circuit", *flagCircuitRepo, *flagCircuitPath)

	Errorf("Updating app repository\n")
	fetchRepo("app", *flagAppRepo, *flagAppPath)

	Errorf("Building circuit binaries\n")
	buildCircuit()

	Errorf("Shipping install package\n")
	bundleDir := shipCircuit()
	Errorf("Build successful!\n")

	// Print temporary directory containing bundle
	Printf("%s\n", bundleDir)
}

func shipCircuit() string {
	tmpdir, err := MakeTempDir()
	if err != nil {
		Fatalf("Problem making packaging directory (%s)\n", err)
	}

	// Copy binaries over to shipping directory
	println("+Packaging", x.binary)
	binpkg := binPkg()
	shipFile := path.Join(tmpdir, x.binary)
	if _, err = CopyFile(path.Join(binpkg, x.binary), shipFile); err != nil {
		Fatalf("Problem copying circuit binary (%s)\n", err)
	}
	if err = os.Chmod(shipFile, 0755); err != nil {
		Fatalf("Problem chmod'ing circuit binary (%s)\n", err)
	}

	// zookeeper lib
	println("+Packaging Zookeeper libraries")
	if err = ShellCopyFile(path.Join(x.zlib, "libzookeeper*"), tmpdir + "/"); err != nil {
		Fatalf("Problem copying Zookeeper library files (%s)\n", err)
	}

	return tmpdir
}

// Source code of a circuit runtime executable
const mainSrc = `
package main
import (
	_ "circuit/load"
	{{range .}}
	_ "{{.}}"
	{{end}}
)
func main() {}
`

func binPkg() string {
	return path.Join(x.goPath["circuit"], "src", "autopkg", x.binary)	
}

func buildCircuit() {

	// Prepare cgo environment for Zookeeper
	// TODO: Add Zookeeper build step. Don't rely on a prebuilt one.
	x.env.Set("CGO_CFLAGS", "-I" + x.zinclude)
	x.env.Set("CGO_LDFLAGS", "-L" + x.zlib + " -lm -lpthread -lzookeeper_mt")
	defer x.env.Unset("CGO_CFLAGS")
	defer x.env.Unset("CGO_LDFLAGS")

	// Remove any installed packages
	if err := os.RemoveAll(path.Join(x.goPath["circuit"], "pkg")); err != nil {
		Fatalf("Problem removing circuit pkg directory (%s)\n", err)
	}
	if err := os.RemoveAll(path.Join(x.goPath["app"], "pkg")); err != nil {
		Fatalf("Problem removing app pkg directory (%s)\n", err)
	}

	// Re-build application package
	/*
	for _, apkg := range x.appPkgs {
		if err := Exec(x.env, path.Join(x.goPath["app"], "src", apkg) , x.goCmd, "install"); err != nil {
			Fatalf("Problem compiling main.go (%s)\n", err)
		}
	}
	*/

	// Create a package for the runtime executable
	binpkg := binPkg()
	if err := os.MkdirAll(binpkg, 0700); err != nil {
		Fatalf("Problem creating runtime package %s (%s)\n", binpkg, err)
	}
	
	// Write main.go
	t := template.New("main")
	template.Must(t.Parse(mainSrc))
	var w bytes.Buffer
	if err := t.Execute(&w, x.appPkgs); err != nil {
		Fatalf("Problem preparing main.go (%s)\n", err)
	}
	if err := ioutil.WriteFile(path.Join(binpkg, "main.go"), w.Bytes(), 0664); err != nil {
		Fatalf("Problem writing main.go (%s)\n", err)
	}

	// Build circuit runtime binary
	println("+Building", x.binary)
	if err := Exec(x.env, binpkg, x.goCmd, "build"); err != nil {
		Fatalf("Problem compiling main.go (%s)\n", err)
	}
}

// repoName returns the top-level name of a GIT repository from its URL
// E.g. git@github.com:tumblr/cirapp.git -> cirapp
func repoName(repo string) string {
	if strings.HasSuffix(repo, ".git") {
		repo = repo[:len(repo) - len(".git")]
	}
__For:
	for i := len(repo) - 1; i >= 0; i-- {
		switch repo[i] {
		case ':', '/', '@':
			repo = repo[i + 1:]
			break __For
		}
	}
	return repo
}

/*
	_build/namespace/src/cloned_user_repo/a/src/a/b/c
	[-------------------]					repoSrc
	[------------------------------------]			repoPath
	[===============]					GOPATH, if gopath == ""
	[======================================]		GOPATH, if gopath == "/a"
*/
func fetchRepo(namespace, repo, gopath string) {

	// Make _build/namespace/src
	repoSrc := path.Join(x.jail, path.Join(namespace, "src"))
	if err := os.MkdirAll(repoSrc, 0700); err != nil {
		Fatalf("Problem creating app source path %s (%s)\n", repoSrc, err)
	}
	repoPath := path.Join(repoSrc, repoName(repo))

	// Check whether repo directory exists
	ok, err := Exists(repoPath)
	if err != nil {
		Fatalf("Problem stat'ing %s (%s)", repoPath, err)
	}
	if !ok {
		// If not, clone the source tree
		if err = Exec(nil, repoSrc, "git", "clone", repo); err != nil {
			Fatalf("Problem cloning repo %s (%s)", repo, err)
		}
	} else {
		// Pull changes
		if err = Exec(nil, repoPath, "git", "pull", "origin", "master"); err != nil {
			Fatalf("Problem pulling repo in %s (%s)", repoPath, err)
		}
	}

	// Create build environment for building in this repo
	oldGoPath := x.env.Get("GOPATH")
	var p string
	if gopath == "" {
		p = path.Join(x.jail, namespace)
	} else {
		p = path.Join(repoPath, gopath)
	}
	x.env.Set("GOPATH", p + ":" + oldGoPath)
	x.goPath[namespace] = p
}

func buildGoCompiler(rebuild bool) {
	// Check whether compiler subdirectory directory exists,
	// $jail/go
	ok, err := Exists(path.Join(x.jail, "/go"))
	if err != nil {
		Fatalf("Problem stat'ing %s (%s)", path.Join(x.jail, "/go"), err)
	}
	//Errorf("EEE %#v\n", x.env.Environ())
	if !ok {
		// If not, fetch the source tree
		if err = Exec(x.env, x.jail, "hg", "clone", "-u", "tip", "https://code.google.com/p/go"); err != nil {
			Fatalf("Problem cloning Go repository (%s)", err)
		}
		// Force rebuild
		rebuild = true
	} else {
		if rebuild {
			// Pull changes
			if err = Exec(x.env, path.Join(x.jail, "/go"), "hg", "pull"); err != nil {
				Fatalf("Problem pulling Go repository changes (%s)", err)
			}
			// Update working copy
			if err = Exec(x.env, path.Join(x.jail, "/go"), "hg", "update"); err != nil {
				Fatalf("Problem updating Go repository changes (%s)", err)
			}
		}
	}
	if rebuild {
		// Build Go compiler
		if err = Exec(x.env, path.Join(x.jail, "/go/src"), path.Join(x.jail, "/go/src/all.bash")); err != nil {
			if !IsExitError(err) {
				Fatalf("Problem building Go (%s)", err)
			}
		}
	}

	// Create build environment for building with this compiler
	x.goRoot = path.Join(x.jail, "/go")
	x.goBin = path.Join(x.goRoot, "/bin")
	x.goCmd = path.Join(x.goBin, "go")
	x.env.Set("PATH", x.goBin + ":" + x.env.Get("PATH"))
	x.env.Set("GOROOT", x.goRoot)
}
