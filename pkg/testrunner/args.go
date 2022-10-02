package testrunner

type Arguements struct {
	Cmdline    string // pytest -v, go test ./..., etc
	StorageUrl string // mongodb://host:7190
	WebUrl     string //
	Pkgs       string // "pkg1 pkg2 pkg3 ./dir4/... ./dir5"
	Token      string // Token to access github account.
	RepoUrl    string // https://github.com/trovalds/linux
	CommitID   string // CommitID is the commit id of the change (in case there is no .git)
	Username   string
	Password   string
	RepoPath   string // Path to the repo, defaults to "."
}
