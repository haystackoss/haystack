package code

type Scope struct {
	Path     string
	Line     int
	FuncName string
	StartLine int
	StartCol int
	EndLine int
	EndCol int
}