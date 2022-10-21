package fixme

import (
	"errors"
	"fmt"
	"hash/fnv"
	"log"
	_ "net/http/pprof"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/briandowns/spinner"
	"github.com/fsnotify/fsnotify"
	junitparser "github.com/joshdk/go-junit"
	parserfactory "github.com/nabaz-io/nabaz/pkg/fixme/diffengine/parser/factory"
	frameworkfactory "github.com/nabaz-io/nabaz/pkg/fixme/framework"
	"github.com/nabaz-io/nabaz/pkg/fixme/paths"
	"github.com/nabaz-io/nabaz/pkg/fixme/reporter"
	"github.com/nabaz-io/nabaz/pkg/fixme/scm/code"
	historyfactory "github.com/nabaz-io/nabaz/pkg/fixme/scm/history/git/factory"
	"github.com/nabaz-io/nabaz/pkg/fixme/storage"
	"github.com/nabaz-io/nabaz/pkg/fixme/testengine"
	"github.com/nabaz-io/nabaz/pkg/fixme/watcher"
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
func Run(cmdline string, repoPath string) {
	nabazSpinner := spinner.New(spinner.CharSets[9], 100*time.Millisecond) // Build our new spinner
	nabazSpinner.Start()
	nabazSpinner.Prefix = "thinking..."

	reporter.SendNabazStarted()

	repoPath, err := filepath.Abs(repoPath)
	if err != nil {
		log.Fatal(err)
	}

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

	framework, err := frameworkfactory.NewFramework(parser, frameworkStr, repoPath, testArgs)
	if err != nil {
		log.Fatal(err)
	}

	storage, err := storage.NewStorage()
	if err != nil {
		panic(err)
	}
	testEngine := testengine.NewTestEngine(localCode, storage, framework, parser, history)

	testsToSkip, _, err := testEngine.TestsToSkip()
	if err != nil {
		nabazSpinner.Disable()
		fmt.Println(err.Error())
		return
	}

	nabazSpinner.Prefix = "running tests..."

	os.Remove(paths.JunitXMLPath())

	testResults, _ := framework.RunTests(testsToSkip)
	nabazSpinner.Disable()

	xmlPath := paths.JunitXMLPath()
	suites, _ := junitparser.IngestFile(xmlPath)

	if len(testResults) == 0 {
		fmt.Println("✅ all good.")
		return
	}

	Red := "\033[31m"
	// Yellow := "\033[33m"
	// Underline := "\033[4m"
	Bold := "\033[1m"
	Reset := "\033[0m"
	firstTest := true
	for _, suite := range suites {
		if suite.Totals.Failed == 0 {
			continue
		}
		// fmt.Printf("📦 %s%s%s\n", Red, suite.Name, Reset)
		for _, test := range suite.Tests {
			if test.Status == "failed" {
				if firstTest {
					fmt.Printf("\n🛠️  %sTODO%s\n\n", Bold, Reset)
					firstTest = false
				}

				fmt.Printf("  ❌ %s%s%s\n", Red, test.Name, Reset)

				testErr := test.Error.Error()
				if testErr != "Failed" {
					errLines := strings.Split(testErr, "\n")
					for _, errLine := range errLines {
						fmt.Printf("    %s\n", errLine)
					}
					fmt.Println()
				}

			}
		}
	}

	testEngine.FillTestCoverageFuncNames(testResults)

	totalDuration := time.Since(startTime).Seconds()
	nabazRun := reporter.CreateNabazRun(testsToSkip, totalDuration, testEngine, history, testResults)
	storage.SaveNabazRun(nabazRun)

	hashedRepoName := hashString("TODO")
	annonymousTelemetry := reporter.NewAnnonymousTelemetry(nabazRun, hashedRepoName)
	reporter.SendAnonymousTelemetry(annonymousTelemetry)

}

// handleFSCreate assumes that the file was just created and not already watched.
func handleFSCreate(w *watcher.Watcher, event fsnotify.Event) {

	info, err := os.Lstat(event.Name)
	if err != nil {
		fmt.Printf("error: %v\n", err)
		return
	}

	if info.IsDir() {
		w.WatchFolder(event.Name)
	}
}

func handleFSEvent(w *watcher.Watcher, cmdline string, repoPath string, event fsnotify.Event) {
	//TODO: Move this to something nicer.
	// do something
	switch event.Op {
	case fsnotify.Create:
		handleFSCreate(w, event)

	default:
		fmt.Printf("Rerunning tests")
		Run(cmdline, repoPath)
	}
}

func Execute(args *Arguements) error {
	absRepoPath, err := filepath.Abs(args.RepoPath)
	if err != nil {
		return err
	}

	oldCwd := getCwd()
	cd(absRepoPath)
	defer cd(oldCwd)

	Run(args.Cmdline, absRepoPath)
	w := watcher.NewWatcher(absRepoPath)
	for {
		select {
		case event := <-w.FileSystemEvents:
			fmt.Print("\033[H\033[2J")
			handleFSEvent(w, args.Cmdline, absRepoPath, event)
		case err := <-w.Errors:
			fmt.Printf("error: %v\n", err)
		}
	}

	return nil
}
