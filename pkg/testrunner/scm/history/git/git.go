package git

import "github.com/nabaz-io/nabaz/pkg/testrunner/scm/code"

// History provided by git
type GitHistory interface {
	GetCommitParents(commitID string) ([]string, error)
	GetFileContent(path string, commitID string) ([]byte, error)
	Diff(currentCommitHash string, olderCommitHash string) ([]code.FileDiff, error)
}
