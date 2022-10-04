package golang

import (
	"errors"

	"github.com/nabaz-io/nabaz/pkg/testrunner/scm/code"
	sitter "github.com/smacker/go-tree-sitter"
	"github.com/smacker/go-tree-sitter/golang"
)

type GolangParser struct {
	golangSyntax *sitter.Language
	parser       *sitter.Parser
}

// TODO: GenerateTree and FindFunction are same in every parser, need to find better abstraction !!!
// TODO: in addition, add caching mechanism for functions from python-test-runner after you find a better abstraction

func NewGolangParser() (*GolangParser, error) {
	return &GolangParser{
		golangSyntax: golang.GetLanguage(),
		parser:       sitter.NewParser(),
	}, nil
}

func (p *GolangParser) GenerateTree(code []byte) *sitter.Tree {
	// TODO: .parse is deprecated: use ParseCtx instead, read about it
	return p.parser.Parse(nil, code)
}

func (p *GolangParser) GetFunctions(code []byte) map[string]*sitter.Node {
	tree := p.GenerateTree(code)
	n := tree.RootNode()

	// funcs query
	q, _ := sitter.NewQuery([]byte(`(function_declaration "func" @structure.anchor)`), p.golangSyntax)
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
			func_name := c.Node.NextNamedSibling().String()
			functions[func_name] = c.Node.Parent()
		}
	}

	// methods query
	q2, _ := sitter.NewQuery([]byte(`(method_declaration "func" @structure.anchor)`), p.golangSyntax)
	qc2 := sitter.NewQueryCursor()
	qc2.Exec(q2, n)
	for {
		m, ok := qc2.NextMatch()
		if !ok {
			break
		}

		for _, c := range m.Captures {
			func_name := c.Node.NextNamedSibling().NextNamedSibling().String()
			functions[func_name] = c.Node.Parent()
		}
	}

	return functions
}

func (p *GolangParser) FindFunction(code []byte, scope *code.Scope) (string, error) {
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
