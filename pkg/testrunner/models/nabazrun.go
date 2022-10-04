package models

// NabazRun is the result of a run with nabaz.
type NabazRun struct {
	RunID                 string            `bson:"result_id"`
	CommitID              string            `bson:"commit_id"`
	TestRuns              []TestRun         `bson:"test_runs,omitempty"`
	PreviousTestRuns      []PreviousTestRun `bson:"previous_test_runs,omitempty"`
	RunDuration           float64           `bson:"run_duration"`
	CachedLongestDuration float64           `bson:"cached_longest_duration,"`
}

// PreviousOutcome returns the previous outcome of the given test.
func (r *NabazRun) PreviousOutcome(testName string) *PreviousTestRun {
	for _, test := range r.PreviousTestRuns {
		if test.Name == testName {
			return &test
		}
	}
	return nil
}

// SkippedTests returns tests that were skipped in this run.
func (r *NabazRun) SkippedTests() []PreviousTestRun {
	var skippedTests []PreviousTestRun
	for _, test := range r.PreviousTestRuns {
		if !test.Ran {
			skippedTests = append(skippedTests, test)
		}
	}
	return skippedTests
}

// FailedTests returns the tests that failed in this Nabaz run.
func (r *NabazRun) FailedTests() []TestRun {
	var failedTests []TestRun
	for _, test := range r.TestRuns {
		if !test.Success {
			failedTests = append(failedTests, test)
		}
	}
	return failedTests
}

// GetTestRun returns the test run info for the given test name.
func (r *NabazRun) GetTestRun(testName string) *TestRun {
	for _, test := range r.TestRuns {
		if test.Name == testName {
			return &test
		}
	}
	return nil
}
