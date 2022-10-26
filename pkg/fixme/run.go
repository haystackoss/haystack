package fixme

import (
	"errors"
	"fmt"
	"log"
	_ "net/http/pprof"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/fsnotify/fsnotify"
	junitparser "github.com/joshdk/go-junit"
	parserfactory "github.com/nabaz-io/nabaz/pkg/fixme/diffengine/parser/factory"
	frameworkfactory "github.com/nabaz-io/nabaz/pkg/fixme/framework"
	"github.com/nabaz-io/nabaz/pkg/fixme/models"
	"github.com/nabaz-io/nabaz/pkg/fixme/paths"
	"github.com/nabaz-io/nabaz/pkg/fixme/reporter"
	"github.com/nabaz-io/nabaz/pkg/fixme/scm/code"
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
func Run(cmdline string, repoPath string, outputChannel chan<- models.NabazOutput) {
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

	nabazOutput := models.NabazOutput{}

	nabazOutput.IsThinking = true
	outputChannel <- nabazOutput

	testsToSkip, _, err := testEngine.TestsToSkip()
	nabazOutput.IsThinking = false
	if err != nil {
		nabazOutput.Err = err.Error()
		outputChannel <- nabazOutput
		return
	}

	os.Remove(paths.JunitXMLPath())

	nabazOutput.IsRunningTests = true
	outputChannel <- nabazOutput
	testResults, _ := framework.RunTests(testsToSkip)

	xmlPath := paths.JunitXMLPath()
	suites, _ := junitparser.IngestFile(xmlPath)

	if len(testResults) == 0 {
		nabazOutput.IsRunningTests = false
		outputChannel <- nabazOutput
	}

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
				nabazOutput.FailedTests = append(nabazOutput.FailedTests, models.FailedTest{
					Name:     test.Name,
					FileLink: testNameToFileLink[test.Name],
					Err:      test.Error.Error(),
				})
			}
		}
	}

	nabazOutput.IsRunningTests = false
	outputChannel <- nabazOutput

	testEngine.FillTestCoverageFuncNames(testResults)

	totalDuration := time.Since(startTime).Seconds()
	nabazRun := reporter.CreateNabazRun(testsToSkip, totalDuration, testEngine, history, testResults)
	storage.SaveNabazRun(nabazRun)

	annonymousTelemetry := reporter.NewAnnonymousTelemetry(nabazRun)
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

func handleOutput(outputChannel <-chan models.NabazOutput) {
	Red := "\033[31m"
	Bold := "\033[1m"
	Reset := "\033[0m"
	Yellow := "\033[33m"

	outputState := models.OutputState{}
	outputState.FailedTests = []models.FailedTest{}

	for {
		select {
		case newOutput := <-outputChannel:
			maxLines := getTerminalHeight()
			if outputState.PreviousTestsFailedOutput == "" {
				fmt.Print("\033[H\033[2J")
			}

			if newOutput.IsThinking || newOutput.IsRunningTests {
				if newOutput.IsThinking {
					fmt.Println("🧠 thinking...")
				} else {
					fmt.Println("🚀 running tests...")
				}
				continue
			}

			if newOutput.Err != "" {
				if outputState.PreviousTestsFailedOutput != "" {
					fmt.Print("\033[H\033[2J")
					buildFailedO := fmt.Sprintf("🛠️  Fix build:\n%s\n", string(newOutput.Err))
					buildFailedlineAmount := len(strings.Split(buildFailedO, "\n"))

					remainingLines := maxLines - buildFailedlineAmount
					relevantLines := strings.Split(outputState.PreviousTestsFailedOutput, "\n")[0:remainingLines]
					buildFailedO += strings.Join(relevantLines, "\n")

					fmt.Print(buildFailedO)
				} else {
					buildOutput := fmt.Sprintf("\n🛠️  Fix build:\n%s\n", string(newOutput.Err))
					splitted := strings.Split(buildOutput, "\n")
					fmt.Print(strings.Join(splitted[0:maxLines], "\n"))
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

				// update / remove tests that failed before
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

				output := fmt.Sprintf("\n🛠️  %sFix tests:%s\n\n", Bold, Reset)
				
				for index, failedTest := range outputState.FailedTests {
					
					testOutput := fmt.Sprintf("  ❌ %s%s%s ", Red, failedTest.Name, Reset)

					testFileExtension := frameworkfactory.TestFileExtension(failedTest.Err)
					if testFileExtension == "" && failedTest.FileLink != "" {
						testOutput += fmt.Sprintf(" (%s%s%s)", Yellow, failedTest.FileLink, Reset) // add file link to output
					}

					testOutput += "\n"

					fileLineSuffix := fmt.Sprintf(".%s:", testFileExtension)
					if failedTest.Err != "Failed" {
						errLines := strings.Split(failedTest.Err, "\n")
						for _, errLine := range errLines {
							if testFileExtension != "" && strings.Contains(errLine, fileLineSuffix) {
								splitted := strings.SplitN(errLine, ":", 3) // x_test.go:123: error message
								fileName := splitted[0]
								lineNumber := splitted[1]
								errorMessage := splitted[2]
								testOutput += fmt.Sprintf("     %s%s:%s%s:%s\n", Yellow, fileName, lineNumber, Reset, errorMessage)

							} else {
								testOutput += fmt.Sprintf("     %s\n", errLine)
							}

						}
						testOutput += fmt.Sprintln()
					}

					if len(strings.Split(testOutput, "\n")) + len(strings.Split(output, "\n")) >= maxLines {
						output += fmt.Sprintf("  %d hidden... (too large, expand terminal or do your TODOs)\n", len(outputState.FailedTests)- index)
						break
					} else {
						output += testOutput
					}
				}

				fmt.Println(output)
				outputState.PreviousTestsFailedOutput = output
			}

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

func runNabazWhenNeeded(cmdline string, repoPath string, pleaseRunChannel <-chan time.Time, outputChannel chan<- models.NabazOutput) {
	previousRunRequestedTime := time.Unix(0, 0)
	previousRunStartedTime := time.Unix(0, 0)

	for {
		select {
		case runRequestTime := <-pleaseRunChannel:
			if runRequestTime.Sub(previousRunRequestedTime) < 250*time.Millisecond {
				// IDEs are making many syscalls, so we need to wait a bit before running
				continue
			}

			// if previous run started after this request, dont run
			if previousRunStartedTime.After(runRequestTime) {
				continue
			}

			previousRunStartedTime = time.Now().UTC()
			Run(cmdline, repoPath, outputChannel)
			previousRunRequestedTime = runRequestTime

		}
	}
}

func Execute(args *Arguements) error {
	reporter.SendNabazStarted()

	absRepoPath, err := filepath.Abs(args.RepoPath)
	if err != nil {
		return err
	}

	oldCwd := getCwd()
	cd(absRepoPath)
	defer cd(oldCwd)

	outputChannel := make(chan models.NabazOutput, 1000)
	go handleOutput(outputChannel)

	pleaseRunChannel := make(chan time.Time, 1000)
	go runNabazWhenNeeded(args.Cmdline, absRepoPath, pleaseRunChannel, outputChannel)

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
