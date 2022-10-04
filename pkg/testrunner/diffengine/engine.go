package diffengine

import (
	"log"

	"github.com/nabaz-io/nabaz/pkg/testrunner/scm/code"
	"github.com/nabaz-io/nabaz/pkg/testrunner/scm/history/git"
	sitter "github.com/smacker/go-tree-sitter"
)

type DiffEngine struct {
	parser  sitter.Parser
	history git.GitHistory
}

func NewDiffEngine(code code.CodeDirectory, history git.GitHistory) *DiffEngine {
	engine := &DiffEngine{}
	engine.parser = nil
	engine.history = history
	return engine
}

func (d *DiffEngine) Affects(path string) bool {
	parser := sitter.NewParser()

	return false
}
func (d *DiffEngine) ChangedFunctions(changedFiles []code.FileDiff) ([]string, error) {
	filePairs := make([]FilePair, 0)
	for _, fileDiff := range changedFiles {
		if fileDiff.IsBinary {
			log.Printf("file %s is binary, skipping...\n", fileDiff.Path)
			continue
		}

		if fileDiff.Status == code.MODIFIED {
			currentFile := d.localCod.GetFile(fileDiff.Path)
			oldFilePath := fileDiff.PreviousPath
			if oldFilePath == "" {
				oldFilePath = fileDiff.Path
			}
			oldFile, err := d.history.GetFileFromCommit(oldFilePath, d.P)
			if err != nil {
				return nil, err
			}
			filePairs = append(filePairs, FilePair{
				CurrentFile: currentFile,
				OldFile:     oldFile,
			})
		}
	}

	changedFunctions := make([]string, 0)

	for _, filePair := range filePairs {
		currFunctions := d.parser.GetFunctions(filePair.CurrentFile)
		oldFunctions := d.parser.GetFunctions(filePair.OldFile)

		for oldFuncName, oldFuncNode := range oldFunctions {
			matchingCurrentFunc := nextFuncNode(currFunctions, oldFuncName, oldFuncNode)
			if matchingCurrentFunc == nil {
				changedFunctions = append(changedFunctions, oldFuncName)
			}
		}
	}
}
