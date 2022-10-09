package factory

import (
	"log"

	"github.com/nabaz-io/nabaz/pkg/testrunner/scm/history/git"
	"github.com/nabaz-io/nabaz/pkg/testrunner/scm/history/git/local"
)

func NewGitHistory(repoPath string) (git.GitHistory, error) {
	repo, err := local.NewLocalGitHistory(repoPath, ".git")
	if err != nil {
		return nil, err
	}
	log.Println("Found .git, using it for code history...")
	return repo, nil
}
