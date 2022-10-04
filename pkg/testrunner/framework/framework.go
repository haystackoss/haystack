package framework

import "github.com/nabaz-io/nabaz/pkg/testrunner/models"

// class TestFramework:
//     def __init__(self):
//         raise NotImplementedError()

//     def list_tests(self) -> list[str]:
//         raise NotImplementedError()

//     def run_tests(self, tests_to_skip: dict[str, CachedTestResult]) -> list[TestResult]:
//         raise NotImplementedError()

//     def base_path(self) -> str:
//         raise NotImplementedError()

type Framework interface {
	ListTests() []string
	RunTests(testsToSkip map[string][]models.PreviousTestRun) []models.TestRun
	BasePath() string
}
