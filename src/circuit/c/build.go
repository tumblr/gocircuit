package c

// Build describes a Go compilation environment
type Build struct {
	WorkingGoPath string    // GOPATH of the user Go repo
	GoPaths       GoPaths   // All GOPATH paths
}

// NewWorkingBuild creates a new build environment, where the working
// gopath is derived from the current working directory.
func NewWorkingBuild() (*Build, error) {
	gopath, err := FindWorkingGoPath()
	if err != nil {
		return nil, err
	}
	return &Build{
		WorkingGoPath: gopath,
		GoPaths:       GetGoPaths(),
	}, nil
}
