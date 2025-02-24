package token

import (
	"baconi.co.uk/oauth/internal/pkg/authentication"
	"baconi.co.uk/oauth/internal/pkg/client"
	"github.com/google/uuid"
)

type InMemoryRepository[T Token] struct {
	store map[uuid.UUID]T
}

func NewInMemoryRepository[T Token]() *InMemoryRepository[T] {
	return &InMemoryRepository[T]{store: make(map[uuid.UUID]T)}
}

func (i InMemoryRepository[T]) Insert(new T) (T, error) {
	i.store[new.GetValue()] = new
	return new, nil
}

func (i InMemoryRepository[T]) FindById(id uuid.UUID) (*T, error) {
	value, ok := i.store[id]
	if ok {
		return &value, nil
	} else {
		return nil, nil
	}
}

func (i InMemoryRepository[T]) FindAllByUsername(username authentication.AuthenticatedUsername) ([]T, error) {
	v := make([]T, 0, len(i.store))
	for _, value := range i.store {
		if value.GetUsername().Value == username.Value {
			v = append(v, value)
		}
	}
	return v, nil
}

func (i InMemoryRepository[T]) FindAllByClientId(clientId client.Id) ([]T, error) {
	v := make([]T, 0, len(i.store))
	for _, value := range i.store {
		if value.GetClientId().Value == clientId.Value {
			v = append(v, value)
		}
	}
	return v, nil
}

func (i InMemoryRepository[T]) DeleteById(id uuid.UUID) error {
	delete(i.store, id)
	return nil
}

func (i InMemoryRepository[T]) DeleteByRecord(record T) error {
	return i.DeleteById(record.GetValue())
}

func (i InMemoryRepository[T]) DeletedExpired() error {
	for _, value := range i.store {
		if value.HasExpired() {
			if err := i.DeleteByRecord(value); err != nil {
				return err
			}
		}
	}
	return nil
}
