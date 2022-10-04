package git

import "github.com/nabaz-io/nabaz/pkg/testrunner/scm/code"

// History provided by git
type GitHistory interface {
	GetCommitParents(commitID string) ([]string, error)
	GetFileContent(filePath string, commitID string) ([]byte, error)
	Diff(currentCommitID string, olderCommitID string) ([]code.FileDiff, error)
}
