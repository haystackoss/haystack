package factory

import (
	"log"

	"github.com/nabaz-io/nabaz/pkg/testrunner/scm/history/git"
	"github.com/nabaz-io/nabaz/pkg/testrunner/scm/history/git/local"
	"github.com/nabaz-io/nabaz/pkg/testrunner/scm/history/git/remote"
)

func NewGitHistory(repoPath, repoUrl, token, commitID string) (git.GitHistory, error) {
	// NewGitHistory creates a new GitHistory
	if repoUrl != "" {
		log.Println("Using GithubAPI for code history")
		history, err := remote.NewGithubAPIHistory(repoUrl, token, commitID)
		return history, err
	} else {
		repo, err := local.NewLocalGitHistory(repoPath)
		if err != nil {
			return nil, err
		}
		if repo != nil {
			log.Println("Found .git, using it for code history...")
			return repo, nil
		}
	}
}
