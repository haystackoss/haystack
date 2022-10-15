package testrunner

import (
	"bufio"
	"errors"
	"hash/fnv"
	"log"
	_ "net/http/pprof"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/briandowns/spinner"
	parserfactory "github.com/nabaz-io/nabaz/pkg/testrunner/diffengine/parser/factory"
	frameworkfactory "github.com/nabaz-io/nabaz/pkg/testrunner/framework"
	"github.com/nabaz-io/nabaz/pkg/testrunner/reporter"
	"github.com/nabaz-io/nabaz/pkg/testrunner/scm/code"
	historyfactory "github.com/nabaz-io/nabaz/pkg/testrunner/scm/history/git/factory"
	"github.com/nabaz-io/nabaz/pkg/testrunner/storage"
	"github.com/nabaz-io/nabaz/pkg/testrunner/testengine"
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

func hashString(spinner string) string {
	algorithm := fnv.New32a()
	algorithm.Write([]byte(spinner))
	hash := algorithm.Sum32()
	return strconv.FormatUint(uint64(hash), 10)
}

func parseCmdline(cmdline string) (string, string, error) {
	supportedFrameworks := []string{"pytest", "go test"}
	for _, framework := range supportedFrameworks {
		if strings.HasPrefix(cmdline, framework) {
			args := strings.TrimPrefix(cmdline, framework)
			return framework, args, nil
		}
	}

	return "", "", errors.New("Unknown test framework provided, nabaz currently supports " + strings.Join(supportedFrameworks, ", ") + " only.")
}

// Run exists mainly for testing purposes
func Run(cmdline string, pkgs string, repoPath string) {
	spinner := spinner.New(spinner.CharSets[9], 100*time.Millisecond)
	go spinner.Start()

	reporter.SendNabazStarted()

	repoPath, err := filepath.Abs(repoPath)
	if err != nil {
		log.Fatal(err)
	}

	oldCwd := getCwd()
	cd(repoPath)
	defer cd(oldCwd)

	startTime := time.Now()

	localCode := code.NewCodeDirectory(repoPath)
	history, err := historyfactory.NewGitHistory(repoPath)
	if err != nil {
		log.Fatal(err)
	}

	err = history.SaveAllFiles()
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

	storage, err := storage.NewStorage()
	if err != nil {
		panic(err)
	}
	testEngine := testengine.NewTestEngine(localCode, storage, framework, parser, history)

	testsToSkip, testsAmount := testEngine.TestsToSkip()

	spinner.Stop()
	spinner.Disable()

	if testEngine.LastNabazRun != nil && testsAmount-len(testsToSkip) == 0 {
		log.Printf("No test were imapcted.")
	} else {
		testResults, _ := framework.RunTests(testsToSkip)
		log.Printf("Ran %d/%d tests.\n", len(testResults), len(testResults)+len(testsToSkip))

		testEngine.FillTestCoverageFuncNames(testResults)

		totalDuration := time.Since(startTime).Seconds()
		nabazRun := reporter.CreateNabazRun(testsToSkip, totalDuration, testEngine, history, testResults)
		storage.SaveNabazRun(nabazRun)

		hashedRepoName := hashString("TODO")
		annonymousTelemetry := reporter.NewAnnonymousTelemetry(nabazRun, hashedRepoName)
		reporter.SendAnonymousTelemetry(annonymousTelemetry)
	}


}

func Execute(args *Arguements) int {
	Run(args.Cmdline, args.Pkgs, args.RepoPath)
	return 0
}
