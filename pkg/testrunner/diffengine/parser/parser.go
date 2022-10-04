package parser

import (
	"github.com/nabaz-io/nabaz/pkg/testrunner/models"
	sitter "github.com/smacker/go-tree-sitter"
)

type Parser interface {
	GenerateTree(code []byte) *sitter.Tree
	GetFunctions(code []byte) map[string]*sitter.Node
	FindFunction(code []byte, scope *models.Scope) string
}

func NewParser() {

}
