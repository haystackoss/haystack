package parser

import (
	"fmt"

	"github.com/nabaz-io/nabaz/pkg/adhdtest/diffengine/parser/golang"
	"github.com/nabaz-io/nabaz/pkg/adhdtest/diffengine/parser/python"
)

func NewParser(framework string) (Parser, error) {
	if framework == "pytest" {
		return python.NewPythonParser()
	} else if framework == "go test" {
		return golang.NewGolangParser()
	} else {
		return nil, fmt.Errorf("UNKNOWN TEST FRAMEWORK \"%s\" PROVIDED, test-runner CURRENTLY SUPPORTS 'pytest' AND 'go test' ONLY", framework)
	}

}
