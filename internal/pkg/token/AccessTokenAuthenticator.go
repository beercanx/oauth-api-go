package token

import (
	"fmt"
	"log"

	"github.com/google/uuid"
)

type accessTokenAuthenticator struct {
	repository RepositoryReadDelete[AccessToken]
}

func (service *accessTokenAuthenticator) Authenticate(token uuid.UUID) (AccessToken, error) {

	accessToken, err := service.repository.FindById(token)
	switch {

	case err != nil:
		return AccessToken{}, fmt.Errorf("authenticate access token failed: %w", err)

	case accessToken.HasExpired():
		if err = service.repository.DeleteById(accessToken.ID); err != nil {
			log.Println("[WARN] AccessTokenAuthenticator.Authenticate failed to delete expired access token:", err)
		}
		return AccessToken{}, ErrTokenHasExpired

	case accessToken.IsBefore():
		return AccessToken{}, ErrTokenIsBefore

	default:
		return accessToken, nil
	}
}

// assert accessTokenAuthenticator implements Authenticator
var _ Authenticator[AccessToken] = (*accessTokenAuthenticator)(nil)

func NewAccessTokenAuthenticator(repository RepositoryReadDelete[AccessToken]) Authenticator[AccessToken] {
	return &accessTokenAuthenticator{
		repository: repository,
	}
}
