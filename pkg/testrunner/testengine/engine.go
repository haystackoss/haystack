package testengine

import (
	"fmt"

	"github.com/nabaz-io/nabaz/pkg/testrunner/diffengine"
	"github.com/nabaz-io/nabaz/pkg/testrunner/diffengine/parser"
	"github.com/nabaz-io/nabaz/pkg/testrunner/framework"
	"github.com/nabaz-io/nabaz/pkg/testrunner/models"
	"github.com/nabaz-io/nabaz/pkg/testrunner/scm/code"
	"github.com/nabaz-io/nabaz/pkg/testrunner/scm/history/git"
	"github.com/nabaz-io/nabaz/pkg/testrunner/storage"
)

type TestEngine struct {
	LocalCode      *code.CodeDirectory
	Storage        storage.Storage
	TestFramework  framework.Framework
	LanguageParser parser.Parser
	GitProvider    git.GitHistory
	CommitId       string
	LastNabazRun   *models.NabazRun
}

func LastNabazRunResult(currentCommitId string, storage storage.Storage, gitProvider git.GitHistory) *models.NabazRun {
	for currentCommitId != "" {
		nabazResult, err := storage.NabazRunByCommitID(currentCommitId)
		if err != nil {
			return nil
		}
		if nabazResult != nil {
			return nabazResult
		}

		commitParents, err := gitProvider.CommitParents(currentCommitId)
		if err != nil || len(commitParents) != 1 {
			return nil
		}
		currentCommitId = commitParents[0]
	}
	return nil
}

func NewTestEngine(localCode *code.CodeDirectory, storage storage.Storage, testFramework framework.Framework,
	languageParser parser.Parser, gitProvider git.GitHistory, commitId string) *TestEngine {

	lastNabazResult := LastNabazRunResult(commitId, storage, gitProvider)
	return &TestEngine{
		LocalCode:      localCode,
		Storage:        storage,
		TestFramework:  testFramework,
		LanguageParser: languageParser,
		GitProvider:    gitProvider,
		CommitId:       commitId,
		LastNabazRun:   lastNabazResult,
	}
}

func (t *TestEngine) listTests() []string {
	tetsMap := t.TestFramework.ListTests()
	return mapKeys(tetsMap)
}

func mapKeys(m map[string]string) []string {
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	return keys
}

func removeDuplications(s []string) []string {
    result := []string{}
    seen := map[string]bool{}
    for _, val := range s {
        if _, ok := seen[val]; !ok {
            result = append(result, val)
            seen[val] = true
        }
    }
    return result
}

func (t *TestEngine) FillTestCoverageFuncNames(testRuns []*models.TestRun) {
	for _, testRun := range testRuns {
		for _, scope := range testRun.CallGraph {
			fullFilePath := t.TestFramework.BasePath() + scope.Path

			code, err := t.LocalCode.GetFileContent(fullFilePath)
			if err != nil {
				panic(fmt.Errorf("failed to get file " + fullFilePath + err.Error()))
			}
			
			// TODO: optimize - same code may contain multiple functions, why parse it everytime?
			funcName, err := t.LanguageParser.FindFunction(code, &scope)
			if err != nil {
				panic(fmt.Errorf("failed to find function name for " + string(code) + err.Error()))
			}
			scope.FuncName = funcName
				
		}

		testRun.CallGraph = removeCallGraphDups(testRun.CallGraph)
	}
}

func removeCallGraphDups(s []code.Scope) []code.Scope {
	result := []code.Scope{}
	seen := make(map[string]bool)
	for _, val := range s {
		if _, ok := seen[val.FuncName]; !ok {
			result = append(result, val)
			seen[val.FuncName] = true
		}
	}
	return result
}

func (t *TestEngine) TestsToSkip() map[string]*models.TestRun {
	if t.LastNabazRun != nil {
		tests := t.listTests()
		diffEngine := diffengine.NewDiffEngine(t.LocalCode, t.GitProvider, t.LanguageParser, t.LastNabazRun.CommitID)
		testsToSkip := t.decideWhichTestsToSkip(tests, diffEngine)
		return testsToSkip
	}

	return map[string]*models.TestRun{}
}

func (engine *TestEngine) decideWhichTestsToSkip(tests []string, diffengine *diffengine.DiffEngine) map[string]*models.TestRun {
	testsToSkip := map[string]*models.TestRun{}

	codeDiff, err := engine.GitProvider.Diff(engine.CommitId, engine.LastNabazRun.CommitID)
	if err != nil {
		panic(err)
	}

	changedFunctions, err := diffengine.ChangedFunctions(codeDiff)
    uniqueChangedFunctions := removeDuplications(changedFunctions)
	if err != nil {
		panic(err)
	}

	for _, test := range tests {

		skippedTest := engine.LastNabazRun.PreviousTestRun(test)
		ranTest := engine.LastNabazRun.GetTestRun(test)

		//  if test is not in last nabaz run (as skipped or ran) we should stop searching and just run it, it's new
		if skippedTest == nil && ranTest == nil {
			continue
		}

		if skippedTest != nil {
			// test skipped in last run, should the NabazRun where it ran.
			relevantNabazResult, err := engine.Storage.NabazRunByRunID(skippedTest.RunIDReference)
			if err != nil {
				// NabazRun where it ran is not found, we should run it.
				continue
			}
			ranTest = relevantNabazResult.GetTestRun(test)
		}

		var scopes []code.Scope = ranTest.CallGraph
		if ranTest.TestFuncScope != (code.Scope{}) {
			scopes = append(scopes, ranTest.TestFuncScope)
		}

		if ranTest.Success == false || diffengine.Affects(uniqueChangedFunctions, scopes) {
			continue
		} else {
			testsToSkip[ranTest.Name] = ranTest
		}
	}

	return testsToSkip
}
