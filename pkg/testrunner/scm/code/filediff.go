package code

import "gopkg.in/src-d/go-git.v4/plumbing/format/diff"

type FileStatus string

const (
	ADDED    FileStatus = "added"
	MODIFIED FileStatus = "modified"
	REMOVED  FileStatus = "removed"
	RENAMED  FileStatus = "renamed"
)

type FileDiff struct {
	Path         string
	Patch        []diff.Chunk
	IsBinary     bool
	Status       FileStatus
	PreviousPath string
}

func (f *FileDiff) IsRenamed() bool {
	return f.Path != f.PreviousPath
}

type ChangedFunction struct {
	Name string
	Path string
}
