package golang

import (
	"context"
	"errors"
	"fmt"

	"github.com/nabaz-io/nabaz/pkg/testrunner/scm/code"
	sitter "github.com/smacker/go-tree-sitter"
	"github.com/smacker/go-tree-sitter/golang"
)

type GolangParser struct {
	golangSyntax *sitter.Language
	parser       *sitter.Parser
}

// TODO: in addition, add caching mechanism for functions from python-test-runner after you find a better abstraction

func NewGolangParser() (*GolangParser, error) {
	golangSyntax := golang.GetLanguage()
	parser := sitter.NewParser()
	parser.SetLanguage(golangSyntax)

	return &GolangParser{
		golangSyntax: golangSyntax,
		parser:       parser,
	}, nil
}

func (p *GolangParser) GenerateTree(code []byte) (*sitter.Tree, error) {
	return p.parser.ParseCtx(context.Background(), nil, code)
}

func (p *GolangParser) GetFunctions(code []byte) map[string]*sitter.Node {
	tree, err := p.GenerateTree(code)
	if err != nil {
		panic(fmt.Errorf("FAILED TO PARSE GOLANG CODE " + err.Error()))
	}

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
			func_name := c.Node.NextNamedSibling().Content(code)
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
			func_name := c.Node.NextNamedSibling().NextNamedSibling().Content(code)
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

		real_lineo := uint32(scope.StartLine - 1)
		if x1 <= real_lineo && real_lineo <= x2 {
			return func_name, nil
		}
	}

	return "", errors.New("FUNCTION NOT FOUND")
}
