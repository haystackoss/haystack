package fixme

import (
	"errors"
	"fmt"
	"log"
	"math"
	_ "net/http/pprof"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/fsnotify/fsnotify"
	junitparser "github.com/joshdk/go-junit"
	parserfactory "github.com/nabaz-io/nabaz/pkg/fixme/diffengine/parser/factory"
	frameworkfactory "github.com/nabaz-io/nabaz/pkg/fixme/framework"
	"github.com/nabaz-io/nabaz/pkg/fixme/limit"
	"github.com/nabaz-io/nabaz/pkg/fixme/models"
	"github.com/nabaz-io/nabaz/pkg/fixme/paths"
	"github.com/nabaz-io/nabaz/pkg/fixme/reporter"
	"github.com/nabaz-io/nabaz/pkg/fixme/scm/code"
	"github.com/nabaz-io/nabaz/pkg/fixme/scm/history/git"
	historyfactory "github.com/nabaz-io/nabaz/pkg/fixme/scm/history/git/factory"
	"github.com/nabaz-io/nabaz/pkg/fixme/storage"
	"github.com/nabaz-io/nabaz/pkg/fixme/testengine"
	"github.com/nabaz-io/nabaz/pkg/fixme/watcher"
	"golang.org/x/term"
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
func Run(cmdline string, repoPath string, history git.GitHistory, outputChannel chan<- models.NabazOutput) {
	startTime := time.Now()

	localCode := code.NewCodeDirectory(repoPath)

	err := history.SaveAllFiles()
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

	nabazOutput := models.NabazOutput{}

	nabazOutput.IsThinking = true
	outputChannel <- nabazOutput

	testsToSkip, listedTestsAmount, err := testEngine.TestsToSkip()
	nabazOutput.IsThinking = false
	if err != nil {
		nabazOutput.Err = err.Error()
		outputChannel <- nabazOutput
		return
	}

	os.Remove(paths.JunitXMLPath())

	if listedTestsAmount == -1 { // didn't list tests => first run
		nabazOutput.IsFirstRun = true
	}

	nabazOutput.IsRunningTests = true
	outputChannel <- nabazOutput
	testResults, _ := framework.RunTests(testsToSkip)

	if len(testResults) > 0 {
		xmlPath := paths.JunitXMLPath()
		suites, _ := junitparser.IngestFile(xmlPath)

		countFailed := 0
		for _, suite := range suites {
			countFailed += suite.Totals.Failed
		}

		testNameToFileLink := frameworkfactory.TestNameToFileLink(frameworkStr, testResults)

		nabazOutput.FailedTests = []models.FailedTest{}
		for _, suite := range suites {
			if suite.Totals.Failed == 0 {
				continue
			}
			for _, test := range suite.Tests {
				if test.Status == "failed" {
					cleanTestName := test.Name
					if frameworkStr == "go test" { // handle sub-test
						cleanTestName = strings.Split(test.Name, "/")[0]
					}

					nabazOutput.FailedTests = append(nabazOutput.FailedTests, models.FailedTest{
						Name:     test.Name,
						FileLink: testNameToFileLink[cleanTestName],
						Err:      test.Error.Error(),
					})
				}
			}
		}
	}

	nabazOutput.IsFirstRun = false
	nabazOutput.IsRunningTests = false
	outputChannel <- nabazOutput

	testEngine.FillTestCoverageFuncNames(testResults)

	totalDuration := time.Since(startTime).Seconds()
	nabazRun := reporter.CreateNabazRun(testsToSkip, totalDuration, testEngine, history, testResults)
	storage.SaveNabazRun(nabazRun)

	annonymousTelemetry := reporter.NewAnnonymousTelemetry(nabazRun)
	reporter.SendAnnonymousUsage(&annonymousTelemetry)
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

func handleFSEvent(w *watcher.Watcher, event fsnotify.Event, pleaseRunChannel chan<- time.Time) {
	//TODO: Move this to something nicer.
	// do something
	switch event.Op {
	case fsnotify.Create:
		handleFSCreate(w, event)

	default:
		pleaseRunChannel <- time.Now().UTC()
	}
}

func FindFailedTest(failedTest string, list []models.FailedTest) *models.FailedTest {
	for _, test := range list {
		if test.Name == failedTest {
			return &test
		}
	}
	return nil
}

func HighestSliceIndex(slice []string, max int) int {
	return int(math.Min(float64(len(slice)), float64(max)))
}

func handleOutput(outputChannel <-chan models.NabazOutput) {
	Red := "\033[31m"
	Bold := "\033[1m"
	Reset := "\033[0m"
	Yellow := "\033[33m"

	outputState := models.OutputState{}
	outputState.FailedTests = []models.FailedTest{}

	for newOutput := range outputChannel {
		maxLines := getTerminalHeight()
		if outputState.PreviousTestsFailedOutput == "" {
			fmt.Print("\033[H\033[2J")
		}

		if newOutput.IsThinking || newOutput.IsRunningTests {
			if newOutput.IsThinking {
				if outputState.PreviousTestsFailedOutput != "" {
					fmt.Print("\n")
				}
				fmt.Println("🧠 thinking...")
			} else {
				str := "🚀 running tests"
				if newOutput.IsFirstRun {
					str += " (first run, running the entire test suite)"
				}
				str += "..."
				fmt.Println(str)
			}
			continue
		}

		if newOutput.Err != "" {
			if outputState.PreviousTestsFailedOutput != "" {
				fmt.Print("\033[H\033[2J")
				buildFailedOutput := fmt.Sprintf("🛠️  %sFix build:%s\n%s\n", Bold, Reset, string(newOutput.Err))
				buildFailedlineAmount := len(strings.Split(buildFailedOutput, "\n")) - 1

				remainingLines := maxLines - buildFailedlineAmount
				splitted := strings.Split(outputState.PreviousTestsFailedOutput, "\n")
				relevantLines := splitted[:HighestSliceIndex(splitted, remainingLines)]
				buildFailedOutput += strings.Join(relevantLines, "\n")

				fmt.Print(buildFailedOutput)
			} else {
				buildOutput := fmt.Sprintf("\n🛠️  %sFix build:%s\n%s\n", Bold, Reset, string(newOutput.Err))
				splitted := strings.Split(buildOutput, "\n")
				fmt.Print(strings.Join(splitted[:HighestSliceIndex(splitted, maxLines)], "\n"))
			}
			continue
		} else if len(newOutput.FailedTests) == 0 {
			fmt.Print("\033[H\033[2J")
			fmt.Println("✔️ All tests passing 🌈")
			outputState.PreviousTestsFailedOutput = ""
			outputState.FailedTests = []models.FailedTest{}
			continue
		} else { // some tests failed

			fmt.Print("\033[H\033[2J")

			output := fmt.Sprintf("🧪 %sFix tests:%s\n\n", Bold, Reset)

			// remove rotton tests, update failed again messages
			for index, cachedFailedTest := range outputState.FailedTests {
				freshMatchingFailedTest := FindFailedTest(cachedFailedTest.Name, newOutput.FailedTests)

				if freshMatchingFailedTest == nil {
					outputState.RemoveRottonTest(index)
				} else if freshMatchingFailedTest.Err != cachedFailedTest.Err {
					outputState.UpdateFailedTestError(index, freshMatchingFailedTest.Err)
				}
			}

			// insert new failed tests
			for _, freshFailedTest := range newOutput.FailedTests {
				if FindFailedTest(freshFailedTest.Name, outputState.FailedTests) == nil {
					outputState.AddFailedTest(freshFailedTest)
				}
			}

			for index, failedTest := range outputState.FailedTests {

				testOutput := fmt.Sprintf("  ❌ %s%s%s ", Red, failedTest.Name, Reset)

				testFileExtensionFromError := frameworkfactory.TestFileExtensionFromError(failedTest.Err)
				if testFileExtensionFromError == "" && failedTest.FileLink != "" {
					testOutput += fmt.Sprintf(" (%s%s%s)", Yellow, failedTest.FileLink, Reset) // add file link to output
				}

				fileLineSuffix := fmt.Sprintf(".%s:", testFileExtensionFromError)
				if failedTest.Err != "Failed" {
					testOutput += "\n"
					errLines := strings.Split(failedTest.Err, "\n")
					for lineInex, errLine := range errLines {
						if testFileExtensionFromError != "" && strings.Contains(errLine, fileLineSuffix) {
							splitted := strings.SplitN(errLine, ":", 3) // x_test.go:123: error message
							fileName := splitted[0]
							lineNumber := splitted[1]
							errorMessage := splitted[2]
							testOutput += fmt.Sprintf("     %s%s:%s%s:%s", Yellow, fileName, lineNumber, Reset, errorMessage)

						} else {
							testOutput += fmt.Sprintf("     %s", errLine)
							if lineInex < len(errLines)-1 {
								testOutput += "\n"
							}
						}

					}
				}

				if index != len(outputState.FailedTests)-1 {
					testOutput += "\n"
				}

				summedTotalLines := len(strings.Split(testOutput, "\n")) + len(strings.Split(output, "\n")) - len(outputState.FailedTests)
				if summedTotalLines > maxLines {
					output += fmt.Sprintf("  %d hidden... (too large, expand terminal or do your TODOs)\n", len(outputState.FailedTests)-index)
					break
				}

				output += testOutput

			}

			fmt.Print(output)
			outputState.PreviousTestsFailedOutput = output
		}
	}
}

func getTerminalHeight() int {
	if !term.IsTerminal(0) {
		return -1
	}

	_, height, err := term.GetSize(0) // cross-platform terminal size
	if err != nil {
		return -1
	}

	return height
}

func detectTerminalSizeChange(pleaseRunChannel chan<- time.Time) {
	termHeight := getTerminalHeight()
	for {
		time.Sleep(100 * time.Millisecond)
		newHeight := getTerminalHeight()
		if newHeight != termHeight {
			termHeight = newHeight
			pleaseRunChannel <- time.Now().UTC()
		}
	}
}

func runNabazWhenNeeded(cmdline string, repoPath string, history git.GitHistory, pleaseRunChannel <-chan time.Time, outputChannel chan<- models.NabazOutput) {
	previousRunRequestedTime := time.Unix(0, 0)
	previousRunStartedTime := time.Unix(0, 0)

	for {
		select {
		case runRequestTime := <-pleaseRunChannel:
			if runRequestTime.Sub(previousRunRequestedTime) < 50*time.Millisecond {
				// IDEs are making many syscalls, so we need to wait a bit before running
				continue
			}

			// if previous run started after this request, dont run
			if previousRunStartedTime.After(runRequestTime) {
				continue
			}

			previousRunStartedTime = time.Now().UTC()
			Run(cmdline, repoPath, history, outputChannel)
			previousRunRequestedTime = runRequestTime

		}
	}
}

func Execute(args *Arguements) error {
	limit.InitLimit()
	reporter.SendAnnonymousStarted()

	absRepoPath, err := filepath.Abs(args.RepoPath)
	if err != nil {
		return err
	}

	oldCwd := getCwd()
	cd(absRepoPath)
	defer cd(oldCwd)

	repoPath, err := filepath.Abs(absRepoPath)
	if err != nil {
		log.Fatal(err)
	}
	history, err := historyfactory.NewGitHistory(repoPath)
	if err != nil {
		log.Fatal(err)
	}

	outputChannel := make(chan models.NabazOutput, 1000)
	go handleOutput(outputChannel)

	pleaseRunChannel := make(chan time.Time, 1000)
	go runNabazWhenNeeded(args.Cmdline, absRepoPath, history, pleaseRunChannel, outputChannel)

	go detectTerminalSizeChange(pleaseRunChannel)

	pleaseRunChannel <- time.Now().UTC()

	w := watcher.NewWatcher(absRepoPath)
	for {
		select {
		case event := <-w.FileSystemEvents:
			handleFSEvent(w, event, pleaseRunChannel)
		case err := <-w.Errors:
			fmt.Printf("error: %v\n", err)
		}
	}
}
