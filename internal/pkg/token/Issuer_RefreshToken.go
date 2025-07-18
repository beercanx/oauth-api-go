package token

import (
	"baconi.co.uk/oauth/internal/pkg/client"
	"baconi.co.uk/oauth/internal/pkg/scope"
	"baconi.co.uk/oauth/internal/pkg/user"
	"fmt"
	"github.com/google/uuid"
	"time"
)

type RefreshTokenIssuer struct {
	repository     Repository[RefreshToken]
	tokenAge       time.Duration
	notBeforeShift time.Duration
}

func (issuer *RefreshTokenIssuer) Issue(
	username user.AuthenticatedUsername,
	clientId client.Id,
	scopes scope.Scopes,
) (RefreshToken, error) {

	issuedAt := time.Now()

	expiresAt := issuedAt.Add(issuer.tokenAge)
	notBefore := issuedAt.Add(-issuer.notBeforeShift)

	refreshToken := RefreshToken{
		value:     uuid.New(),
		username:  username,
		clientId:  clientId,
		scopes:    scopes,
		issuedAt:  issuedAt,
		expiresAt: expiresAt,
		notBefore: notBefore,
	}

	if err := issuer.repository.Insert(refreshToken); err != nil {
		return RefreshToken{}, fmt.Errorf("issue refresh token failed: %w", err)
	}

	return refreshToken, nil
}

// assert RefreshTokenService implements Issuer
var _ Issuer[RefreshToken] = (*RefreshTokenIssuer)(nil)

func NewRefreshTokenIssuer(repository Repository[RefreshToken]) *RefreshTokenIssuer {
	return &RefreshTokenIssuer{
		repository:     repository,
		notBeforeShift: 1 * time.Minute,
		tokenAge:       4 * time.Hour,
	}
}
