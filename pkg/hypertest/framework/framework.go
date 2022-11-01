package framework

import (
	"fmt"
	"strings"

	"github.com/nabaz-io/nabaz/pkg/hypertest/diffengine/parser"
	"github.com/nabaz-io/nabaz/pkg/hypertest/framework/pytest"

	"github.com/nabaz-io/nabaz/pkg/hypertest/models"
)

type Framework interface {
	ListTests() (map[string]string, error) // map[testName]packageName
	RunTests(testsToSkip map[string]models.SkippedTest) (testRuns []models.TestRun, exitCode int)
	BasePath() string
}

func NewFramework(languageParser parser.Parser, framework, repoPath, testArgs string) (Framework, error) {
	if framework == "pytest" {
		return pytest.NewPytestFramework(repoPath, testArgs), nil
	} else if framework == "go test" {
		return NewGoTestFramework(languageParser, repoPath, testArgs), nil
	}

	return nil, fmt.Errorf("UNKNOWN FRAMEWORK \"%s\" PROVIDED, nabaz CURRENTLY SUPPORTS PYTEST AND GO TEST ONLY", framework)
}

func TestFileExtensionFromError(err string) string {
	if strings.Contains(err, ".py") {
		return "py"
	} else if strings.Contains(err, ".go") {
		return "go"
	} else {
		return ""
	}
}

func TestNameToFileLink(framework string, testResults []models.TestRun) map[string]string {
	testNameToTestFilePath := make(map[string]string)
	for _, test := range testResults {
		if !test.Success {
			filePath := ""
			line := -1

			if framework == "pytest" {
				if len(test.CallGraph) > 0 {
					scope := test.CallGraph[0]
					filePath = scope.Path
					line = scope.StartLine
				}
			} else if framework == "go test" {
				if test.TestFuncScope != nil {
					scope := test.TestFuncScope
					filePath = scope.Path
					line = scope.StartLine
				}
			}

			if filePath != "" {
				splitted := strings.Split(filePath, "/")
				fileName := splitted[len(splitted)-1]
				testNameToTestFilePath[test.Name] = fmt.Sprintf("%s:%d", fileName, line+1)
			}
		}
	}

	return testNameToTestFilePath
}
