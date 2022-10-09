package reporter_test

import (
	"testing"

	"github.com/nabaz-io/nabaz/pkg/testrunner/models"
	"github.com/nabaz-io/nabaz/pkg/testrunner/reporter"
)

func TestSendTelemetry(t *testing.T) {
	telemetry := models.Telemetry{
		RepoName: "test",
		RunDuration: 0.1,
		LongestDuration: 100,
		TestsSkipped: 100,
		TestsRan: 1,
		TestsFailed: 1,
	}

	err := reporter.SendAnonymousTelemetry(telemetry)
	if err != nil {
		t.Errorf("failed to send telemetry")
	}

}
