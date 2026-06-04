package token

import (
	"errors"

	"github.com/google/uuid"
)

var (
	ErrTokenHasExpired = errors.New("token has expired")
	ErrTokenIsBefore   = errors.New("token is before")
)

type Authenticator[T any] interface {
	Authenticate(token uuid.UUID) (T, error)
}
