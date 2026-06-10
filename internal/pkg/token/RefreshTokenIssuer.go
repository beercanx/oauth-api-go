package token

import (
	"fmt"
	"time"

	"baconi.co.uk/oauth/internal/pkg/client"
	"baconi.co.uk/oauth/internal/pkg/scope"
	"baconi.co.uk/oauth/internal/pkg/user"
	"github.com/google/uuid"
)

type refreshTokenIssuer struct {
	repository     RepositoryCreate[RefreshToken]
	tokenAge       time.Duration
	notBeforeShift time.Duration
}

func (issuer *refreshTokenIssuer) Issue(
	username user.AuthenticatedUsername,
	clientId client.Id,
	scopes scope.Scopes,
) (RefreshToken, error) {

	issuedAt := time.Now()

	expiresAt := issuedAt.Add(issuer.tokenAge)
	notBefore := issuedAt.Add(-issuer.notBeforeShift)

	refreshToken := RefreshToken{
		ID:        uuid.New(),
		Username:  username,
		ClientId:  clientId,
		Scopes:    scopes,
		IssuedAt:  issuedAt,
		ExpiresAt: expiresAt,
		NotBefore: notBefore,
	}

	if err := issuer.repository.Insert(refreshToken); err != nil {
		return RefreshToken{}, fmt.Errorf("issue refresh token failed: %w", err)
	}

	return refreshToken, nil
}

// assert RefreshTokenService implements Issuer
var _ Issuer[RefreshToken] = (*refreshTokenIssuer)(nil)

func NewRefreshTokenIssuer(repository RepositoryCreate[RefreshToken]) Issuer[RefreshToken] {
	return &refreshTokenIssuer{
		repository:     repository,
		notBeforeShift: 1 * time.Minute,
		tokenAge:       4 * time.Hour,
	}
}
