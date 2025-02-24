package token

import (
	"baconi.co.uk/oauth/internal/pkg/authentication"
	"baconi.co.uk/oauth/internal/pkg/client"
	"github.com/google/uuid"
)

type Repository[T Token] interface {
	Insert(new T) (T, error)
	FindById(id uuid.UUID) (*T, error)
	FindAllByUsername(username authentication.AuthenticatedUsername) ([]T, error)
	FindAllByClientId(clientId client.Id) ([]T, error)
	DeleteById(id uuid.UUID) error
	DeleteByRecord(record T) error
	DeletedExpired() error
}
