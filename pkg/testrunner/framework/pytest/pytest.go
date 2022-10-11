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
	jsonPath := tmpdir + "/output.json"

	// TODO validate suggest installing packages if not installed
	// python3 plugin.py /tmp/output.json '{"test_file.py::test_validate_user_agent_bad":true,"test_file.py::test_validate_user_agent_chrome_good":true}

	formattedTestsToSkip := "'{"
	for _, test := range testsToSKip {
		formattedTestsToSkip += fmt.Sprintf("\"%s\":true,", test)
	}
	formattedTestsToSkip = formattedTestsToSkip[:len(formattedTestsToSkip)-1] + "}'" // remove last comma

	args := []string{"plugin.py", jsonPath, formattedTestsToSkip, "--rootdir", p.repoPath}
	args = injectArgs(args, p.args...)

	cmd := exec.Command("python3", args...)

	stdout, err := cmd.Output()
	// TODO: we don't get stdout from python, get it.
	// TOOO: handle err, 1 means test failed, 0 means test passed, don't panic on 1?
	fmt.Println("stdout: ", string(stdout))
	if err != nil {
		panic(fmt.Errorf("WHILE RUNNING PYTEST TEST WITH ARGS %s GOT ERROR: %s", args, err))
	}
	

	// parse json file
	rawCoverage := readFileString(jsonPath)
	// print
	fmt.Println(rawCoverage)

	return string(stdout[:])
}

func injectArgs(args []string, argsToInject ...string) []string {
	argsCopy := make([]string, len(args))
	copy(argsCopy, args)
	argsCopy = append(argsCopy, argsToInject...)
	return argsCopy

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