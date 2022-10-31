package storage

import "github.com/nabaz-io/nabaz/pkg/adhdtest/models"

type Storage interface {
	NabazRunByRunID(runID int64) (*models.NabazRun, error)
	NabazRunByCommitID(commitID string) (*models.NabazRun, error)
	SaveNabazRun(nabazRun *models.NabazRun) error
	Reset() error
}

func NewStorage() (Storage, error) {
	return NewLocalStorage()
}
