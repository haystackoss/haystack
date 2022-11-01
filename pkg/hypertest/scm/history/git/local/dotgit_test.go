package local_test

import (
	"os/exec"
	"strings"
	"testing"

	"github.com/nabaz-io/nabaz/pkg/hypertest/scm/history/git/local"
)

func bashGitHead() string {
	cmd := "git rev-parse HEAD"
	out, err := exec.Command("bash", "-c", cmd).Output()
	if err != nil {
		panic(err)
	}
	return strings.TrimSpace((string(out)))
}

func TestHEADCommitID(t *testing.T) {
	// TestDotGit tests the dotgit implementation
	_, err := local.NewLocalGitHistory(".")
	if err != nil {
		t.Error(err)
	}
}

func TestCommitParents(t *testing.T) {
	// TestCommitParents tests the commit parents
	_, err := local.NewLocalGitHistory(".")
	if err != nil {
		t.Error(err)
	}
	// commitID := "ca536266ec486665a94cf1a409f18d5d4da90c59"
	// parentCommitID := "9322a4f4e42460659e0cd7d4ef4b716e2096f3bd"
	// parents, err := localGit.CommitParents(commitID)
	// if err != nil {
	// 	t.Error(err)
	// }

	// if len(parents) != 1 {
	// 	t.Error("Unexpected number of parents")
	// }

	// if parents[0] != parentCommitID {
	// 	t.Error("Unexpected parent commit ID")
	// }
}
