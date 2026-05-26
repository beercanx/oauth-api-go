package token

import (
	"fmt"
	"time"

	"baconi.co.uk/oauth/internal/pkg/client"
	"baconi.co.uk/oauth/internal/pkg/db"
	"baconi.co.uk/oauth/internal/pkg/scope"
	"baconi.co.uk/oauth/internal/pkg/user"
	"github.com/google/uuid"
)

type AccessTokenIssuer struct {
	repository     Repository[db.AccessToken]
	tokenAge       time.Duration
	notBeforeShift time.Duration
}

func (issuer *AccessTokenIssuer) Issue(
	username user.AuthenticatedUsername,
	clientId client.Id,
	scopes scope.Scopes,
) (db.AccessToken, error) {

	issuedAt := time.Now()

	expiresAt := issuedAt.Add(issuer.tokenAge)
	notBefore := issuedAt.Add(-issuer.notBeforeShift)

	accessToken := db.AccessToken{
		ID:        uuid.New(),
		Username:  username,
		ClientID:  clientId,
		Scopes:    scopes,
		IssuedAt:  issuedAt,
		ExpiresAt: expiresAt,
		NotBefore: notBefore,
	}

	if err := issuer.repository.Insert(accessToken); err != nil {
		return db.AccessToken{}, fmt.Errorf("issue access token failed: %w", err)
	}

	return accessToken, nil
}

// assert AccessTokenIssuer implements Issuer
var _ Issuer[db.AccessToken] = (*AccessTokenIssuer)(nil)

func NewAccessTokenIssuer(repository Repository[db.AccessToken]) *AccessTokenIssuer {
	return &AccessTokenIssuer{
		repository:     repository,
		notBeforeShift: 1 * time.Minute,
		tokenAge:       2 * time.Hour,
	}
}
