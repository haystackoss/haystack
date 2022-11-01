package reporter_test

import (
	"testing"

	"github.com/nabaz-io/nabaz/pkg/hypertest/models"
	"github.com/nabaz-io/nabaz/pkg/hypertest/reporter"
)

func TestSendResultTelemetry(t *testing.T) {
	telemetry := models.ResultTelemetry{
		HashedId:        "test",
		RunDuration:     0,
		Os:              "linux",
		Arch:            "amd64",
		LongestDuration: 0,
		TestsSkipped:    0,
		TestsRan:        0,
		TestsFailed:     0,
	}

	err := reporter.SendAnnonymousUsage(&telemetry)
	if err != nil {
		t.Errorf("failed to send1 telemetry")
	}

}

func TestSendExecutionTelemetry(t *testing.T) {
	err := reporter.SendAnnonymousStarted()
	if err != nil {
		t.Errorf("failed to send telemetry")
	}
}
