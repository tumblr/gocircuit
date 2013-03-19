package main

import (
	"encoding/json"
	"flag"
	"io/ioutil"
	"path"
	"os"
	"strings"
)

var (
	flagBinary      = flag.String("binary", "4r", "Preferred name for the resulting runtime binary")
	flagJail        = flag.String("jail", "", "Build jail directory")

	flagAppRepo     = flag.String("app", "", "App repository")
	flagAppPath     = flag.String("appsrc", "", "GOPATH relative to app repository")

	flagPkg         = flag.String("pkg", "", "Package to import for side-effects in circuit runtime binary")

	flagShow        = flag.Bool("show", false, "Show output of underlying build commands")
	flagRebuildGo   = flag.Bool("rebuildgo", false, "Force fetch and rebuild of the Go compiler")

	flagZInclude    = flag.String("zinclude", "", "Zookeeper C headers directory")
	flagZLib        = flag.String("zlib", "", "Zookeeper libraries directory")

	flagCircuitRepo = flag.String("cir", "git@github.com:tumblr/gocircuit.git", "Circuit repository")
	flagCircuitPath = flag.String("cirsrc", ".", "GOPATH relative to circuit repository")

	flagPrefixPath  = flag.String("prefixpath", "", "Prefix to add to default PATH environment")
)

// Flags is used to persist the state of command-line flags in the jail
type Flags struct {
	Binary      string
	Jail        string
	AppRepo     string
	AppPath     string
	Pkg         string
	Show        bool
	RebuildGo   bool
	ZInclude    string
	ZLib        string
	CircuitRepo string
	CircuitPath string
	PrefixPath  string
}

func (flags *Flags) FlagsFile() string {
	return path.Join(flags.Jail, "flags")
}

// FlagsChanged indicates which flag groups have changed since the previous
// invocation of the build tool
type FlagsChanged struct {
	Binary      bool
	Jail        bool
	AppRepo     bool
	Pkg         bool
	CircuitRepo bool
}

func getFlags() *Flags {
	return &Flags{
		Binary:      strings.TrimSpace(*flagBinary),
		Jail:        strings.TrimSpace(*flagJail),
		AppRepo:     strings.TrimSpace(*flagAppRepo),
		AppPath:     strings.TrimSpace(*flagAppPath),
		Pkg:         strings.TrimSpace(*flagPkg),
		Show:        *flagShow,
		RebuildGo:   *flagRebuildGo,
		ZInclude:    strings.TrimSpace(*flagZInclude),
		ZLib:        strings.TrimSpace(*flagZLib),
		CircuitRepo: strings.TrimSpace(*flagCircuitRepo),
		CircuitPath: strings.TrimSpace(*flagCircuitPath),
		PrefixPath:  strings.TrimSpace(*flagPrefixPath),
	}
}

func LoadFlags() (*Flags, *FlagsChanged) {
	flag.Parse()
	flags := getFlags()

	// Read old flags from jail
	oldFlags := &Flags{}
	hbuf, err := ioutil.ReadFile(flags.FlagsFile())
	if err != nil {
		println("No previous build flags found in jail.")
		goto __Diff
	}
	if err = json.Unmarshal(hbuf, oldFlags); err != nil {
		println("Previous flags cannot parse: ", err.Error())
		goto __Diff
	}

	// Compare old and new flags
__Diff:
	flagsChanged := &FlagsChanged{
		Binary:      flags.Binary != oldFlags.Binary,
		Jail:        flags.Jail != oldFlags.Jail,
		AppRepo:     flags.AppRepo != oldFlags.AppRepo || flags.AppPath != oldFlags.AppPath,
		Pkg:         flags.Pkg != oldFlags.Pkg,
		CircuitRepo: flags.CircuitRepo != oldFlags.CircuitRepo || flags.CircuitPath != oldFlags.CircuitPath,
	}

	return flags, flagsChanged
}

func SaveFlags(flags *Flags) {
	fbuf, err := json.Marshal(flags)
	if err != nil {
		println("Problems marshaling flags: ", err.Error())
		os.Exit(1)
	}
	if err = ioutil.WriteFile(flags.FlagsFile(), fbuf, 0600); err != nil {
		println("Problems writing flags: ", err.Error())
		os.Exit(1)
	}
}
