package testengine

/*
from __future__ import annotations

from logging import Logger
from typing import List, Dict

from diff_engine.engine import DiffEngine
from diff_engine.parser.base import LanguageParser
from framework.base import TestFramework
from models.nabaz_result import NabazResult, CachedTestResult, TestResult
from scm.git_provider import GitProvider
from scm.local_code import LocalCodeDirectory
from storage.base import Storage


class TestEngine:
    def __init__(self, local_code: LocalCodeDirectory, storage: Storage, test_framework: TestFramework,
                 language_parser: LanguageParser, git_provider: GitProvider, commit_id: str, logger: Logger):
        self._local_code = local_code
        self._storage = storage
        self._test_framework = test_framework
        self._language_parser = language_parser
        self._git_provider = git_provider
        self._logger = logger
        self._commit_id = commit_id
        self.last_nabaz_result = self._get_last_nabaz_run_result(current_commit_id=self._commit_id)

    def get_tests_to_skip(self) -> Dict[str, CachedTestResult]:
        try:

            self._logger.info("Searching for the last nabaz run result...")
            if self.last_nabaz_result:
                tests = self._list_tests()
                self._logger.debug("listed tests: " + str(tests))
                self._logger.info("Found nabaz result {} for commit {}. thinking..."
                                  .format(self.last_nabaz_result.result_id, self.last_nabaz_result.commit_id))

                diff_engine = DiffEngine(local_code=self._local_code, language_parser=self._language_parser,
                                         old_commit_id=self.last_nabaz_result.commit_id,
                                         git_provider=self._git_provider, logger=self._logger)
                tests_to_skip = self._decide_which_tests_to_skip(tests, self.last_nabaz_result, diff_engine)

                self._logger.info(f"decided amount to skip: {len(tests_to_skip)} (out of {len(tests)})")
                return tests_to_skip
            else:
                self._logger.info("No previous nabaz run result found, running all tests.")
                return {}
        except Exception as e:
            self._logger.error("Failed to get tests to skip (critical).")
            raise e

    def populate_test_results_with_metadata(self, test_runs: List[TestResult]):
        for test_run in test_runs:
            for scope in test_run.call_graph:
                full_file_path = self._test_framework.base_path() + scope.path
                code = self._local_code.get_file(file_path=full_file_path)
                # TODO: optimize - same code may contain multiple functions, why parse it everytime?
                scope.func_name = self._language_parser.find_function(scope=scope, code=code)

            test_run.call_graph = [scope for scope in test_run.call_graph if scope.func_name]

    def _list_tests(self) -> List[str]:
        self._logger.info(f"Listing tests...")
        return self._test_framework.list_tests()

    def _get_last_nabaz_run_result(self, current_commit_id: str) -> NabazResult | None:
        try:
            while current_commit_id is not None:
                nabaz_result = self._storage.get_run_result(commit_id=current_commit_id)
                if nabaz_result:
                    return nabaz_result

                commit_parents = self._git_provider.get_commit_parents(commit_id=current_commit_id)
                if len(commit_parents) != 1:
                    return None  # >1 - We don't support multiple parents yet, git history is too complicated to handle.

                current_commit_id = commit_parents[0]
        except:
            self._logger.warning(f"Failed to get commit {current_commit_id} parents make sure that the provided"
                                 f" --commit-id exists in git.")

        return None

    def _decide_which_tests_to_skip(self, tests: List[str], last_nabaz_result: NabazResult,
                                    diff_engine: DiffEngine) -> Dict[str, CachedTestResult]:
        self._logger.info(f"deciding which tests to skip...")
        tests_to_skip: Dict[str, CachedTestResult] = {}

        self._logger.info(f"diffing between {self._commit_id} and {last_nabaz_result.commit_id}")

        try:
            code_diff = self._git_provider.diff(current_commit=self._commit_id,
                                                older_commit=last_nabaz_result.commit_id)
        except:
            raise Exception(f"Failed to diff. Make sure {self._commit_id} it's a valid commit id.")

        changed_functions = set(diff_engine.get_changed_function_in_files(code_diff))

        for test_name in tests:
            cached_test_result: CachedTestResult = last_nabaz_result.get_cached_test(test_name)
            if cached_test_result is None:
                # if test is not in last nabaz run test results, we should stop searching and just run it (it's new)
                continue

            if cached_test_result.ran:
                ran_test: TestResult = last_nabaz_result.get_test_run(test_name)
            else:
                relevant_nabaz_result = self._storage.get_run_result(result_id=cached_test_result.ran_result_id)
                ran_test: TestResult = relevant_nabaz_result.get_test_run(test_name)

            # TODO: add failed message, or error message caching so we can skip failed tests.
            tests_scopes = [*ran_test.call_graph]
            if ran_test.test_func_scope:
                tests_scopes.append(ran_test.test_func_scope)
            # should run test if diff affects its call graph or if test failed
            if ran_test.success is False or diff_engine.affects(changed_functions, tests_scopes, self._logger):
                continue
            else:
                # test not affected by diff, it should be skipped and let's use last result
                tests_to_skip[ran_test.name] = cached_test_result

        return tests_to_skip

*/
import (
	"github.com/nabaz-io/nabaz/pkg/testrunner/diffengine"
	"github.com/nabaz-io/nabaz/pkg/testrunner/diffengine/parser"
	"github.com/nabaz-io/nabaz/pkg/testrunner/framework"
	"github.com/nabaz-io/nabaz/pkg/testrunner/models"
	"github.com/nabaz-io/nabaz/pkg/testrunner/scm/code"
	"github.com/nabaz-io/nabaz/pkg/testrunner/scm/history/git"
)

