package token

import (
	"errors"
	"github.com/google/uuid"
)

var (
	ErrAccessTokenHasExpired = errors.New("access token has expired")
	ErrAccessTokenIsBefore   = errors.New("access token is before")
)

type Authenticator[T Token] interface {
	Authenticate(token uuid.UUID) (T, error)
}
