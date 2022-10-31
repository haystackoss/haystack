package code

type Scope struct {
	Path      string `json:"path"`
	FuncName  string `json:"func_name"`
	StartLine int    `json:"startline"`
	StartCol  int    `json:"startcol"`
	EndLine   int    `json:"endline"`
	EndCol    int    `json:"endcol"`
}
