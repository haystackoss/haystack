package code

type FileStatus string

const (
	ADDED    FileStatus = "added"
	MODIFIED FileStatus = "modified"
	REMOVED  FileStatus = "removed"
	RENAMED  FileStatus = "renamed"
)

type FileDiff struct {
	Path         string
	Patch        string
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
