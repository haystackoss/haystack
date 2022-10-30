package parser

import (
	sitter "github.com/smacker/go-tree-sitter"
)

type Parser interface {
	GenerateTree(code []byte) (*sitter.Tree, error)
	GetFunctions(code []byte) map[string]*sitter.Node
}
