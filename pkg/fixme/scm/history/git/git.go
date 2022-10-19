package git

import "github.com/nabaz-io/nabaz/pkg/fixme/scm/code"

// History provided by git
type GitHistory interface {
	SaveAllFiles() error
	CommitParents(commitID string) ([]string, error)
	GetFileContent(filePath string, commitID string) ([]byte, error)
	Diff(currentCommitID string, olderCommitID string) ([]code.FileDiff, error)
	HEAD() string
}
