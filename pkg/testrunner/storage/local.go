package storage

import (
	"database/sql"
	"encoding/json"

	_ "github.com/mattn/go-sqlite3"
	"github.com/nabaz-io/nabaz/pkg/testrunner/models"
)

type LocalStorage struct {
	db              *sql.DB
	cacheByRunID    map[uint64]*models.NabazRun
	cacheByCommitID map[string]*models.NabazRun
}

func NewLocalStorage(path string) (*LocalStorage, error) {
	db, err := sql.Open("sqlite3", path)
	if err != nil {
		return nil, err
	}

	cacheByCommitID := make(map[string]*models.NabazRun)
	cacheByRunID := make(map[uint64]*models.NabazRun)

	err = createNabazRunTable(db)
	if err != nil {
		return nil, err
	}
	return &LocalStorage{db: db, cacheByRunID: cacheByRunID, cacheByCommitID: cacheByCommitID}, nil
}

func createNabazRunTable(db *sql.DB) error {
	_, err := db.Exec(`
		CREATE TABLE IF NOT EXISTS nabaz_runs (
			run_id INTEGER PRIMARY KEY,
			commit_id TEXT,
			tests_ran BLOB,
			tests_skipped BLOB,
			run_duration REAL,
			longest_duration REAL
		)
	`)
	return err
}
func (s *LocalStorage) NabazRunByRunID(runID uint64) (*models.NabazRun, error) {
	if cachedRun, ok := s.cacheByRunID[runID]; ok {
		return cachedRun, nil
	}

	row := s.db.QueryRow("SELECT * FROM nabaz_runs WHERE run_id = ?", runID)

	run := models.NabazRun{}

	unmarshaledTestsSkipped := make([]byte, 0)
	unmarshaledTestsRan := make([]byte, 0)
	err := row.Scan(&run.RunID, &run.CommitID, &unmarshaledTestsRan, &unmarshaledTestsSkipped, &run.RunDuration, &run.LongestDuration)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(unmarshaledTestsRan, &run.TestsRan)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(unmarshaledTestsSkipped, &run.TestsSkipped)
	if err != nil {
		return nil, err
	}

	s.cacheByRunID[runID] = &run

	return &run, nil
}

func (s *LocalStorage) NabazRunByCommitID(commitID string) (*models.NabazRun, error) {
	if cachedRun, ok := s.cacheByCommitID[commitID]; ok {
		return cachedRun, nil
	}

	row := s.db.QueryRow("SELECT * FROM nabaz_runs WHERE commit_id = ?", commitID)

	run := models.NabazRun{}
	err := row.Scan(&run.RunID, &run.CommitID, run.TestsRan, run.TestsSkipped, &run.RunDuration, &run.LongestDuration)
	if err != nil {
		return nil, err
	}

	s.cacheByCommitID[commitID] = &run

	return &run, nil
}

func (s *LocalStorage) SaveNabazRun(nabazRun *models.NabazRun) error {
	marshaledTestsRan, err := json.Marshal(nabazRun.TestsRan)
	if err != nil {
		return err
	}

	marshaledTestsSkipped, err := json.Marshal(nabazRun.TestsSkipped)
	if err != nil {
		return err
	}

	_, err = s.db.Exec(
		`INSERT INTO nabaz_runs (
			run_id,
			commit_id,
			tests_ran,
			tests_skipped,
			run_duration,
			longest_duration
		) VALUES (?, ?, ?, ?, ?, ?)`,
		nabazRun.RunID, nabazRun.CommitID, marshaledTestsRan, marshaledTestsSkipped, nabazRun.RunDuration, nabazRun.LongestDuration)
	return err
}

func (s *LocalStorage) Reset() error {
	_, err := s.db.Exec("DROP TABLE IF EXISTS nabaz_runs")
	return err
}

func (s *LocalStorage) Close() error {
	return s.db.Close()
}