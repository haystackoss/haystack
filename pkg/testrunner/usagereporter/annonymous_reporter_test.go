package usagereporter_test

import (
	"testing"

	"github.com/nabaz-io/nabaz/pkg/testrunner/models"
	"github.com/nabaz-io/nabaz/pkg/testrunner/usagereporter"
)

func TestSendTelemetry(t *testing.T) {
	telemetry := models.Telemetry{
		RepoName: "test",
		RunDuration: 0.1,
		LongestDuration: 100,
		SkippedTests: 100,
		RanTests: 1,
		FailedTests: 0,
	}

	err := usagereporter.SendTelemetry(telemetry)
	if err != nil {
		t.Errorf("failed to send telemetry")
	}

}