type TestEngine struct {
	LocalCode       *code.CodeDirectory
	Storage         *Storage
	TestFramework   framework.Framework
	LanguageParser  parser.Parser
	GitProvider     git.GitHistory
	CommitId        string
	LastNabazResult *models.NabazRun
}

func LastNabazRunResult(currentCommitId string, storage *Storage, gitProvider git.GitHistory) *models.NabazRun {
	for currentCommitId != "" {
		nabazResult := storage.GetRunResult(currentCommitId)
		if nabazResult != nil {
			return nabazResult
		}

		commitParents, err := gitProvider.GetCommitParents(currentCommitId)
		if err != nil || len(commitParents) != 1 {
			return nil
		}
		currentCommitId = commitParents[0]
	}
	return nil
}

func NewTestEngine(localCode *code.CodeDirectory, storage *Storage, testFramework framework.Framework,
	 languageParser parser.Parser, gitProvider git.GitHistory, commitId string) *TestEngine {

	lastNabazResult := LastNabazRunResult(storage, commitId, gitProvider)
	return &TestEngine{
		LocalCode:       localCode,
		Storage:         storage,
		TestFramework:   testFramework,
		LanguageParser:  languageParser,
		GitProvider:     gitProvider,
		CommitId:        commitId,
		LastNabazResult: lastNabazResult,
	}
}


func (t *TestEngine) ListTests() []string {
	return t.TestFramework.ListTests()
}

func (t *TestEngine) GetTestsToSkip(tests []string) map[string]models.CachedTestResult {
	if t.LastNabazResult != nil {
		diffEngine := diffengine.NewDiffEngine( t.LocalCode, t.LanguageParser, t.LastNabazResult.CommitId, t.GitProvider)
		testsToSkip := t.DecideWhichTestsToSkip(tests, t.LastNabazResult, diffEngine)
		return testsToSkip
	}
	return map[string]models.CachedTestResult{}
}

func (engine *TestEngine) GetTestsToSkip() []string {
	if engine.LastNabazResult == nil {
		// TODO: add log
		tests := engine.TestFramework.ListTests()
		diffengine := diffengine.NewDiffEngine(engine.LocalCode, engine.GitProvider, engine.CommitId, engine.LanguageParser)
		testsToSkip := engine.decideWhichTestsToSkip(tests, diffengine)
		return testsToSkip
	}
	return []string{}
}

func (engine *TestEngine) decideWhichTestsToSkip(tests []string, diffengine *diffengine.DiffEngine) []string {
}