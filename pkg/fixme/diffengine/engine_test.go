package diffengine_test

import (
	"github.com/nabaz-io/nabaz/pkg/fixme/diffengine"
	"github.com/nabaz-io/nabaz/pkg/fixme/diffengine/parser"
	"github.com/nabaz-io/nabaz/pkg/fixme/scm/code"
	historyfactory "github.com/nabaz-io/nabaz/pkg/fixme/scm/history/git/factory"
)

func NewTestDiffEngine() (*diffengine.DiffEngine, error) {
	codeDir := code.NewCodeDirectory(".")
	history, err := historyfactory.NewGitHistory(".")
	if err != nil {
		return nil, err
	}

	parser, err := parser.NewParser("go test")
	if err != nil {
		return nil, err
	}

	engine := diffengine.NewDiffEngine(codeDir, history, parser, "")
	if err != nil {
		return nil, err
	}

	return engine, nil
}
