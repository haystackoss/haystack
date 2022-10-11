package pytest

import (
	"fmt"
	"os/exec"
	"strings"
)

type Pytest struct {
	repoPath string
	args     []string
}

func NewPytestFramework(repoPath string, args []string) *Pytest {
	return &Pytest{
		repoPath: repoPath,
		args:    args,
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

func (p *Pytest) RunTest(testsToSKip []string) string {
	args := []string{"-v", "--cov", "--cov-context=test", "-p", "pytest-json-report", "--json-report-file=/tmp/nabaz-pytest.json"}
	skipstr := ""
	if len(testsToSKip) > 0 {
		skipstr += fmt.Sprintf("not %s", testsToSKip[0])

		if len(testsToSKip) > 1 {
			for _, test := range testsToSKip[1:] {
				skipstr += fmt.Sprintf(" and not %s", test)
			}
		}
		args = append(args, "-k", skipstr)
	}

	args = append(args, p.args...)
	cmd := exec.Command("pytest", args...)

	stdout, err := cmd.Output()
	if err != nil {
		panic(fmt.Errorf("WHILE RUNNING PYTEST TEST WITH ARGS %s GOT ERROR: %s", args, err))
	}

	fmt.Print(string(stdout))

	return string(stdout[:])
}