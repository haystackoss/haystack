package framework

import (
	"fmt"

	"github.com/nabaz-io/nabaz/pkg/fixme/diffengine/parser"
	"github.com/nabaz-io/nabaz/pkg/fixme/framework/pytest"

	"github.com/nabaz-io/nabaz/pkg/fixme/models"
)

type Framework interface {
	ListTests() map[string]string // map[testName]packageName
	RunTests(testsToSkip map[string]models.SkippedTest) ([]models.TestRun, int)
	BasePath() string
}

func NewFramework(languageParser parser.Parser, framework, repoPath, testArgs, pkgs string) (Framework, error) {
	if framework == "pytest" {
		return pytest.NewPytestFramework(repoPath, testArgs), nil
	} else if framework == "go test" {
		return NewGoTestFramework(languageParser, repoPath, testArgs, pkgs), nil
	}

	return nil, fmt.Errorf("UNKNOWN FRAMEWORK \"%s\" PROVIDED, nabaz CURRENTLY SUPPORTS PYTEST AND GO TEST ONLY", framework)
}
