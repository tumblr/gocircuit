// 4crossbuild builds the circuit on a remote host (usually with a different OS
// and/or architecture) and ships the result locally.
package main

import (
	"bytes"
	"io/ioutil"
	"path"
	"os"
	"strings"
	"text/template"
)

/*
	Build jail layout:
		/flags
		/go
		/app/src
		/circuit/src
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

// Command-line tools to be built
var cmdPkg = []string{"4clear-helper"}

func main() {
	flags, flagsChanged := LoadFlags()

	// Initialize build environment
	x.binary = flags.Binary
	if strings.TrimSpace(x.binary) == "" {
		println("Missing name of target binary")
		os.Exit(1)
	}
	x.env = OSEnv()
	if flags.PrefixPath != "" {
		x.env.Set("PATH", flags.PrefixPath + ":" + x.env.Get("PATH"))
	}
	x.jail = flags.Jail
	x.appPkgs = []string{flags.Pkg}
	x.zinclude = flags.ZInclude
	x.zlib = flags.ZLib
	x.goPath = make(map[string]string)

	// Make jail if not present
	var err error
	if err = os.MkdirAll(x.jail, 0700); err != nil {
		Fatalf("Problem creating build jail (%s)\n", err)
	}

	Errorf("Building Go compiler\n")
	buildGoCompiler(flags.RebuildGo)

	Errorf("Updating circuit repository\n")
	// If repo name or fetch method has changed, remove any pre-existing clone
	fetchRepo("circuit", flags.CircuitRepo, flags.CircuitPath, flagsChanged.CircuitRepo)

	Errorf("Updating app repository\n")
	fetchRepo("app", flags.AppRepo, flags.AppPath, flagsChanged.AppRepo)

	Errorf("Building circuit binaries\n")
	buildCircuit()

	Errorf("Shipping install package\n")
	bundleDir := shipCircuit()
	Errorf("Build successful!\n")

	// Print temporary directory containing bundle
	Printf("%s\n", bundleDir)

	SaveFlags(flags)
}

func shipCircuit() string {
	tmpdir, err := MakeTempDir()
	if err != nil {
		Fatalf("Problem making packaging directory (%s)\n", err)
	}

	// Copy worker binary over to shipping directory
	println("+Packaging", x.binary)
	binpkg := workerPkgPath()
	shipFile := path.Join(tmpdir, x.binary)
	if _, err = CopyFile(path.Join(binpkg, x.binary), shipFile); err != nil {
		Fatalf("Problem copying circuit worker binary (%s)\n", err)
	}
	if err = os.Chmod(shipFile, 0755); err != nil {
		Fatalf("Problem chmod'ing circuit worker binary (%s)\n", err)
	}

	// Copy command-line helper tools over to shipping directory
	for _, cpkg := range cmdPkg {
		shipHelper := path.Join(tmpdir, cpkg)
		if _, err = CopyFile(path.Join(helperPkgPath(cpkg), cpkg), shipHelper); err != nil {
			Fatalf("Problem copying circuit helper binary (%s)\n", err)
		}
		if err = os.Chmod(shipHelper, 0755); err != nil {
			Fatalf("Problem chmod'ing circuit helper binary (%s)\n", err)
		}
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

func workerPkgPath() string {
	return path.Join(x.goPath["circuit"], "src", "autopkg", x.binary)	
}

func helperPkgPath(helper string) string {
	return path.Join(x.goPath["circuit"], "src/circuit/cmd", helper)
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

	// Re-build command-line tools
	for _, cpkg := range cmdPkg {
		if err := Exec(x.env, path.Join(x.goPath["circuit"], "src/circuit/cmd", cpkg) , x.goCmd, "build"); err != nil {
			Fatalf("Problem compiling %s (%s)\n", cpkg, err)
		}
	}

	// Create a package for the runtime executable
	binpkg := workerPkgPath()
	if err := os.RemoveAll(binpkg); err != nil {
		Fatalf("Problem removing old autopkg directory %s (%s)\n", binpkg, err)
	}
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

func repoSchema(s string) (schema, url string) {
	switch {
	case strings.HasPrefix(s, "{git}"):
		return "git", s[len("{git}"):]
	case strings.HasPrefix(s, "{rsync}"):
		return "rsync", s[len("{rsync}"):]
	}
	Fatalf("Repo '%s' has unrecognizable schema\n", s)
	panic("unr")
}

func cloneGitRepo(repo, parent string) {
	// If not, clone the source tree
	if err := Exec(nil, parent, "git", "clone", repo); err != nil {
		Fatalf("Problem cloning repo '%s' (%s)", repo, err)
	}
}

func pullGitRepo(dir string) {
	if err := Exec(nil, dir, "git", "pull", "origin", "master"); err != nil {
		Fatalf("Problem pulling repo in %s (%s)", dir, err)
	}
}

func rsyncRepo(src, dstparent string) {
	if err := Exec(nil, "", "rsync", "-acrv", "--exclude", ".git", "--exclude", ".hg", "--exclude", "*.a", src, dstparent); err != nil {
		Fatalf("Problem rsyncing dir '%s' to within '%s' (%s)", src, dstparent, err)
	}
}

func fetchRepo(namespace, repo, gopath string, fetchFresh bool) {

	schema, repo := repoSchema(repo)

	// If fetching fresh, remove pre-existing clones
	if fetchFresh {
		if err := os.RemoveAll(path.Join(x.jail, namespace)); err != nil {
			Fatalf("Problem removing old repo clone (%s)\n", err)
		}
	}

	// Make _build/namespace/src
	repoSrc := path.Join(x.jail, namespace, "src")
	if err := os.MkdirAll(repoSrc, 0700); err != nil {
		Fatalf("Problem creating app source path %s (%s)\n", repoSrc, err)
	}
	repoPath := path.Join(repoSrc, repoName(repo))

	// Check whether repo directory exists
	ok, err := Exists(repoPath)
	if err != nil {
		Fatalf("Problem stat'ing %s (%s)", repoPath, err)
	}
	switch schema {
	case "git":
		if !ok {
			cloneGitRepo(repo, repoSrc)
		} else {
			pullGitRepo(repoPath)
		}
	case "rsync":
		rsyncRepo(repo, repoSrc)
	default:
		Fatalf("Unrecognized repo schema: %s\n", schema)
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
	// Oddly, On Linux the dynamic linker seems to need Zookeeper libs even to compile Go,
	// likely in effect of an existing CGO_LDFLAGS.
	x.env.Set("LD_LIBRARY_PATH", x.zlib)
	// So, to be safe, we provision this for OSX as well.
	x.env.Set("DYLD_LIBRARY_PATH", x.zlib)

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
