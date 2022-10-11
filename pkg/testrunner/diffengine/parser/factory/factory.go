package factory

import (
	"fmt"

	"github.com/nabaz-io/nabaz/pkg/testrunner/diffengine/parser"
	"github.com/nabaz-io/nabaz/pkg/testrunner/diffengine/parser/golang"
	"github.com/nabaz-io/nabaz/pkg/testrunner/diffengine/parser/python"
)

func NewParser(framework string) (parser.Parser, error) {
	if framework == "pytest" {
		return python.NewPythonParser()
	} else if framework == "go test" {
		return golang.NewGolangParser()
	} else {
		return nil, fmt.Errorf("UNKNOWN TEST FRAMEWORK \"%s\" PROVIDED, test-runner CURRENTLY SUPPORTS pytest AND gotest ONLY.", framework)
	}
}
