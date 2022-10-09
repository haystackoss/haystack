package framework

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"io/fs"
	"io/ioutil"
	"os"
	"os/exec"
	"strconv"
	"strings"

	"github.com/nabaz-io/nabaz/pkg/testrunner/diffengine/parser"
	"github.com/nabaz-io/nabaz/pkg/testrunner/models"
	"github.com/nabaz-io/nabaz/pkg/testrunner/scm/code"
	sitter "github.com/smacker/go-tree-sitter"
)

const (
	// GoTestFramework is the name of the go test framework
	goTestFramework      = "go test"
	START_NEW_TEST_MAGIC = "_testName"
)

type GoTest struct {
	parser       parser.Parser
	repoPath     string
	args         []string
	env          []string
	pkgs         []string
	GOPATH       string
	tests        map[string]string // map[testName]pkgPath
	testRunTime  float64
	pkgsCache    map[string]packageParseCache
	coveragePath string
}
type goTestResult struct {
	Action   string  `json:"Action"`
	Output   string  `json:"Output"`
	Package  string  `json:"Package"`
	Test     string  `json:"Test"`
	Duration float64 `json:"Duration"`
}
type functionCache struct {
	node     *sitter.Node
	fileName string
}
type packageParseCache struct {
	testFilesToParse []string
	functionsCache   map[string]functionCache
}

func setupGoEnv() []string {
	os.Setenv("GOROOT", "/usr/local/nabaz-go")
	os.Setenv("PATH", "/usr/local/nabaz-go/bin:"+os.ExpandEnv("$PATH"))
	return os.Environ()
}
func injectGoTestArgs(args []string, argsToInject ...string) []string {
	argsCopy := make([]string, len(args))
	copy(argsCopy, args)
	argsCopy = append(argsCopy, argsToInject...)
	return argsCopy

}
func isSubTest(name string) bool {
	// TestXXX/SubTestXXX
	return strings.Contains(name, "/")
}

func isTestFile(fileName string) bool {
	if !strings.HasSuffix(fileName, "_test.go") {
		return false
	}

	// main entrypoint of test package we don't care about it.
	if fileName == "main_test.go" {
		return false
	}

	return true

}

func NewGoTestFramework(languageParser parser.Parser, repoPath string, args string, pkgs string) *GoTest {
	framework := &GoTest{}
	framework.testRunTime = 0
	framework.repoPath = repoPath
	framework.args = strings.Split(args, " ")
	framework.env = setupGoEnv()
	framework.pkgs = strings.Split(pkgs, " ")
	framework.parser = languageParser
	framework.GOPATH = ""
	framework.tests = make(map[string]string)
	framework.pkgsCache = make(map[string]packageParseCache)
	framework.coveragePath = ""
	return framework
}

func run(args []string, env []string) (stdout []byte, stderr []byte, exitCode int) {
	var stdoutBuf bytes.Buffer
	var stderrBuf bytes.Buffer
	cmd := &exec.Cmd{}
	cmd.Env = env
	cmd.Path = args[0]
	cmd.Args = args
	cmd.Stdout = &stdoutBuf
	cmd.Stderr = &stderrBuf

	err := cmd.Run()
	if err != nil {
		if exitError, ok := err.(*exec.ExitError); ok {
			exitCode = exitError.ExitCode()
			return stdoutBuf.Bytes(), stderrBuf.Bytes(), exitCode
		}
	}

	stdout = stdoutBuf.Bytes()
	stderr = stderrBuf.Bytes()
	exitCode = cmd.ProcessState.ExitCode()
	return stdout, stderr, exitCode

}
func (g *GoTest) ListTests() map[string]string {
	if len(g.tests) > 0 {
		return g.tests
	}

	baseGoTestCmdline := []string{"go", "test", "-list", "Test", "-json"}
	finalCmdline := injectGoTestArgs(baseGoTestCmdline, g.args...)
	finalCmdline = injectGoTestArgs(finalCmdline, g.pkgs...)

	stdout, stderr, exitCode := run(finalCmdline, g.env)
	if exitCode != 0 {
		panic(fmt.Errorf("LISTING TESTS FAILED WITH EXIT CODE %d AND STDERR: %s", exitCode, string(stderr)))
	}

	unparsedEvents := bytes.Split(stdout, []byte("\n"))
	events := make([]*goTestResult, len(unparsedEvents))
	for _, unparsedEvent := range unparsedEvents {
		if len(unparsedEvent) == 0 {
			continue
		}
		event := goTestResult{}
		err := json.Unmarshal([]byte(unparsedEvent), &event)
		if err != nil {
			panic(err)
		}

		events = append(events, &event)
	}

	for _, event := range events {
		output := event.Output
		if strings.HasPrefix(output, "Test") {
			uniqueTestName := strings.TrimSpace(output)
			g.tests[uniqueTestName] = event.Package
		}
	}

	return g.tests
}

