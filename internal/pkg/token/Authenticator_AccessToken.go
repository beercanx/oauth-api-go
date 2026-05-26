package token

import (
	"fmt"

	"baconi.co.uk/oauth/internal/pkg/db"
	"github.com/google/uuid"
)

type AccessTokenAuthenticator struct {
	repository Repository[db.AccessToken]
}

func (service *AccessTokenAuthenticator) Authenticate(token uuid.UUID) (db.AccessToken, error) {

	accessToken, err := service.repository.FindById(token)
	switch {

	case err != nil:
		return db.AccessToken{}, fmt.Errorf("authenticate access token failed: %w", err)

	case accessToken.HasExpired():
		if err = service.repository.DeleteByRecord(accessToken); err != nil {
			return db.AccessToken{}, fmt.Errorf("delete expired access token failed: %w", err)
		}
		return db.AccessToken{}, ErrAccessTokenHasExpired

	case accessToken.IsBefore():
		return db.AccessToken{}, ErrAccessTokenIsBefore

	default:
		return accessToken, nil
	}
}

// assert AccessTokenAuthenticator implements Authenticator
var _ Authenticator[db.AccessToken] = (*AccessTokenAuthenticator)(nil)

func NewAccessTokenAuthenticator(repository Repository[db.AccessToken]) *AccessTokenAuthenticator {
	return &AccessTokenAuthenticator{
		repository: repository,
	}
}
