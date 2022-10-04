package git

import "github.com/nabaz-io/nabaz/pkg/testrunner/scm/code"

// History provided by git
type GitHistory interface {
	GetCommitParents(commitID string) ([]string, error)
	GetFileContent(filePath string, ref string) ([]byte, error)
	Diff(currentCommit string, olderCommit string) ([]code.FileDiff, error)
}
