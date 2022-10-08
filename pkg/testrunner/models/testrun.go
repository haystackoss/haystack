package models

// Scope is the scope of a test.
type Scope struct {
	Path      string `json:"path"`
	Line      int    `json:"line"`
	FuncName  string `json:"func_name"`
	File      string `json:"file"`
	StartLine int    `json:"start_line"`
	StartCol  int    `json:"start_col"`
	EndLine   int    `json:"end_line"`
	EndCol    int    `json:"end_col"`
}

type TestRun struct {
	Name          string  `json:"name"`
	Success       bool    `json:"success"`
	TimeInMs      float64 `json:"time_in_ms"`
	CallGraph     []Scope `json:"call_graph"`
	TestFuncScope Scope   `json:"test_func_scope"`
}
