// 4build automates the process of building a circuit application locally
package main

import (
	"os"
	"path"
	"strings"
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
	workerPkg string
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
		x.env.Set("PATH", flags.PrefixPath+":"+x.env.Get("PATH"))
	}
	//println(fmt.Sprintf("%#v\n", x.env))
	x.jail = flags.Jail
	x.workerPkg = flags.WorkerPkg
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
	println("--Packaging", x.binary)
	binpkg := workerPkgPath()
	_, workerName := path.Split(binpkg)
	shipFile := path.Join(tmpdir, x.binary)	// Destination binary location and name
	if _, err = CopyFile(path.Join(binpkg, workerName), shipFile); err != nil {
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

	// Place the zookeeper dynamic libraries in the shipment
	// Shipping Zookeeper is not necessary when static linking (currently enabled).
	/*
	println("--Packaging Zookeeper libraries")
	if err = ShellCopyFile(path.Join(x.zlib, "libzookeeper*"), tmpdir+"/"); err != nil {
		Fatalf("Problem copying Zookeeper library files (%s)\n", err)
	}
	*/

	return tmpdir
}

// workerPkgPath returns the absolute path to the app package that should be compiled as a circuit worker binary
func workerPkgPath() string {
	return path.Join(x.goPath["app"], "src", x.workerPkg)
}

func helperPkgPath(helper string) string {
	return path.Join(x.goPath["circuit"], "src/circuit/cmd", helper)
}

func buildCircuit() {

	// Prepare cgo environment for Zookeeper
	// TODO: Add Zookeeper build step. Don't rely on a prebuilt one.
	x.env.Set("CGO_CFLAGS", "-I" + x.zinclude)

	// Static linking (not available in Go1.0.3, available later, in +4ad21a3b23a4, for example)
	x.env.Set("CGO_LDFLAGS", path.Join(x.zlib, "libzookeeper_mt.a"))
	// Dynamic linking
	// x.env.Set("CGO_LDFLAGS", x.zlib + " -lzookeeper_mt"))

	// Cleanup set CGO_* flags at end
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
		if err := Shell(x.env, path.Join(x.goPath["circuit"], "src/circuit/cmd", cpkg), x.goCmd+" build -a"); err != nil {
			Fatalf("Problem compiling %s (%s)\n", cpkg, err)
		}
	}

	// Create a package for the runtime executable
	binpkg := workerPkgPath()

	// Build circuit runtime binary
	println("--Building", x.binary)
	// TODO: The -a flag here seems necessary. Otherwise changes in
	// circuit/sys do not seem to be reflected in recompiled tutorials when
	// the synchronization method for all repositories is rsync.
	// Understand what is going on. The flag should not be needed as the
	// circuit should see the changes in the sources inside the build jail.
	// Is this a file timestamp problem introduced by rsync?
	if err := Shell(x.env, binpkg, x.goCmd + " build -a"); err != nil {
		Fatalf("Problem with ‘(working directory %s) %s build’ (%s)\n", binpkg, x.goCmd, err)
	}
}

func buildGoCompiler(rebuild bool) {
	// Unset lingering CGO_* flags as they mess with the build of the Go compiler
	x.env.Unset("CGO_CFLAGS")
	x.env.Unset("CGO_LDFLAGS")

	// Check whether compiler subdirectory directory exists,
	// $jail/go
	ok, err := Exists(path.Join(x.jail, "/go"))
	if err != nil {
		Fatalf("Problem stat'ing %s (%s)", path.Join(x.jail, "/go"), err)
	}
	//Exec(x.env, x.jail, "which", "hg")
	if !ok {
		// If not, fetch the source tree
		if err = Shell(x.env, x.jail, "hg clone -u tip https://code.google.com/p/go"); err != nil {
			Fatalf("Problem cloning Go repository (%s)", err)
		}
		// Force rebuild
		rebuild = true
	} else {
		if rebuild {
			// Pull changes
			if err = Shell(x.env, path.Join(x.jail, "/go"), "hg pull"); err != nil {
				Fatalf("Problem pulling Go repository changes (%s)", err)
			}
			// Update working copy
			if err = Shell(x.env, path.Join(x.jail, "/go"), "hg update"); err != nil {
				Fatalf("Problem updating Go repository changes (%s)", err)
			}
		}
	}
	if rebuild {
		// Build Go compiler
		if err = Shell(x.env, path.Join(x.jail, "/go/src"), path.Join(x.jail, "/go/src/all.bash")); err != nil {
			if !IsExitError(err) {
				Fatalf("Problem building Go (%s)", err)
			}
		}
	}

	// Create build environment for building with this compiler
	x.goRoot = path.Join(x.jail, "/go")
	x.goBin = path.Join(x.goRoot, "/bin")
	x.goCmd = path.Join(x.goBin, "go")
	x.env.Set("PATH", x.goBin+":"+x.env.Get("PATH"))
	x.env.Set("GOROOT", x.goRoot)
}
