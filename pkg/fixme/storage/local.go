package storage

import (
	"database/sql"
	"encoding/json"
	"os"

	_ "github.com/mattn/go-sqlite3"
	"github.com/nabaz-io/nabaz/pkg/fixme/models"
)

type LocalStorage struct {
	db              *sql.DB
	cacheByRunID    map[int64]*models.NabazRun
	cacheByCommitID map[string]*models.NabazRun
}

func NewLocalStorage() (*LocalStorage, error) {
	var err error
	tmpdir := os.TempDir()
	if tmpdir == "" {
		tmpdir, err = os.UserHomeDir()
		if err != nil {
			tmpdir = "."
		}
	}

	db, err := sql.Open("sqlite3", tmpdir+"/nabaz.db")
	if err != nil {
		return nil, err
	}

	cacheByCommitID := make(map[string]*models.NabazRun)
	cacheByRunID := make(map[int64]*models.NabazRun)

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
func (s *LocalStorage) NabazRunByRunID(runID int64) (*models.NabazRun, error) {
	if cachedRun, ok := s.cacheByRunID[runID]; ok {
		return cachedRun, nil
	}

	row := s.db.QueryRow("SELECT * FROM nabaz_runs WHERE run_id = ? ORDER BY run_id DESC", runID)

	run, err := parseNabazRun(row)
	if err != nil {
		return nil, err
	}

	s.cacheByRunID[runID] = run

	return run, nil
}

func (s *LocalStorage) NabazRunByCommitID(commitID string) (*models.NabazRun, error) {
	if cachedRun, ok := s.cacheByCommitID[commitID]; ok {
		return cachedRun, nil
	}

	row := s.db.QueryRow("SELECT * FROM nabaz_runs WHERE commit_id = ? ORDER BY run_id DESC", commitID)

	run, err := parseNabazRun(row)
	if err != nil {
		return nil, err
	}

	s.cacheByCommitID[commitID] = run

	return run, nil
}

func parseNabazRun(row *sql.Row) (*models.NabazRun, error) {
	run := models.NabazRun{}

	unmarshaledTestsSkipped := make([]byte, 0)
	unmarshaledTestsRan := make([]byte, 0)
	err := row.Scan(&run.RunID, &run.CommitID, &unmarshaledTestsRan, &unmarshaledTestsSkipped, &run.RunDuration, &run.LongestDuration)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
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
