package pytest

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"os"
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
	tmpdir := os.TempDir()
	if tmpdir == "" {
		nomedir, err := os.UserHomeDir()
		if err != nil {
			tmpdir = "."
		} else {
			tmpdir = nomedir
		}
	}
	jsonPath := tmpdir + "/nabaz-pytest.json"
	jsonReportFile := fmt.Sprintf("--json-report-file=%s", jsonPath)
	// TODO validate json-report and cov plugin and suggest installing them if not installed
	args := []string{"-v", "--cov", "--cov-report=html", "--cov-context=test", "--json-report", jsonReportFile }
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

	stdout, _ := cmd.Output()
	// TODO: thats bad, handle when exit code is not 0
	// if err != nil {
	// 	panic(fmt.Errorf("WHILE RUNNING PYTEST TEST WITH ARGS %s GOT ERROR: %s", args, err))
	// }
	fmt.Print(string(stdout))

	// parse json file
	rawCoverage := readFileString(jsonPath)
	// print
	fmt.Println(rawCoverage)

	return string(stdout[:])
}

func readFileString(path string) string {
	file, err := os.Open(path)
	if err != nil {
		panic(fmt.Errorf("FAILED OT OPEN PER TEST CODE COVERAGE FILE: %s", err))
	}
	defer file.Close()

	reader := bufio.NewReader(file)
	bytes, err := ioutil.ReadAll(reader)
	if err != nil {
		panic(fmt.Errorf("FAILED TO READ PER TEST CODE COVERAGE FILE: %s", err))
	}

	return string(bytes)
}