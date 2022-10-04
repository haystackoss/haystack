package remote

import (
	"errors"

	"github.com/nabaz-io/nabaz/pkg/testrunner/scm/code"
	"github.com/nabaz-io/nabaz/pkg/testrunner/scm/history/git"
)

type GithubAPIHistory struct {
	// contains filtered or unexported fields
	token        string
	remoteURL    string
	headCommitID string
}

func NewGithubAPIHistory(token, remoteURL, commitID string) (git.GitHistory, error) {
	if commitID == "" || remoteURL == "" {
		return nil, errors.New("--commitid and --repo-url are all required for remote history to work")
	}
	// NewGithubAPI creates a new GithubAPI
	return &GithubAPIHistory{
		token:        token,
		remoteURL:    remoteURL,
		headCommitID: commitID,
	}, nil
}

func (g *GithubAPIHistory) GetCommitParents(commitID string) ([]string, error) {
	panic("implement me")
	// GetCommitParents returns the parents of a commit
	return nil, nil
}

func (g *GithubAPIHistory) GetFileFromCommit(filePath string, ref string) (string, error) {
	panic("implement me")
	// GetFileFromCommit returns the content of a file from a commit
	return "", nil
}

func (g *GithubAPIHistory) Diff(currentCommit string, olderCommit string) ([]code.FileDiff, error) {
	panic("implement me")
	// Diff returns the diff between two commits
	return nil, nil
}
