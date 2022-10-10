package models

type SkippedTest struct {
	Name     string `json:"name"`
	RunIDRef uint64 `json:"run_id_reference"`
}
