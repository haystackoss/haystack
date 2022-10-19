package python

import (
	"context"
	"errors"
	"fmt"

	"github.com/nabaz-io/nabaz/pkg/fixme/scm/code"
	sitter "github.com/smacker/go-tree-sitter"
	"github.com/smacker/go-tree-sitter/python"
)

type PythonParser struct {
	pySyntax *sitter.Language
	parser   *sitter.Parser
}

func NewPythonParser() (*PythonParser, error) {
	pySyntax := python.GetLanguage()
	parser := sitter.NewParser()
	parser.SetLanguage(pySyntax)

	return &PythonParser{
		pySyntax: pySyntax,
		parser:   parser,
	}, nil
}

func (p *PythonParser) GenerateTree(code []byte) (*sitter.Tree, error) {
	return p.parser.ParseCtx(context.Background(), nil, code)
}

func (p *PythonParser) GetFunctions(code []byte) map[string]*sitter.Node {
	tree, err := p.GenerateTree(code)
	if err != nil {
		panic(fmt.Errorf("failed to parse python code " + err.Error()))
	}

	n := tree.RootNode()

	q, _ := sitter.NewQuery([]byte(`(function_definition name: (identifier) @function.def)`), p.pySyntax)
	qc := sitter.NewQueryCursor()
	qc.Exec(q, n)

	functions := make(map[string]*sitter.Node)
	for {
		m, ok := qc.NextMatch()
		if !ok {
			break
		}

		for _, c := range m.Captures {
			functions[c.Node.Content(code)] = c.Node.Parent()
		}
	}

	return functions
}

func (p *PythonParser) FindFunction(code []byte, scope *code.Scope) (string, error) {
	functions := p.GetFunctions(code)

	for func_name, func_node := range functions {
		x1 := func_node.StartPoint().Row
		x2 := func_node.EndPoint().Row

		real_lineo := uint32(scope.StartLine - 1)
		if x1 <= real_lineo && real_lineo <= x2 {
			return func_name, nil
		}
	}

	return "", errors.New("Function not found")
}
