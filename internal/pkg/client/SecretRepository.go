package client

import (
	"github.com/google/uuid"
)

type SecretRepository interface {
	FindById(id uuid.UUID) (Secret, bool)
	FindByClient(client Id) ([]Secret, bool)
	FindByClientId(clientId string) ([]Secret, bool)
}