func (g *GoTest) RunTests(testsToSkip map[string][]models.SkippedTest) ([]models.TestRun, int) {
	fullRun := true
	pertestcoverprofile, err := ioutil.TempFile("", "*") // "" means use default temp dir native to OS
	if err != nil {
		panic(err)
	}
	defer os.Remove(pertestcoverprofile.Name())

	g.coveragePath = pertestcoverprofile.Name()

	if len(testsToSkip) > 0 {
		fullRun = false
	}

	testsFound := g.tests
	testsToRun := make([]string, 0, len(testsFound))
	pkgsToRun := make([]string, 0, len(testsFound))
	for test, pkg := range testsFound {
		// if test doesn't exist in testsToSkip, add it to the tests to run.
		if _, exists := testsToSkip[test]; !exists {
			testsToRun = append(testsToRun, test)
			pkgsToRun = append(pkgsToRun, pkg)
		}
	}

	for i, test := range testsToRun {
		testsToRun[i] = fmt.Sprintf("^%s$", test)
	}

	testsToRunCmd := strings.Join(testsToRun, "|")

	// we want to remove duplicates from pkgsToRun
	pkgsToRun = removeDuplicates(pkgsToRun)
	pkgsToRunCmd := strings.Join(pkgsToRun, " ")

	injectableTestsToRun := ""
	if testsToRunCmd != "" {
		injectableTestsToRun = testsToRunCmd
	} else {
		injectableTestsToRun = "^$"
	}

	args := injectGoTestArgs(g.args, "-coverpkg", "./...", "-cover", "-pertestcoverprofile", g.coveragePath, "-json")
	if !fullRun {
		args = injectGoTestArgs(args, "-run", injectableTestsToRun, pkgsToRunCmd)
	} else {
		args = injectGoTestArgs(args, g.pkgs...)
	}

	args = injectGoTestArgs([]string{"go", "test"}, args...)

	stdout, stderr, exitCode := run(args, g.env)

	if exitCode != 0 {
		fmt.Println(stderr)
	}

	output := ""
	testResults := make([]goTestResult, 0, len(testsFound))
	for _, jsonEvent := range bytes.Split(stdout, []byte("\n")) {
		testResult := goTestResult{}
		if err := json.Unmarshal(jsonEvent, &testResult); err != nil {
			fmt.Println(err)
			continue
		}

		if !isSubTest(testResult.Test) && (testResult.Action == "pass" || testResult.Action == "fail" || testResult.Action == "skip") {
			testResults = append(testResults, testResult)
		} else if testResult.Action == "Output" {
			output += testResult.Output
		}
	}

	// Print Output
	fmt.Println(output)

	// Get coverage data
	cov := g.getCoverageData()

	// Get test results
	ranTests := make([]models.TestRun, 0, len(testsFound))
	for _, testResult := range testResults {
		ranTests = append(ranTests, models.TestRun{
			Name:          testResult.Test,
			Success:       testResult.Action == "pass",
			TimeInMs:      testResult.Duration,
			CallGraph:     cov[testResult.Test],
			TestFuncScope: g.findTestScopeInPkg(testResult),
		})
	}
	return ranTests, exitCode
}

func (g *GoTest) getCoverageData() map[string][]code.Scope {

	rawCoverage := readFileString(g.coveragePath)
	lines := strings.Split(rawCoverage, "\n")
	modeLine := lines[0]
	_ = strings.Split(modeLine, ":")[1]

	testName := ""
	coverageLines := lines[1:]
	coverageData := make(map[string][]code.Scope)
	for _, line := range coverageLines {
		splittedLine := strings.Split(line, ":")
		if len(splittedLine) != 2 {
			continue
		}

		if strings.TrimSpace(splittedLine[0]) == START_NEW_TEST_MAGIC {
			testName = splittedLine[1]
			coverageData[testName] = make([]code.Scope, 0)
			continue
		}

		splittedInfo := strings.Split(splittedLine[1], " ")
		rawCoordinates, _, count := splittedInfo[0], splittedInfo[1], splittedInfo[2]

		coordinates := strings.Split(rawCoordinates, ",")
		startCoordinates := strings.Split(coordinates[0], ".")
		endCoordinates := strings.Split(coordinates[1], ".")
		countInt, err := strconv.Atoi(count)
		if err != nil {
			panic(err)
		}

		if countInt <= 0 {
			continue
		}

		if _, exists := coverageData[testName]; !exists {
			coverageData[testName] = make([]code.Scope, 0)
		}

		startLine, err := strconv.Atoi(startCoordinates[0])
		if err != nil {
			panic(fmt.Errorf("WHILE PARSING go test COVERAGE FILE %s GOT ERROR: %s", g.coveragePath, err))
		}

		startColumn, err := strconv.Atoi(startCoordinates[1])
		if err != nil {
			panic(fmt.Errorf("WHILE PARSING go test COVERAGE FILE %s GOT ERROR: %s", g.coveragePath, err))
		}

		endLine, err := strconv.Atoi(endCoordinates[0])
		if err != nil {
			panic(fmt.Errorf("WHILE PARSING go test COVERAGE FILE %s GOT ERROR: %s", g.coveragePath, err))
		}

		endColumn, err := strconv.Atoi(endCoordinates[1])
		if err != nil {
			panic(fmt.Errorf("WHILE PARSING go test COVERAGE FILE %s GOT ERROR: %s", g.coveragePath, err))
		}

		coverageData[testName] = append(coverageData[testName], code.Scope{
			Path:      splittedLine[0],
			StartLine: startLine,
			StartCol:  startColumn,
			EndLine:   endLine,
			EndCol:    endColumn,
		})
	}

	return coverageData

}

