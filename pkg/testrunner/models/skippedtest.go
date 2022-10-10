package models

type SkippedTest struct {
	Name     string `json:"name"`
	RunIDRef int64 `json:"run_id_reference"`
}
