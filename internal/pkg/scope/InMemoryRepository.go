package scope

import "errors"

var (
	ErrNoSuchScope = errors.New("scope does not exist")
)

type InMemoryRepository struct {
	store map[string]Scope
}

// assert InMemoryRepository implements Repository
var _ Repository = &InMemoryRepository{}

func NewInMemoryRepository() *InMemoryRepository {
	repository := &InMemoryRepository{make(map[string]Scope)}
	repository.store["basic"] = Scope{"basic"}
	return repository
}

func (repository *InMemoryRepository) FindById(id string) (Scope, error) {
	scope, ok := repository.store[id]
	if ok {
		return scope, nil
	} else {
		return Scope{}, ErrNoSuchScope
	}
}
