package python

import (
	"errors"

	"github.com/nabaz-io/nabaz/pkg/testrunner/scm/code"
	sitter "github.com/smacker/go-tree-sitter"
	"github.com/smacker/go-tree-sitter/python"
)

// TODO: GenerateTree and FindFunction are same in every parser, need to find better abstraction !!!

type PythonParser struct {
	python *sitter.Language
	parser *sitter.Parser
}

func NewPythonParser() (*PythonParser, error) {
	return &PythonParser{
		python: python.GetLanguage(),
		parser: sitter.NewParser(),
	}, nil
}

func (p *PythonParser) GenerateTree(code []byte) *sitter.Tree {
	// TODO: .parse is deprecated: use ParseCtx instead, read about it  
	return p.parser.Parse(nil, code)
}

func (p *PythonParser) GetFunctions(code []byte) map[string]*sitter.Node {
	tree := p.GenerateTree(code)
	n := tree.RootNode()

	q, _ := sitter.NewQuery([]byte(`(function_definition name: (identifier) @function.def)`) , p.python)
	qc := sitter.NewQueryCursor()
	qc.Exec(q, n)

	functions := make(map[string]*sitter.Node)
	for {
		m, ok := qc.NextMatch()
		if !ok {
			break
		}

		for _, c := range m.Captures {
			// TODO: make sure to use func_name as key and func_node as value, need to run it
			functions[c.Node.String()] = c.Node.Parent()
		}
	}	

	return functions
}

func (p *PythonParser) FindFunction(code []byte, scope *code.Scope) (string, error) {
	functions := p.GetFunctions(code)

	for func_name, func_node := range functions {
		x1 := func_node.StartPoint().Row
		x2 := func_node.EndPoint().Row

		real_lineo := uint32(scope.Line - 1)
		if x1 <= real_lineo && real_lineo <= x2 {
			return func_name, nil
		}
	}

	return "", errors.New("Function not found")
}