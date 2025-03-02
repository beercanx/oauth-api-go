package token

import (
	"github.com/google/uuid"
)

type Authenticator[T Token] interface {
	Authenticate(token uuid.UUID) (*T, error)
}
