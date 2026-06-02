package token

import (
	"errors"

	"github.com/google/uuid"
)

var (
	ErrNoSuchToken = errors.New("token does not exist")
)

type Repository[T any] interface {
	Insert(new T) error
	FindById(id uuid.UUID) (T, error)
	DeleteById(id uuid.UUID) error
	DeleteByRecord(record T) error
	DeletedExpired() error
}
