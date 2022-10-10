package pytest

import (
	"fmt"
	"os/exec"
	"strings"
)

type Pytest struct {
	repoPath string
}

func NewPytestFramework(repoPath string) *Pytest {
	return &Pytest{
		repoPath: repoPath,
	}
}


func (p *Pytest) ListTests() map[string]string {

	cmd := exec.Command("pytest", "--collect-only", "--quiet", "--rootdir", p.repoPath)
	stdout, err := cmd.Output()
	if err != nil {
		panic(fmt.Errorf("WHILE LISTING PYTEST TESTS FOR %s GOT ERROR: %s", p.repoPath, err))
	}

	strStdout := string(stdout[:])
	lines := strings.Split(strStdout, "\n")
	lines = lines[:len(lines)-3]

	tests := make(map[string]string)
	for _, test := range lines {
		tests[test] = test
	}

	return tests
}