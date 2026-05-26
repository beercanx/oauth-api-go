package token

import (
	"errors"

	"baconi.co.uk/oauth/internal/pkg/client"
	"baconi.co.uk/oauth/internal/pkg/user"
	"github.com/google/uuid"
)

var (
	ErrNoSuchToken = errors.New("token does not exist")
)

type Repository[T any] interface {
	Insert(new T) error
	FindById(id uuid.UUID) (T, error)
	FindAllByUsername(username user.AuthenticatedUsername) ([]T, error)
	FindAllByClientId(clientId client.Id) ([]T, error)
	DeleteById(id uuid.UUID) error
	DeleteByRecord(record T) error
	DeletedExpired() error
}
