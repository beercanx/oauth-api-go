package token

import (
	"fmt"

	"baconi.co.uk/oauth/internal/pkg/db"
	"github.com/google/uuid"
)

type accessTokenAuthenticator struct {
	repository RepositoryReadDelete[db.AccessToken]
}

func (service *accessTokenAuthenticator) Authenticate(token uuid.UUID) (db.AccessToken, error) {

	accessToken, err := service.repository.FindById(token)
	switch {

	case err != nil:
		return db.AccessToken{}, fmt.Errorf("authenticate access token failed: %w", err)

	case accessToken.HasExpired():
		if err = service.repository.DeleteByRecord(accessToken); err != nil {
			return db.AccessToken{}, fmt.Errorf("delete expired access token failed: %w", err)
		}
		return db.AccessToken{}, ErrTokenHasExpired

	case accessToken.IsBefore():
		return db.AccessToken{}, ErrTokenIsBefore

	default:
		return accessToken, nil
	}
}

// assert accessTokenAuthenticator implements Authenticator
var _ Authenticator[db.AccessToken] = (*accessTokenAuthenticator)(nil)

func NewAccessTokenAuthenticator(repository RepositoryReadDelete[db.AccessToken]) Authenticator[db.AccessToken] {
	return &accessTokenAuthenticator{
		repository: repository,
	}
}
