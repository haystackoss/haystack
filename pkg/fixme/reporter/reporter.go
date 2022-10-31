package reporter

import (
	"bytes"
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"math"
	"net/http"
	"os"
	"os/user"
	"runtime"
	"time"

	"github.com/nabaz-io/nabaz/pkg/fixme/models"
	"github.com/nabaz-io/nabaz/pkg/fixme/scm/history/git"
	"github.com/nabaz-io/nabaz/pkg/fixme/testengine"
)

func CreateNabazRun(testsToSkip map[string]models.SkippedTest, totalDuration float64, testEngine *testengine.TestEngine, history git.GitHistory, testResults []models.TestRun) *models.NabazRun {
	skippedTests := make([]models.SkippedTest, 0, len(testsToSkip))
	for _, v := range testsToSkip {
		skippedTests = append(skippedTests, v)
	}

	longestDuration := totalDuration
	if testEngine.LastNabazRun != nil {
		longestDuration = math.Max(totalDuration, testEngine.LastNabazRun.LongestDuration)
	}

	return &models.NabazRun{
		RunID:           time.Now().UnixNano(),
		CommitID:        history.HEAD(),
		TestsRan:        testResults,
		TestsSkipped:    skippedTests,
		RunDuration:     totalDuration,
		LongestDuration: longestDuration,
	}
}


func UniqueHash() string {
	id := ""

	username, err := user.Current()
	if err == nil {
		id +=  username.Username
	}

	hostname, err := os.Hostname()
	if err == nil {
		id += "@" + hostname
	}

	return md5String(id)
}


func md5String(s string) string {
	algorithm := md5.New()
	algorithm.Write([]byte(s))
	hash := algorithm.Sum(nil)
	return hex.EncodeToString(hash)
}

func NewAnnonymousTelemetry(nabazRun *models.NabazRun) models.ResultTelemetry {
	return models.ResultTelemetry{
		HashedId:        UniqueHash(),
		Os:              runtime.GOOS,
		Arch:            runtime.GOARCH,
		RunDuration:     nabazRun.RunDuration,
		LongestDuration: nabazRun.LongestDuration,
		TestsSkipped:    len(nabazRun.TestsSkipped),
		TestsRan:        len(nabazRun.TestsRan),
		TestsFailed:     len(nabazRun.FailedTests()),
	}
}

func SendAnnonymousStarted() error {
	t := models.ExecutionTelemtry{
		HashedId: UniqueHash(),
		Os:   runtime.GOOS,
		Arch: runtime.GOARCH,
	}

	return sendAnnonymousTelemetry("/started", t)
}


func SendAnnonymousUsage(usage *models.ResultTelemetry) error {
	return sendAnnonymousTelemetry("/usage", &usage)
}

func sendAnnonymousTelemetry(endpoint string, telemetry models.Telemetry) error {
	if os.Getenv("NO_TELEMETRY") == "1" {
		return nil
	}

	j, err := json.Marshal(telemetry)
	if err != nil {
		return err
	}

	res, err := http.Post("https://api.nabaz.io/stats" + endpoint, "application/json", bytes.NewBuffer(j))

	if err != nil {
		return err
	} else if res.StatusCode != 200 {
		return fmt.Errorf("bad status code %d", res.StatusCode)
	}

	return nil
}
