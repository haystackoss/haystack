package local

import (
	"errors"

	git "github.com/nabaz-io/go-git.v4"
	"github.com/nabaz-io/go-git.v4/plumbing"
	"github.com/nabaz-io/go-git.v4/plumbing/format/diff"
	"github.com/nabaz-io/nabaz/pkg/testrunner/scm/code"
)

// LocalGitHistory is history supplied by .git
type LocalGitHistory struct {
	// Path is the path to the local git repository.
	*git.Repository
	headCommitID string
	RootPath     string
}

// NewLocalGitHistory creates a new LocalGitRepo.
func NewLocalGitHistory(path string) (*LocalGitHistory, error) {
	git.GitDirName = ".git"
	originalDotGit, err := git.PlainOpenWithOptions(path, &git.PlainOpenOptions{DetectDotGit: true})
	if err != nil {
		return nil, err
	}
	conf, err := originalDotGit.Config()
	if err != nil {
		return nil, err
	}
	originalDotGitPath := conf.Core.Worktree

	git.GitDirName = ".nabazgit"
	gitRepo, err := git.PlainInit(originalDotGitPath, false)
	switch err {
	case nil:
		// do nothing
	case git.ErrRepositoryAlreadyExists:
		gitRepo, err = git.PlainOpen(originalDotGitPath)
		if err != nil {
			return nil, err
		}
	default:
		return nil, err
	}

	rootPath := originalDotGitPath
	localRepo := &LocalGitHistory{
		Repository: gitRepo,
		RootPath:   rootPath,
	}

	return localRepo, nil
}

// HeadCommitID returns the commit ID of the HEAD of the repository.
func (l *LocalGitHistory) HeadCommitID() (string, error) {
	if l.headCommitID == "" {
		head, err := l.Repository.Head()
		if err != nil {
			return "", errors.New("could not get HEAD commit")
		}
		l.headCommitID = head.Hash().String()
	}
	return l.headCommitID, nil
}

func (r *LocalGitHistory) CommitParents(commitID string) ([]string, error) {
	commit, err := r.Repository.CommitObject(plumbing.NewHash(commitID))
	if err != nil {
		return nil, err
	}

	parents := make([]string, 0)
	for _, parent := range commit.ParentHashes {
		parents = append(parents, parent.String())

	}
	return parents, nil
}

func (r *LocalGitHistory) GetFileContent(path string, commitID string) ([]byte, error) {
	hash := plumbing.NewHash(commitID)
	commit, err := r.Repository.CommitObject(hash)
	if err != nil {
		return nil, err
	}

	file, err := commit.File(path)
	if err != nil {
		return nil, err
	}

	content, err := file.Contents()
	if err != nil {
		return nil, err
	}

	return []byte(content), nil
}

func (l *LocalGitHistory) Diff(currentCommitID string, olderCommitID string) ([]code.FileDiff, error) {
	currentCommit, err := l.Repository.CommitObject(plumbing.NewHash(currentCommitID))
	if err != nil {
		return nil, err
	}

	olderCommit, err := l.Repository.CommitObject(plumbing.NewHash(olderCommitID))
	if err != nil {
		return nil, err
	}

	patch, err := currentCommit.Patch(olderCommit)
	if err != nil {
		return nil, err
	}

	patches := patch.FilePatches()

	fileDiffs := make([]code.FileDiff, len(patches))
	for i, patch := range patches {
		isBinary := patch.IsBinary()
		from, to := patch.Files()
		status := fileChangeNature(from, to)

		fileDiffs[i] = code.FileDiff{
			Path:         to.Path(),
			Patch:        patch.Chunks(),
			IsBinary:     isBinary,
			Status:       status,
			PreviousPath: from.Path(),
		}
	}
	return fileDiffs, nil

}

// fileChangeNature figures out whats the nature of the change, i.e. if the file was added, deleted or modified.
func fileChangeNature(from diff.File, to diff.File) code.FileStatus {
	if from == nil {
		return code.ADDED
	}
	if to == nil {
		return code.REMOVED
	}

	if from.Path() != to.Path() {
		return code.RENAMED
	}

	return code.MODIFIED
}

func (l *LocalGitHistory) HEAD() string {
	commitID, _ := l.HeadCommitID()
	return commitID
}
