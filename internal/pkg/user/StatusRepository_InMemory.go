package user

import (
	"strings"
)

type inMemoryStatusRepository struct {
	store map[string]Status
}

func (repository *inMemoryStatusRepository) Insert(status Status) error {
	repository.store[strings.ToLower(status.username)] = status
	return nil
}

func (repository *inMemoryStatusRepository) FindByUsername(username string) (Status, error) {
	if status, ok := repository.store[strings.ToLower(username)]; ok {
		return status, nil
	}
	return Status{}, ErrNoSuchStatus
}

// assert inMemoryStatusRepository implements StatusRepository
var _ StatusRepository = (*inMemoryStatusRepository)(nil)

func NewInMemoryStatusRepository() StatusRepository {
	repository := &inMemoryStatusRepository{make(map[string]Status)}

	// TODO - Remove once we've got a means of creating new users
	_ = repository.Insert(Status{"aardvark", false})

	return repository
}
