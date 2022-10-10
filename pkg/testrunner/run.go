package testrunner

import (
	"errors"
	"hash/fnv"
	"log"
	"net/http"
	_ "net/http/pprof"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	parserfactory "github.com/nabaz-io/nabaz/pkg/testrunner/diffengine/parser/factory"
	frameworkfactory "github.com/nabaz-io/nabaz/pkg/testrunner/framework"
	"github.com/nabaz-io/nabaz/pkg/testrunner/models"
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

func hashString(s string) string {
	algorithm := fnv.New32a()
	algorithm.Write([]byte(s))
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
func Run(cmdline string, pkgs string, repoPath string) (*models.NabazRun, int) {
	go func() {
		log.Println(http.ListenAndServe("localhost:6060", nil))
	}()

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

	testsToSkip := testEngine.TestsToSkip()

	if len(testsToSkip) > 0 {
		log.Print("Nabaz skipping ", len(testsToSkip), " tests")
	}

	testResults, exitCode := framework.RunTests(testsToSkip)

	log.Printf("Ran %d/%d tests\n", len(testResults), len(testResults)+len(testsToSkip))

	testEngine.FillTestCoverageFuncNames(testResults)

	totalDuration := time.Since(startTime).Seconds()
	nabazRun := reporter.CreateNabazRun(testsToSkip, totalDuration, testEngine, history, testResults)
	storage.SaveNabazRun(nabazRun)

	hashedRepoName := hashString("TODO")
	annonymousTelemetry := reporter.NewAnnonymousTelemetry(nabazRun, hashedRepoName)
	reporter.SendAnonymousTelemetry(annonymousTelemetry)

	log.Printf("Total duration: %2f\n", totalDuration)
	log.Printf("longest duration %2f\n", nabazRun.LongestDuration)
	if totalDuration < nabazRun.LongestDuration {
		log.Printf("Saved %2f seconds of your life!\n", nabazRun.LongestDuration-totalDuration)
	}

	return nabazRun, exitCode

}

func Execute(args *Arguements) int {
	_, exitCode := Run(args.Cmdline, args.Pkgs, args.RepoPath)
	return exitCode
}
