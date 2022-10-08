package framework

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"os"
	"os/exec"
	"strings"

	"github.com/nabaz-io/nabaz/pkg/testrunner/diffengine/parser"
	"github.com/nabaz-io/nabaz/pkg/testrunner/models"
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
	pkgsCache    map[string][]string
	coveragePath string
}
type GoTestEvent struct {
	Action   string  `json:"Action"`
	Output   string  `json:"Output"`
	Package  string  `json:"Package"`
	Test     string  `json:"Test"`
	Duration float64 `json:"Duration"`
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
	framework.pkgsCache = make(map[string][]string)
	framework.coveragePath = ""
	return framework
}

type FunctionCache struct {
	node     *sitter.Node
	fileName string
}

type PackageCache struct {
	testFilesToLoad []string
	functionsCache  map[string]FunctionCache
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
		panic("Listing tests failed with exit code " + string(exitCode) + " and stderr: " + string(stderr))
	}

	unparsedEvents := bytes.Split(stdout, []byte("\n"))
	events := make([]*GoTestEvent, len(unparsedEvents))
	for _, unparsedEvent := range unparsedEvents {
		if len(unparse 	dEvent) == 0 {
			continue
		}
		event := GoTestEvent{}
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

func (g *GoTest) RunTests(testsToSkip map[string][]*models.PreviousTestRun) ([]*models.TestRun, int) {
	/*
		def run_tests(self, tests_to_skip: dict[str, CachedTestResult]) -> list[str]:
			full_run = False

			with tempfile.NamedTemporaryFile(delete=False) as f:
				self._coverage_path = f.name
				self._logger.info(f"Using {self._coverage_path} as coverage file")

			if not tests_to_skip:
				full_run = True

			tests_found = self.tests.keys()  # Test names are keys in the dict
			tests_to_run = [test for test in tests_found if test not in tests_to_skip]
			tests_to_run_cmd = '|'.join([f"^{test}$" for test in tests_to_run])  # go test -run accepts regex
			pkgs_to_run_cmd = list(set([package for test_name, package in self.tests.items() if
							   test_name in tests_to_run]))  # go test -p accepts package names

			injectable_tests_to_run = tests_to_run_cmd if tests_to_run_cmd else "^$"

			args = inject_gotest_args(self.args, "-coverpkg", './...',
									  "-cover",
									  "-pertestcoverprofile", self._coverage_path,
									  "-json")
			if not full_run:
				# If we are not running all tests, we need to specify which tests to run
				args = inject_gotest_args(args, "-run", injectable_tests_to_run, *pkgs_to_run_cmd)
			else:
				args = inject_gotest_args(args, *self.pkgs)

			args = inject_gotest_args(["go", "test"], *args)

			self.test_run_time = time.time()
			p = run(args, capture_output=True, text=True, env=self.env)
			self.test_run_time = time.time() - self.test_run_time

			if p.returncode != 0:
				print(p.stderr)
				# we dont know if it's because of total failure or because of a specific test failing
				# so we wont exit just yet.

			output = ""
			tests_run_results = []
			for event in [json.loads(json_event) for json_event in p.stdout.splitlines()]:
				test_name = event.get('Test', '/')
				if 'Elapsed' in event and not is_subtest(test_name) and event.get('Action') in ['pass', 'fail', 'skip']:
					tests_run_results.append(event)
				if 'Output' in event:
					output += event.get('Output', '')

			print(output)

			cov = self._get_coverage_data()
			ran_tests = [
				TestResult(
					name=run_result.get('Test'),
					success=run_result.get('Action') == 'pass',
					time_in_ms=run_result.get('Elapsed'),
					call_graph=[Scope(**scope) for scope in cov.get(run_result.get('Test'), [])],
					test_func_scope=self.find_test_scope_in_pkg(run_result)
				)
				for run_result in tests_run_results
			]
			return ran_tests, p.returncode
	*/
	fullRun := false
	pertestcoverprofile, err := ioutil.TempFile("", "*")
	if err != nil {
		panic(err)
	}
	defer os.Remove(pertestcoverprofile.Name())

	g.coveragePath = pertestcoverprofile.Name()

	if len(testsToSkip) == 0 {
		fullRun = true
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
}

/*
def base_path(self):
		def GOPATH():
			p = run(["go", "env", "GOPATH"], capture_output=True, text=True, env=self.env)
			return p.stdout.strip()
		if not self.GOPATH:
			self.GOPATH = GOPATH()
		return self.GOPATH + "/src/"
*/
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

/*
	def find_test_scope_in_pkg(self, run_result: dict) -> Scope:
		pkg = run_result.get('Package')
		test_name = run_result.get('Test')

		# load package cache
		if pkg in self._pkgs_cache:
			package_cache = self._pkgs_cache[pkg]
		else:
			all_files = os.listdir(self.base_path() + pkg)
			test_files = [file for file in all_files if is_test_file(file)]
			package_cache = PackageCache(test_files_to_load=test_files, functions_cache={})
			self._pkgs_cache[pkg] = package_cache

		# if func already loaded in cache
		if matching_func := package_cache.functions_cache.get(test_name):
			path = pkg + '/' + matching_func.file_name
			return self.create_scope(node=matching_func.node, file_path=path, func_name=test_name)
		else:
			for test_file in package_cache.test_files_to_load:
				path = pkg + '/' + test_file
				with open(self.base_path() + path, 'rb') as f:
					code = f.read()
					# continue loading package's files into cache
					loaded_functions = self.parser.get_functions(code)
					new_function_to_cache = {func_name: FunctionCache(node=node, file_name=test_file)
											 for func_name, node in loaded_functions.items()}
					package_cache.functions_cache.update(new_function_to_cache)

					# remove file from files to load
					package_cache.test_files_to_load = [f for f in package_cache.test_files_to_load if f != test_file]

					if matching_func := new_function_to_cache.get(test_name):
						return self.create_scope(node=matching_func.node, file_path=path, func_name=test_name)

		raise Exception(f"Couldn't find scope for {test_name}")

	@staticmethod
	def create_scope(node: Node, file_path: str, func_name: str):
		return Scope(path=file_path, func_name=func_name, line=node.start_point[0],
					 startline=node.start_point[0], startcol=node.start_point[1],
					 endline=node.end_point[0], endcol=node.end_point[1])

	def base_path(self):
		def GOPATH():
			p = run(["go", "env", "GOPATH"], capture_output=True, text=True, env=self.env)
			return p.stdout.strip()
		if not self.GOPATH:
			self.GOPATH = GOPATH()
		return self.GOPATH + "/src/"

	def _get_coverage_data(self):
		raw_cov = parse_gotest_coverage_file(self._coverage_path)
		return raw_cov

def parse_gotest_coverage_file(path):
	with open(path, 'r') as f:
		buf = f.read()

	lines = buf.splitlines()
	mode_line = lines[0]
	mode = mode_line.split(':')[1]

	test_name = ""
	coverage_lines = lines[1:]
	coverage_data = {}
	for line in coverage_lines:
		splitted_line = line.split(':')
		if len(splitted_line) != 2:
			continue

		if splitted_line[0].strip() == START_NEW_TEST_MAGIC:
			test_name = splitted_line[1].strip()
			coverage_data[test_name] = []
			continue

		# pathto/file/name.go:line.column,line.column numberOfStatements count
		raw_coordinates, number_of_stmts, count = splitted_line[1].strip().split(' ')
		coordinates = raw_coordinates.split(',')
		start_coordinates = coordinates[0].split('.')
		end_coordinates = coordinates[1].split('.')
		count = int(count)

		if not (count > 0):
			continue

		if test_name not in coverage_data:
			coverage_data[test_name] = []

		coverage_data[test_name].append({
			"path": splitted_line[0].strip(),  # scope is the path
			"line": int(start_coordinates[0]),
			"startline": int(start_coordinates[0]),
			"startcol": int(start_coordinates[1]),
			"endline": int(end_coordinates[0]),
			"endcol": int(end_coordinates[1]),
		})

	return coverage_data
*/
