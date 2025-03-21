package token

import (
	"baconi.co.uk/oauth/internal/pkg/client"
	"baconi.co.uk/oauth/internal/pkg/scope"
	"baconi.co.uk/oauth/internal/pkg/user"
	"fmt"
	"github.com/google/uuid"
	"time"
)

type RefreshTokenService struct { // TODO - Consider refactoring into Issuer_RefreshToken
	repository     Repository[RefreshToken]
	tokenAge       time.Duration
	notBeforeShift time.Duration
}

func (service *RefreshTokenService) Issue(
	username user.AuthenticatedUsername,
	clientId client.Id,
	scopes []scope.Scope,
) RefreshToken {

	issuedAt := time.Now()

	expiresAt := issuedAt.Add(service.tokenAge)
	notBefore := issuedAt.Add(-service.notBeforeShift)

	refreshToken := RefreshToken{
		Value:     uuid.New(),
		Username:  username,
		ClientId:  clientId,
		Scopes:    scopes,
		IssuedAt:  issuedAt,
		ExpiresAt: expiresAt,
		NotBefore: notBefore,
	}

	if err := service.repository.Insert(refreshToken); err != nil {
		panic(fmt.Errorf("issue refresh token failed: %w", err))
	}

	return refreshToken
}

// assert RefreshTokenService implements Issuer
var _ Issuer[RefreshToken] = &RefreshTokenService{}

func NewRefreshTokenService(repository Repository[RefreshToken]) *RefreshTokenService {
	return &RefreshTokenService{
		repository:     repository,
		notBeforeShift: 1 * time.Minute,
		tokenAge:       4 * time.Hour,
	}
}
