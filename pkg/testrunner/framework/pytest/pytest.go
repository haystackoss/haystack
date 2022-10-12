package pytest

import (
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"strings"

	"github.com/nabaz-io/nabaz/pkg/testrunner/models"
)

type Pytest struct {
	repoPath string
	args     []string
}

func NewPytestFramework(repoPath string, args string) *Pytest {
	return &Pytest{
		repoPath: repoPath,
		args:     strings.Split(args, " "),
	}
}

func (p *Pytest) BasePath() string {
	return ""
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

func (p *Pytest) RunTests(testsToSKip map[string]models.SkippedTest) ([]models.TestRun, int) {
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

	// TODO  suggest installing packages if not installed

	formattedTestsToSkip := "{"
	for testName := range testsToSKip {
		formattedTestsToSkip += fmt.Sprintf("\"%s\":true,", testName)
	}
	if len(formattedTestsToSkip) > 1 {
		formattedTestsToSkip = formattedTestsToSkip[:len(formattedTestsToSkip)-1] + "}" // remove last comma
	} else {
		formattedTestsToSkip = "{}"
	}

	// TODO: cp plugin to tmp
	args := []string{"/usr/local/bin/_nabazpytestplugin.py", jsonPath, formattedTestsToSkip, "--rootdir", p.repoPath}
	args = injectArgs(args, p.args...)

	cmd := exec.Command("python3", args...)
	cmdReader, err := cmd.StdoutPipe()
	cmd.Stderr = cmd.Stdout

	cmd.Start()
	ch := make(chan struct{}, 1)
	ctx, cancel := context.WithCancel(context.Background())

	go func(ctx context.Context) {
		scanner := bufio.NewScanner(cmdReader)

		for {
			select {
			case <-ctx.Done():
				ch <- struct{}{}
				return

			default:
				if scanner.Scan() {
					fmt.Println(scanner.Text())
				}
			}
		}
	}(ctx)

	err = cmd.Wait()
	cancel()
	<-ch

	exitCode := cmd.ProcessState.ExitCode()

	if err != nil {
		if strings.Contains(err.Error(), "exit status 1") {
			// do nth
		} else {
			panic(fmt.Errorf("WHILE RUNNING PYTEST TEST WITH ARGS %s GOT ERROR: %s", args, err))
		}
	}

	rawMapOfStrToTestRun := readFileString(jsonPath)

	testRuns := make(map[string]models.TestRun)
	json.Unmarshal([]byte(rawMapOfStrToTestRun), &testRuns)

	var tests []models.TestRun
	for _, test := range testRuns {
		tests = append(tests, test)
	}

	return tests, exitCode
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
