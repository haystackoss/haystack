package models

type SkippedTest struct {
	Name           string `json:"name"`
	RunIDReference uint64 `json:"run_id_reference"`
}
