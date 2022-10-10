package storage

import "github.com/nabaz-io/nabaz/pkg/testrunner/models"

type Storage interface {
	NabazRunByRunID(runID uint64) (*models.NabazRun, error)
	NabazRunByCommitID(commitID string) (*models.NabazRun, error)
	SaveNabazRun(nabazRun *models.NabazRun) error
	Reset() error
}

func NewStorage() (Storage, error) {
	return NewLocalStorage()
}
