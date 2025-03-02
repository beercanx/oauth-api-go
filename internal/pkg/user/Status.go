package user

import "strings"

type Status struct {
	username string
	locked   bool
}

func (status Status) isLocked() bool {
	return status.locked
}

type StatusRepository interface {
	Insert(status Status) error
	FindByUsername(username string) (*Status, error)
}

type InMemoryStatusRepository struct {
	store map[string]*Status
}

func NewInMemoryStatusRepository() *InMemoryStatusRepository {
	repository := &InMemoryStatusRepository{make(map[string]*Status)}

	// TODO - Remove once we've got a means of creating new users
	_ = repository.Insert(Status{"aardvark", false})

	return repository
}

func (repository InMemoryStatusRepository) Insert(status Status) error {
	repository.store[strings.ToLower(status.username)] = &status
	return nil
}

func (repository InMemoryStatusRepository) FindByUsername(username string) (*Status, error) {
	return repository.store[strings.ToLower(username)], nil
}
