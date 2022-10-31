package models

// NabazRun is the result of a run with nabaz.
type NabazRun struct {
	RunID           int64        `json:"result_id"`
	CommitID        string        `json:"commit_id"`
	TestsRan        []TestRun     `json:"tests_ran"`
	TestsSkipped    []SkippedTest `json:"tests_skipped"`
	RunDuration     float64       `json:"run_duration"`
	LongestDuration float64       `json:"longest_duration"`
}

// PreviousTestRun returns the previous test run outcome of a test that wasn't run this time around.
func (r *NabazRun) PreviousTestRun(testName string) *SkippedTest {
	for _, test := range r.TestsSkipped {
		if test.Name == testName {
			return &test
		}
	}
	return nil
}

// FailedTests returns the tests that failed in this Nabaz run.
func (r *NabazRun) FailedTests() []TestRun {
	var failedTests []TestRun
	for _, test := range r.TestsRan {
		if !test.Success {
			failedTests = append(failedTests, test)
		}
	}
	return failedTests
}

// GetTestRun returns the test run info for the given test name.
func (r *NabazRun) GetTestRun(testName string) *TestRun {
	for _, test := range r.TestsRan {
		if test.Name == testName {
			return &test
		}
	}
	return nil
}
