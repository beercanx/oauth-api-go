package token

import (
	"baconi.co.uk/oauth/internal/pkg/client"
	"baconi.co.uk/oauth/internal/pkg/user"
	"github.com/google/uuid"
)

type InMemoryRepository[T Token] struct {
	store map[uuid.UUID]T
}

func (i *InMemoryRepository[T]) Insert(new T) error {
	i.store[new.GetValue()] = new
	return nil
}

func (i *InMemoryRepository[T]) FindById(id uuid.UUID) (T, error) {
	value, ok := i.store[id]
	if ok {
		return value, nil
	} else {
		return *new(T), ErrNoSuchToken
	}
}

func (i *InMemoryRepository[T]) FindAllByUsername(username user.AuthenticatedUsername) ([]T, error) {
	v := make([]T, 0, len(i.store))
	for _, value := range i.store {
		if value.GetUsername() == username {
			v = append(v, value)
		}
	}
	return v, nil
}

func (i *InMemoryRepository[T]) FindAllByClientId(clientId client.Id) ([]T, error) {
	v := make([]T, 0, len(i.store))
	for _, value := range i.store {
		if value.GetClientId() == clientId {
			v = append(v, value)
		}
	}
	return v, nil
}

func (i *InMemoryRepository[T]) DeleteById(id uuid.UUID) error {
	delete(i.store, id)
	return nil
}

func (i *InMemoryRepository[T]) DeleteByRecord(record T) error {
	return i.DeleteById(record.GetValue())
}

func (i *InMemoryRepository[T]) DeletedExpired() error {
	for _, value := range i.store {
		if value.HasExpired() {
			if err := i.DeleteByRecord(value); err != nil {
				return err
			}
		}
	}
	return nil
}

// assert InMemoryRepository implements Repository
var _ Repository[AccessToken] = &InMemoryRepository[AccessToken]{}
var _ Repository[RefreshToken] = &InMemoryRepository[RefreshToken]{}

func NewInMemoryRepository[T Token]() *InMemoryRepository[T] {
	return &InMemoryRepository[T]{store: make(map[uuid.UUID]T)}
}
