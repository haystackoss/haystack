package models

type Telemetry struct {
	RepoName string `json:"repo_name"`
	RunDuration float64 `json:"run_duration"`
	LongestDuration float64 `json:"longest_duration"`
	TestsSkipped int `json:"skipped_tests"`
	TestsRan int `json:"ran_tests"`
	TestsFailed int `json:"failed_tests"`
}
