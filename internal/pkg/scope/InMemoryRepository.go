package scope

type InMemoryRepository struct {
	store map[string]Scope
}

func NewInMemoryRepository() *InMemoryRepository {
	repository := &InMemoryRepository{make(map[string]Scope)}
	repository.store["basic"] = Scope{"basic"}
	return repository
}

func (repository InMemoryRepository) FindById(id string) *Scope {
	scope, ok := repository.store[id]
	if ok {
		return &scope
	} else {
		return nil
	}
}
