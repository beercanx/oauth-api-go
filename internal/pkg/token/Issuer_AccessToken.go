package token

import (
	"baconi.co.uk/oauth/internal/pkg/client"
	"baconi.co.uk/oauth/internal/pkg/scope"
	"baconi.co.uk/oauth/internal/pkg/user"
	"fmt"
	"github.com/google/uuid"
	"time"
)

type AccessTokenIssuer struct {
	repository     Repository[AccessToken]
	tokenAge       time.Duration
	notBeforeShift time.Duration
}

func (issuer *AccessTokenIssuer) Issue(
	username user.AuthenticatedUsername,
	clientId client.Id,
	scopes scope.Scopes,
) AccessToken {

	issuedAt := time.Now()

	expiresAt := issuedAt.Add(issuer.tokenAge)
	notBefore := issuedAt.Add(-issuer.notBeforeShift)

	accessToken := AccessToken{
		Value:     uuid.New(),
		Username:  username,
		ClientId:  clientId,
		Scopes:    scopes,
		IssuedAt:  issuedAt,
		ExpiresAt: expiresAt,
		NotBefore: notBefore,
	}

	if err := issuer.repository.Insert(accessToken); err != nil {
		panic(fmt.Errorf("issue access token failed: %w", err))
	}

	return accessToken
}

// assert AccessTokenIssuer implements Issuer
var _ Issuer[AccessToken] = (*AccessTokenIssuer)(nil)

func NewAccessTokenIssuer(repository Repository[AccessToken]) *AccessTokenIssuer {
	return &AccessTokenIssuer{
		repository:     repository,
		notBeforeShift: 1 * time.Minute,
		tokenAge:       2 * time.Hour,
	}
}
