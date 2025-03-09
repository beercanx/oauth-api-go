package token

import (
	"baconi.co.uk/oauth/internal/pkg/client"
	"baconi.co.uk/oauth/internal/pkg/scope"
	"baconi.co.uk/oauth/internal/pkg/user"
	"errors"
	"fmt"
	"github.com/google/uuid"
	"time"
)

type AccessTokenService struct {
	repository     Repository[AccessToken]
	tokenAge       time.Duration
	notBeforeShift time.Duration
}

// assert InMemoryRepository implements Issuer and Authenticator
var _ Issuer[AccessToken] = &AccessTokenService{}
var _ Authenticator[AccessToken] = &AccessTokenService{}

func NewAccessTokenService(repository Repository[AccessToken]) *AccessTokenService {
	return &AccessTokenService{
		repository:     repository,
		notBeforeShift: 1 * time.Minute,
		tokenAge:       2 * time.Hour,
	}
}

func (service *AccessTokenService) Issue(
	username user.AuthenticatedUsername,
	clientId client.Id,
	scopes []scope.Scope,
) AccessToken {

	issuedAt := time.Now()

	expiresAt := issuedAt.Add(service.tokenAge)
	notBefore := issuedAt.Add(-service.notBeforeShift)

	accessToken := AccessToken{
		Value:     uuid.New(),
		Username:  username,
		ClientId:  clientId,
		Scopes:    scopes,
		IssuedAt:  issuedAt,
		ExpiresAt: expiresAt,
		NotBefore: notBefore,
	}

	if err := service.repository.Insert(accessToken); err != nil {
		panic(fmt.Errorf("issue access token failed: %w", err))
	}

	return accessToken
}

var (
	ErrAccessTokenHasExpired = errors.New("access token has expired")
	ErrAccessTokenIsBefore   = errors.New("access token is before")
)

func (service *AccessTokenService) Authenticate(token uuid.UUID) (AccessToken, error) {

	accessToken, err := service.repository.FindById(token)
	switch {

	case err != nil:
		return AccessToken{}, fmt.Errorf("authenticate access token failed: %w", err)

	case accessToken.HasExpired():
		if err = service.repository.DeleteByRecord(accessToken); err != nil {
			return AccessToken{}, fmt.Errorf("delete expired access token failed: %w", err)
		}
		return AccessToken{}, ErrAccessTokenHasExpired

	case accessToken.IsBefore():
		return AccessToken{}, ErrAccessTokenIsBefore

	default:
		return accessToken, nil
	}
}
