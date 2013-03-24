package main


func parseRepo(s string) (schema, key, value, url string) {
	switch {
	case strings.HasPrefix(s, "{hg}"):
		schema, s = "hg", s[len("{hg}"):]
	case strings.HasPrefix(s, "{git}"):
		schema, s = "git", s[len("{git}"):]
	case strings.HasPrefix(s, "{rsync}"):
		schema, s = "rsync", s[len("{rsync}"):]
	default:
		Fatalf("Repo '%s' has unrecognizable schema\n", s)
	}
	i := strings.Index(s, "}")
	if len(s) > 1 && s[0] = '{' && i > 0 {
		var arg string
		arg, s = s[1:i], s[i+1:]
		key, value = parseArg(arg)
	}
	url = s
	return
}

//
//	{hg}{changeset:4ad21a3b23a4}
//	{hg}{id:4ad21a3b23a4}
//	{hg}{rev:3452}
//	{hg}{tag:weekly}
//	{hg}{tip}
//	{hg}{branch:master}
//
//	{git}{rev:51e592253000600d586408f3e36a3f4692011086}
//
func parseArg(arg string) (key, value string) {
	part := strings.SplitN(arg, ":", 2)
	if len(part) > 0 {
		key = part[0]
	}
	if len(part) > 1 {
		value = part[1]
	}
	return
}

func cloneMercurialRepo(repo, arg, parent string) {
	// If not, clone the source tree
	if err := Shell(x.env, parent, "hg clone " + arg + " " + repo); err != nil {
		Fatalf("Problem cloning repo '%s' (%s)", repo, err)
	}
}

func cloneGitRepo(repo, arg, parent string) {
	// If not, clone the source tree
	if err := Shell(x.env, parent, "git clone " + repo); err != nil {
		Fatalf("Problem cloning repo '%s' (%s)", repo, err)
	}
}

func pullGitRepo(dir string) {
	if err := Shell(x.env, dir, "git pull origin master"); err != nil {
		Fatalf("Problem pulling repo in %s (%s)", dir, err)
	}
}

func rsyncRepo(src, dstparent string) {
	if err := Shell(x.env, "", "rsync -acrv --delete --exclude .git --exclude .hg --exclude *.a "+src+" "+dstparent); err != nil {
		Fatalf("Problem rsyncing dir '%s' to within '%s' (%s)", src, dstparent, err)
	}
}

func fetchRepo(namespace, repo, gopath string, fetchFresh bool) {

	schema, key, value, url := parseRepo(repo)

	// If fetching fresh, remove pre-existing clones
	if fetchFresh {
		if err := os.RemoveAll(path.Join(x.jail, namespace)); err != nil {
			Fatalf("Problem removing old repo clone (%s)\n", err)
		}
	}

	// Make jail/namespace/src
	repoSrc := path.Join(x.jail, namespace, "src")
	if err := os.MkdirAll(repoSrc, 0700); err != nil {
		Fatalf("Problem creating app source path %s (%s)\n", repoSrc, err)
	}
	??
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
	x.env.Set("GOPATH", p+":"+oldGoPath)
	x.goPath[namespace] = p
}

