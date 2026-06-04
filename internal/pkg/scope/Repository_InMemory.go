package scope

type InMemoryRepository struct {
	store map[string]Scope
}

func (repository *InMemoryRepository) FindById(id string) (Scope, error) {
	scope, ok := repository.store[id]
	if ok {
		return scope, nil
	}
	return "", ErrNoSuchScope
}

// assert InMemoryRepository implements Repository
var _ Repository = (*InMemoryRepository)(nil)

func NewInMemoryRepository() *InMemoryRepository {
	repository := &InMemoryRepository{make(map[string]Scope)}
	repository.store["basic"] = "basic"
	repository.store["read"] = "read"
	repository.store["write"] = "write"
	return repository
}
