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
) RefreshToken {

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
		panic(fmt.Errorf("issue refresh token failed: %w", err)) // TODO - Decide if this is an anti pattern in Go.
	}

	return refreshToken
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
