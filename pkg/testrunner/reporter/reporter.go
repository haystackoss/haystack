package reporter

import (
	"bytes"
	"encoding/json"
	"fmt"
	"math"
	"math/rand"
	"net/http"
	"time"

	"github.com/nabaz-io/nabaz/pkg/testrunner/models"
	"github.com/nabaz-io/nabaz/pkg/testrunner/scm/history/git"
	"github.com/nabaz-io/nabaz/pkg/testrunner/testengine"
)

func CreateNabazRun(testsToSkip map[string]models.SkippedTest, totalDuration time.Duration, testEngine *testengine.TestEngine, history git.GitHistory, testResults []models.TestRun) *models.NabazRun {
	rndInt := rand.Uint64()
	skippedTests := make([]models.SkippedTest, 0, len(testsToSkip))
	for _, v := range testsToSkip {
		skippedTests = append(skippedTests, v)
	}

	longestDuration := totalDuration.Seconds()
	if testEngine.LastNabazRun != nil {
		longestDuration = math.Max(longestDuration, testEngine.LastNabazRun.LongestDuration)
	}

	return &models.NabazRun{
		RunID:           rndInt,
		CommitID:        history.HEAD(),
		TestsRan:        testResults,
		TestsSkipped:    skippedTests,
		RunDuration:     totalDuration.Seconds(),
		LongestDuration: longestDuration,
	}
}


func NewAnnonymousTelemetry(nabazRun *models.NabazRun, hashedRepoName string) models.Telemetry {
	return models.Telemetry{
		RepoName: hashedRepoName,
		RunDuration:   nabazRun.RunDuration,
		LongestDuration: nabazRun.LongestDuration,
		TestsSkipped: len(nabazRun.TestsSkipped),
		TestsRan:     len(nabazRun.TestsRan),
		TestsFailed: len(nabazRun.FailedTests()),
	}
}

func SendAnonymousTelemetry(telemetry models.Telemetry) error {
	j, err := json.Marshal(telemetry)
    if err != nil {
        fmt.Println("failed to marshal telemetry")
        return err
    }

	res, err := http.Post("https://api.nabaz.io/stats", "application/json", bytes.NewBuffer(j))
	if err != nil || res.StatusCode != 200 {
		fmt.Println("failed to send telemetry")
		return err
	}

	return nil
}