package token

import (
	"fmt"
	"time"

	"baconi.co.uk/oauth/internal/pkg/client"
	"baconi.co.uk/oauth/internal/pkg/scope"
	"baconi.co.uk/oauth/internal/pkg/user"
	"github.com/google/uuid"
)

type accessTokenIssuer struct {
	repository     RepositoryCreate[AccessToken]
	tokenAge       time.Duration
	notBeforeShift time.Duration
}

func (issuer *accessTokenIssuer) Issue(
	username user.AuthenticatedUsername,
	clientId client.Id,
	scopes scope.Scopes,
) (AccessToken, error) {

	issuedAt := time.Now()

	expiresAt := issuedAt.Add(issuer.tokenAge)
	notBefore := issuedAt.Add(-issuer.notBeforeShift)

	accessToken := AccessToken{
		ID:        uuid.New(),
		Username:  username,
		ClientID:  clientId,
		Scopes:    scopes,
		IssuedAt:  issuedAt,
		ExpiresAt: expiresAt,
		NotBefore: notBefore,
	}

	if err := issuer.repository.Insert(accessToken); err != nil {
		return AccessToken{}, fmt.Errorf("issue access token failed: %w", err)
	}

	return accessToken, nil
}

// assert accessTokenIssuer implements Issuer
var _ Issuer[AccessToken] = (*accessTokenIssuer)(nil)

func NewAccessTokenIssuer(repository RepositoryCreate[AccessToken]) Issuer[AccessToken] {
	return &accessTokenIssuer{
		repository:     repository,
		notBeforeShift: 1 * time.Minute,
		tokenAge:       2 * time.Hour,
	}
}
