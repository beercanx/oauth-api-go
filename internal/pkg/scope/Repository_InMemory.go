package scope

type InMemoryRepository struct {
	store map[string]Scope
}

func (repository *InMemoryRepository) FindById(id string) (Scope, error) {
	scope, ok := repository.store[id]
	if ok {
		return scope, nil
	} else {
		return Scope{}, ErrNoSuchScope
	}
}

// assert InMemoryRepository implements Repository
var _ Repository = &InMemoryRepository{}

func NewInMemoryRepository() *InMemoryRepository {
	repository := &InMemoryRepository{make(map[string]Scope)}
	repository.store["basic"] = Scope{"basic"}
	return repository
}
