package user

import (
	"strings"
)

type InMemoryStatusRepository struct {
	store map[string]Status
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

// assert InMemoryStatusRepository implements StatusRepository
var _ StatusRepository = &InMemoryStatusRepository{}

func NewInMemoryStatusRepository() *InMemoryStatusRepository {
	repository := &InMemoryStatusRepository{make(map[string]Status)}

	// TODO - Remove once we've got a means of creating new users
	_ = repository.Insert(Status{"aardvark", false})

	return repository
}
