package models

type Telemetry struct {
	RepoName string `json:"repo_name"`
	RunDuration float64 `json:"run_duration"`
	LongestDuration float64 `json:"longest_duration"`
	SkippedTests int `json:"skipped_tests"`
	RanTests int `json:"ran_tests"`
	FailedTests int `json:"failed_tests"`
}
