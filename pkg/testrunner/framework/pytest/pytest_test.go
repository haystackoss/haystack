package pytest_test

import (
	"testing"

	"github.com/nabaz-io/nabaz/pkg/testrunner/framework/pytest"
)


func TestListTests(t *testing.T) {
	framework := pytest.NewPytestFramework("./pythonrepo/", []string{"-v"})
	tests := framework.ListTests()
	if len(tests) != 3 {
		t.Errorf("Expected 3 tests, got %d", len(tests))
	}
}

func TestRunTest(t *testing.T) {
	framework := pytest.NewPytestFramework("./pythonrepo/", []string{"-v"})
	_ = framework.RunTest([]string{"test_validate_user_agent_bad***"})
	t.Error()
}