package token

import (
	"baconi.co.uk/oauth/internal/pkg/client"
	"baconi.co.uk/oauth/internal/pkg/user"
	"github.com/google/uuid"
)

type Repository[T Token] interface {
	Insert(new T) error                                                 // TODO - Verify if Go DB libraries panic or return errors
	FindById(id uuid.UUID) (T, error)                                   // TODO - Verify if Go DB libraries panic or return errors
	FindAllByUsername(username user.AuthenticatedUsername) ([]T, error) // TODO - Verify if Go DB libraries panic or return errors
	FindAllByClientId(clientId client.Id) ([]T, error)                  // TODO - Verify if Go DB libraries panic or return errors
	DeleteById(id uuid.UUID) error                                      // TODO - Verify if Go DB libraries panic or return errors
	DeleteByRecord(record T) error                                      // TODO - Verify if Go DB libraries panic or return errors
	DeletedExpired() error                                              // TODO - Verify if Go DB libraries panic or return errors
}
