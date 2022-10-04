package models

// Scope is the scope of a test.
type Scope struct {
	Path      string
	Line      int
	FuncName  string
	File      string
	StartLine int
	StartCol  int
	EndLine   int
	EndCol    int
}

// TestRun is the result of a test that was run in this Nabaz run.
type TestRun struct {
	Name          string
	Success       bool
	TimeInMs      float64
	CallGraph     []Scope
	TestFuncScope *Scope
}
