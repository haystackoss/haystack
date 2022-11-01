package storage_test

import (
	"os"
	"testing"

	"github.com/nabaz-io/nabaz/pkg/hypertest/models"
	"github.com/nabaz-io/nabaz/pkg/hypertest/scm/code"
	"github.com/nabaz-io/nabaz/pkg/hypertest/storage"
)

func TestLocalStorage(t *testing.T) {
	s, err := storage.NewLocalStorage()
	if err != nil {
		panic(err)
	}
	defer os.Remove(os.TempDir() + "/nabaz.db")

	runID := int64(1337)
	commitID := "abcdef1234567890"
	scope := code.Scope{Path: "/tmp", FuncName: "nabaz", StartLine: 42, StartCol: 42, EndLine: 42, EndCol: 42}
	callGraph := []*code.Scope{&scope}
	err = s.SaveNabazRun(&models.NabazRun{
		RunID:    runID,
		CommitID: commitID,
		TestsRan: []models.TestRun{
			{Name: "Test1", Success: true, TimeInMs: 10, TestFuncScope: &scope, CallGraph: callGraph},
		},
		TestsSkipped: []models.SkippedTest{
			{Name: "Test2", RunIDRef: 12345},
		},
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

	if runByRunID.TestsRan[0].TestFuncScope.Path != "/tmp" {
		t.Errorf("TestFuncScope doesn't match")
	}

	if runByRunID.TestsRan[0].CallGraph[0].Path != "/tmp" {
		t.Errorf("CallGraph doesn't match")
	}

	_, err = s.NabazRunByCommitID(commitID)
	if err != nil {
		t.Errorf("Error getting NabazRun by CommitID: %s", err)
	}

	_ = s.SaveNabazRun(&models.NabazRun{
		RunID:    runID + 1,
		CommitID: commitID,
		TestsRan: []models.TestRun{
			{Name: "Test1", Success: true, TimeInMs: 10, TestFuncScope: &scope, CallGraph: callGraph},
		},
		TestsSkipped: []models.SkippedTest{
			{Name: "Test2", RunIDRef: 12345},
		},
		RunDuration:     10,
		LongestDuration: 10,
	})
	runByRunID, err = s.NabazRunByRunID(runID)
	if err != nil {
		t.Errorf("Error getting NabazRun by RunID: %s", err)
	}

	runByRunID, err = s.NabazRunByRunID(runID + 1)
	if err != nil {
		t.Errorf("Error getting NabazRun by RunID: %s", err)
	}

}
