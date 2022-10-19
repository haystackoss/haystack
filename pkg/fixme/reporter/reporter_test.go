package reporter_test

import (
	"testing"

	"github.com/nabaz-io/nabaz/pkg/fixme/models"
	"github.com/nabaz-io/nabaz/pkg/fixme/reporter"
)

func TestSendResultTelemetry(t *testing.T) {
	telemetry := models.ResultTelemetry{
		RepoName:        "test",
		RunDuration:     0.1,
		Os:              "linux",
		Arch:            "amd64",
		LongestDuration: 100,
		TestsSkipped:    100,
		TestsRan:        1,
		TestsFailed:     1,
	}

	err := reporter.SendAnonymousTelemetry(telemetry)
	if err != nil {
		t.Errorf("failed to send telemetry")
	}

}

func TestSendExecutionTelemetry(t *testing.T) {
	err := reporter.SendNabazStarted()
	if err != nil {
		t.Errorf("failed to send telemetry")
	}
}
