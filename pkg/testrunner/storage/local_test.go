package storage_test

import (
	"testing"

	"github.com/nabaz-io/nabaz/pkg/testrunner/models"
	"github.com/nabaz-io/nabaz/pkg/testrunner/scm/code"
	"github.com/nabaz-io/nabaz/pkg/testrunner/storage"
)

func TestLocalStorage(t *testing.T) {
	s, err := storage.NewLocalStorage("./db.sqlite3")
	if err != nil {
		panic(err)
	}
	defer s.Reset()

	runID := uint64(1337)
	commitID := "abcdef1234567890"

	callGraph := make([]code.Scope, 0)
	err = s.SaveNabazRun(&models.NabazRun{
		RunID:    runID,
		CommitID: commitID,
		TestsRan: []models.TestRun{
			{Name: "Test1", Success: true, TimeInMs: 10, TestFuncScope: nil, CallGraph: callGraph},
		},
		TestsSkipped:    []models.SkippedTest{},
		RunDuration:     10,
		LongestDuration: 10,
	})
	if err != nil {
		t.Errorf("Failed to save NabazRun: %v", err)
	}

	runByRunID, err := s.NabazRunByRunID(runID)
	if err != nil {
		t.Errorf("Error getting NabazRun by RunID: %s", err)
	}

	if runByRunID.CommitID != commitID {
		t.Errorf("RunID and CommitID don't match")
	}

	if runByRunID.TestsRan[0].Name != "Test1" {
		t.Errorf("Test name doesn't match")
	}
}
