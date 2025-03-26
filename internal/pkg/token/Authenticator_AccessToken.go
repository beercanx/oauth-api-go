package token

import (
	"fmt"
	"github.com/google/uuid"
)

type AccessTokenAuthenticator struct {
	repository Repository[AccessToken]
}

func (service *AccessTokenAuthenticator) Authenticate(token uuid.UUID) (AccessToken, error) {

	accessToken, err := service.repository.FindById(token)
	switch {

	case err != nil:
		return AccessToken{}, fmt.Errorf("authenticate access token failed: %w", err)

	case HasExpired(accessToken):
		if err = service.repository.DeleteByRecord(accessToken); err != nil {
			return AccessToken{}, fmt.Errorf("delete expired access token failed: %w", err)
		}
		return AccessToken{}, ErrAccessTokenHasExpired

	case IsBefore(accessToken):
		return AccessToken{}, ErrAccessTokenIsBefore

	default:
		return accessToken, nil
	}
}

// assert AccessTokenAuthenticator implements Authenticator
var _ Authenticator[AccessToken] = (*AccessTokenAuthenticator)(nil)

func NewAccessTokenAuthenticator(repository Repository[AccessToken]) *AccessTokenAuthenticator {
	return &AccessTokenAuthenticator{
		repository: repository,
	}
}
