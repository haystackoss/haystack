package git

import "github.com/nabaz-io/nabaz/pkg/testrunner/scm/code"

// History provided by git
type GitHistory interface {
	GetCommitParents(commitID string) ([]string, error)
	GetFileFromCommit(filePath string, ref string) ([]byte, error)
	Diff(currentCommit string, olderCommit string) ([]code.FileDiff, error)
}
