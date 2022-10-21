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

	"github.com/nabaz-io/nabaz/pkg/fixme/models"
	"github.com/nabaz-io/nabaz/pkg/fixme/paths"
)

type Pytest struct {
	repoPath string
	args     []string
}

func NewPytestFramework(repoPath string, args string) *Pytest {
	// validate any python3 installed
	if _, err := exec.LookPath("python3"); err != nil {
		fmt.Println("Error: can't run nabaz with pytest, python3 is not installed")
		os.Exit(1)
	}

	return &Pytest{
		repoPath: repoPath,
		args:     strings.Split(args, " "),
	}
}

func (p *Pytest) BasePath() string {
	return ""
}

func (p *Pytest) ListTests() (map[string]string, error) {

	cmd := exec.Command("pytest", "--collect-only", "--quiet", "--rootdir", p.repoPath)
	stdout, err := cmd.Output()
	exitCode := cmd.ProcessState.ExitCode()

	if exitCode == 2 || exitCode == 3 || exitCode == 4 { // pytest error
		panic(fmt.Errorf("FAILED TO LIST TESTS, USER ERROR: %s, stdout: %s", err, stdout))
	}

	if exitCode == 5 { // no tests collected
		return map[string]string{}, nil
	}

	if err != nil {
		panic(fmt.Errorf("WHILE LISTING PYTEST TESTS FOR %s GOT ERROR: %s, stdout %s, exit code %d", p.repoPath, err, stdout, exitCode))
	}

	strStdout := string(stdout[:])
	lines := strings.Split(strStdout, "\n")
	lines = lines[:len(lines)-3]

	tests := make(map[string]string)
	for _, test := range lines {
		tests[test] = test
	}

	return tests, nil
}

func (p *Pytest) RunTests(testsToSKip map[string]models.SkippedTest) (testRuns []models.TestRun, exitCode int) {
	tmpdir := paths.TempDir()
	jsonPath := tmpdir + "/pytest-results.json"

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
	args := []string{"/usr/local/bin/_nabazpytestplugin.py", jsonPath, paths.JunitXMLPath(), formattedTestsToSkip, "--rootdir", p.repoPath}
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
					// DON'T 
					// fmt.Println(scanner.Text())
				}
			}
		}
	}(ctx)

	err = cmd.Wait()
	cancel()
	<-ch

	exitCode = cmd.ProcessState.ExitCode()

	if err != nil {
		if exitCode == 1 || exitCode == 5 { // tests failed or no tests collected
			// do nth
		} else {
			panic(fmt.Errorf("WHILE RUNNING PYTEST TEST WITH ARGS %s GOT ERROR: %s", args, err))
		}
	}

	rawMapOfStrToTestRun := readFileString(jsonPath)

	testMap := make(map[string]models.TestRun)
	json.Unmarshal([]byte(rawMapOfStrToTestRun), &testMap)

	var tests []models.TestRun
	for _, test := range testMap {
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
