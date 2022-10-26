package reporter_test

import (
	"testing"

	"github.com/nabaz-io/nabaz/pkg/fixme/models"
	"github.com/nabaz-io/nabaz/pkg/fixme/reporter"
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
