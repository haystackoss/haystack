package pytest_test

import (
	"testing"

	"github.com/nabaz-io/nabaz/pkg/testrunner/framework/pytest"
)


func TestListTests(t *testing.T) {
	framework := pytest.NewPytestFramework("./pythonrepo/")
	tests := framework.ListTests()
	if len(tests) != 3 {
		t.Errorf("Expected 3 tests, got %d", len(tests))
	}
}