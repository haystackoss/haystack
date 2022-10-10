package factory

import (
	"github.com/nabaz-io/nabaz/pkg/testrunner/scm/history/git"
	"github.com/nabaz-io/nabaz/pkg/testrunner/scm/history/git/local"
)

func NewGitHistory(repoPath string) (git.GitHistory, error) {
	repo, err := local.NewLocalGitHistory(repoPath)
	if err != nil {
		return nil, err
	}
	return repo, nil
}
