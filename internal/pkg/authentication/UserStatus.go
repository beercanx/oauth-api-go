package authentication

import "strings"

type UserStatus struct {
	username string
	locked   bool
}

func (status UserStatus) isLocked() bool {
	return status.locked
}

type UserStatusRepository interface {
	Insert(status UserStatus) error
	FindByUsername(username string) (*UserStatus, error)
}

type InMemoryUserStatusRepository struct {
	store map[string]*UserStatus
}

func NewInMemoryUserStatusRepository() *InMemoryUserStatusRepository {
	repository := &InMemoryUserStatusRepository{make(map[string]*UserStatus)}

	// TODO - Remove once we've got a means of creating new users
	_ = repository.Insert(UserStatus{"aardvark", false})

	return repository
}

func (repository InMemoryUserStatusRepository) Insert(status UserStatus) error {
	repository.store[strings.ToLower(status.username)] = &status
	return nil
}

func (repository InMemoryUserStatusRepository) FindByUsername(username string) (*UserStatus, error) {
	return repository.store[strings.ToLower(username)], nil
}
