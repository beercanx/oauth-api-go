package user

import (
	"errors"
	"strings"
)

type Status struct {
	username string
	locked   bool
}

func (status Status) isLocked() bool {
	return status.locked
}

type StatusRepository interface {
	Insert(status Status) error                     // TODO - Verify if Go DB libraries panic or return errors
	FindByUsername(username string) (Status, error) // TODO - Verify if Go DB libraries panic or return errors
}

var (
	ErrNoSuchStatus = errors.New("status does not exist")
)

type InMemoryStatusRepository struct {
	store map[string]Status
}

// assert InMemoryStatusRepository implements StatusRepository
var _ StatusRepository = &InMemoryStatusRepository{}

func NewInMemoryStatusRepository() *InMemoryStatusRepository {
	repository := &InMemoryStatusRepository{make(map[string]Status)}

	// TODO - Remove once we've got a means of creating new users
	_ = repository.Insert(Status{"aardvark", false})

	return repository
}

func (repository *InMemoryStatusRepository) Insert(status Status) error {
	repository.store[strings.ToLower(status.username)] = status
	return nil
}

func (repository *InMemoryStatusRepository) FindByUsername(username string) (Status, error) {
	if status, ok := repository.store[strings.ToLower(username)]; ok {
		return status, nil
	}
	return Status{}, ErrNoSuchStatus
}
