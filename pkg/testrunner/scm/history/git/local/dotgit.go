package local

import (
	gitCli "github.com/go-git/go-git"
	"github.com/nabaz-io/nabaz/pkg/testrunner/scm/code"
)

// LocalGitHistory is history supplied by .git
type LocalGitHistory struct {
	// Path is the path to the local git repository.
	Path         string
	HeadCommitID string
	*gitCli.Repository
}

// NewLocalGitRepo creates a new LocalGitRepo.
func NewLocalGitRepo(path string) (*LocalGitHistory, error) {
	localRepo := &LocalGitHistory{
		Path: path,
	}
	git, err := gitCli.PlainOpen(path)
	if err != nil {
		return nil, err
	}

	head, err := git.Head()
	if err != nil {
		return nil, err
	}

	localRepo.HeadCommitID = head.String()
	localRepo.Repository = git

	return localRepo, nil
}

// GetHeadCommitID returns the commit ID of the HEAD of the repository.
func (r *LocalGitHistory) GetHeadCommitID() string {
	return r.HeadCommitID
}

func (r *LocalGitHistory) GetCommitParents(commitID string) ([]string, error) {
	panic("implement me")
	return []string{}, nil
}

func (r *LocalGitHistory) GetFileContent(filePath string, ref string) ([]byte, error) {
	panic("implement me")
	return []byte(""), nil
}
func (r *LocalGitHistory) Diff(currentCommit string, olderCommit string) ([]code.FileDiff, error) {
	panic("implement me")
	return []code.FileDiff{}, nil
}