func (g *GoTest) findTestScopeInPkg(testResult goTestResult) *code.Scope {
	pkg := testResult.Package
	testName := testResult.Test

	// load package
	var currentPkgCache packageParseCache
	if _, exists := g.pkgsCache[pkg]; exists {
		currentPkgCache = g.pkgsCache[pkg]
	} else {
		allFiles, err := ioutil.ReadDir(g.BasePath() + pkg)
		if err != nil {
			panic(fmt.Errorf("WHILE READING DIRECTORY %s GOT ERROR: %s", g.BasePath()+pkg, err))
		}

		testFiles := filterTestFiles(allFiles)
		testFileNames := getTestFileNames(testFiles)

		currentPkgCache = packageParseCache{
			testFilesToParse: testFileNames,
			functionsCache:   make(map[string]functionCache),
		}
		g.pkgsCache[pkg] = currentPkgCache
	}

	// if func already parsed and loaded in cache
	if matchingFunc, exists := currentPkgCache.functionsCache[testName]; exists {
		path := pkg + "/" + matchingFunc.fileName
		return g.createScope(matchingFunc.node, path, testName)
	}

	for _, testFile := range currentPkgCache.testFilesToParse {
		path := pkg + "/" + testFile
		content, err := ioutil.ReadFile(g.BasePath() + path)
		if err != nil {
			panic(fmt.Errorf("WHILE READING FILE %s GOT ERROR: %s", g.BasePath()+path, err))
		}

		// continue loading package's files into cache
		loadedFunctions := g.parser.GetFunctions(content)
		newFunctionsToCache := make(map[string]functionCache)
		for funcName, node := range loadedFunctions {
			newFunctionsToCache[funcName] = functionCache{
				node:     node,
				fileName: testFile,
			}
		}
		currentPkgCache.functionsCache = mergeMaps(currentPkgCache.functionsCache, newFunctionsToCache)

		// remove file from files to parse
		currentPkgCache.testFilesToParse = removeElemFromList(currentPkgCache.testFilesToParse, testFile)

		if matchingFunc, exists := newFunctionsToCache[testName]; exists {
			return g.createScope(matchingFunc.node, path, testName)
		}
	}

	panic(fmt.Errorf("COULDNT FIND SCOPE FOR %s", testName))
}

func removeElemFromList(list []string, elem string) []string {
	for i, v := range list {
		if v == elem {
			return append(list[:i], list[i+1:]...)
		}
	}
	return list
}

func mergeMaps(m1, m2 map[string]functionCache) map[string]functionCache {
	for k, v := range m2 {
		m1[k] = v
	}
	return m1
}

func (g *GoTest) createScope(node *sitter.Node, filePath string, funcName string) *code.Scope {
	return &code.Scope{
		Path:      filePath,
		FuncName:  funcName,
		StartLine: int(node.StartPoint().Row),
		StartCol:  int(node.StartPoint().Column),
		EndLine:   int(node.EndPoint().Row),
		EndCol:    int(node.EndPoint().Column),
	}
}

func getTestFileNames(testFiles []os.FileInfo) []string {
	testFileNames := make([]string, 0)
	for _, testFile := range testFiles {
		testFileNames = append(testFileNames, testFile.Name())
	}
	return testFileNames
}

func filterTestFiles(allFiles []fs.FileInfo) []fs.FileInfo {
	testFiles := make([]fs.FileInfo, 0)
	for _, file := range allFiles {
		if isTestFile(file.Name()) {
			testFiles = append(testFiles, file)
		}
	}

	return testFiles
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

func removeDuplicates(slice []string) []string {
	keys := make(map[string]bool)
	list := []string{}
	for _, entry := range slice {
		if _, value := keys[entry]; !value {
			keys[entry] = true
			list = append(list, entry)
		}
	}
	return list
}

func (g *GoTest) BasePath() string {
	if g.GOPATH == "" {
		stdout, stderr, exitCode := run([]string{"go", "env", "GOPATH"}, g.env)
		if exitCode != 0 {
			panic(stderr)
		}
		g.GOPATH = strings.TrimSpace(string(stdout))
	}
	return g.GOPATH + "/src/"
}
