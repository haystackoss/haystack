package testrunner

type Arguements struct {
	Cmdline    string // pytest -v, go test ./..., etc
	Pkgs       string // "pkg1 pkg2 pkg3 ./dir4/... ./dir5"
	RepoPath   string // Path to the repo, defaults to "."
}
