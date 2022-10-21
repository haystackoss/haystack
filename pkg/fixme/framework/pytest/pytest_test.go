package pytest_test

import (
	"testing"

	"github.com/nabaz-io/nabaz/pkg/fixme/framework/pytest"
	"github.com/nabaz-io/nabaz/pkg/fixme/models"
)

func TestListTests(t *testing.T) {
	framework := pytest.NewPytestFramework("./", "-v")
	tests, _ := framework.ListTests()
	if len(tests) != 3 {
		t.Errorf("Expected 3 tests, got %d", len(tests))
	}
}

func TestRunTests(t *testing.T) {
	framework := pytest.NewPytestFramework(".", "-v")
	testsToSKip := make(map[string]models.SkippedTest)
	// append "test_file.py::test_validate_user_agent_bad"
	testsToSKip["test_file.py::test_validate_user_agent_bad"] = models.SkippedTest{}

	testsRuns, exitCode := framework.RunTests(testsToSKip)
	if len(testsRuns) != 2 {
		t.Errorf("Expected 2 test run, got %d", len(testsRuns))
	}
	if exitCode != 1 {
		t.Errorf("Expected exit code 1, got %d", exitCode)
	}
}
