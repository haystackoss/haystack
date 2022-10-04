package diffengine

import (
	"log"

	"github.com/nabaz-io/nabaz/pkg/testrunner/scm/code"
	"github.com/nabaz-io/nabaz/pkg/testrunner/scm/history/git"
	sitter "github.com/smacker/go-tree-sitter"
)

type DiffEngine struct {
	parser    *sitter.Parser
	history   git.GitHistory
	localCode *code.CodeDirectory
}

type FilePair struct {
	CurrentFile string
	OldFile     string
}

func NewDiffEngine(code *code.CodeDirectory, history git.GitHistory, oldCommitID string, parser *LanaguageParser *DiffEngine {
	engine := &DiffEngine{}
	engine.parser = nil
	engine.history = history
	engine.localCode = code

	return engine
}

func (d *DiffEngine) Affects(changedFunctions []string, coverage) bool {
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
			currentFile, err := d.localCode.GetFileContent(fileDiff.Path)
			if err != nil {
				return nil, err
			}

			oldFilePath := fileDiff.PreviousPath
			if oldFilePath == "" {
				oldFilePath = fileDiff.Path
			}
			commitID := "TODO: REPLACE ME"
			oldFile, err := d.history.GetFileContent(oldFilePath, commitID)
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
