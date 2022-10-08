package models

import "github.com/nabaz-io/nabaz/pkg/testrunner/scm/code"

type TestRun struct {
	Name          string       `json:"name"`
	Success       bool         `json:"success"`
	TimeInMs      float64      `json:"time_in_ms"`
	CallGraph     []code.Scope `json:"call_graph"`
	TestFuncScope code.Scope   `json:"test_func_scope"`
}
