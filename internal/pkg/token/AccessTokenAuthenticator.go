package token

import (
	"fmt"
	"log/slog"

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
			slog.Warn("Failed to manually delete expired access token", slog.Any("error", err))
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
