package testrunner

import (
	"errors"
	"log"
	"os"
	"strings"
	"time"

	parserfactory "github.com/nabaz-io/nabaz/pkg/testrunner/diffengine/parser/factory"
	frameworkfactory "github.com/nabaz-io/nabaz/pkg/testrunner/framework"
	"github.com/nabaz-io/nabaz/pkg/testrunner/models"
	"github.com/nabaz-io/nabaz/pkg/testrunner/scm/code"
	historyfactory "github.com/nabaz-io/nabaz/pkg/testrunner/scm/history/git/factory"
	storagefactory "github.com/nabaz-io/nabaz/pkg/testrunner/storage/factory"
)

func getCwd() string {
	dir, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}
	return dir
}

func cd(path string) {
	os.Chdir(path)
}

func parseCmdline(cmdline string) (string, string, error) {

	/*
			def parse_cmdline(cmdline):
				supported_frameworks = ["pytest", "go test"]
		    	if framework := next((framework for framework in supported_frameworks if cmdline.startswith(framework)), None):
		        	args = cmdline.removeprefix(framework)
		        	return framework, args

			    raise click.UsageError(f"Unknown test framework provided, "
		        	                   f"test-runner currently supports {supported_frameworks} only.")
	*/

	supportedFrameworks := []string{"pytest", "go test"}
	for _, framework := range supportedFrameworks {
		if framework == cmdline {
			args := strings.TrimPrefix(cmdline, framework)
			return framework, args, nil
		}
	}

	return "", "", errors.New("Unknown test framework provided, nabaz currently supports " + strings.Join(supportedFrameworks, ", ") + " only.")
}

// Run exists mainly for testing purposes
func Run(cmdline string, pkgs string, repoPath string) (*models.NabazRun, int) {
	oldCwd := getCwd()
	cd(repoPath)
	defer cd(oldCwd)

	log.Println("Starting nabaz.io test-runner...")

	startTime := time.Now()

	localCode := code.NewCodeDirectory(repoPath)
	history, err := historyfactory.NewGitHistory(repoPath)
	if err != nil {
		log.Fatal(err)
	}

	frameworkStr, testArgs, err := parseCmdline(cmdline)
	if err != nil {
		log.Fatal(err)
	}

	parser, err := parserfactory.NewParser(frameworkStr)
	if err != nil {
		log.Fatal(err)
	}

	framework, err := frameworkfactory.NewFramework(parser, frameworkStr, repoPath, testArgs, pkgs)
	if err != nil {
		log.Fatal(err)
	}

	storage := storagefactory.NewStorage("")

	testEngine = NewTestEngine(localCode, framework, storage, parser, gitProvider, history, logger)

	testsToSkip := testEngine.TestsToSkip()

	log.Printf("Running tests %d tests (skipping %d tests)\n")
	testResults, exitCode := framework.RunTests()

	log.Printf("Ran %d/%d tests\n", len(testResults), len(testResults)+len(testsToSkip))

	testEngine.PopulateTestResultsWithMetadata(testResults)

	totalDuration := time.Now().Sub(startTime)

	return testResults, exitCode

}

func Execute(args *Arguements) int {
	exitCode, _ := Run(args.Cmdline, args.Pkgs, args.RepoPath)
	return exitCode
}
