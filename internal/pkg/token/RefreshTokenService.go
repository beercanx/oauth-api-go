package token

import (
	"baconi.co.uk/oauth/internal/pkg/authentication"
	"baconi.co.uk/oauth/internal/pkg/client"
	"github.com/google/uuid"
	"time"
)

type RefreshTokenService struct {
	repository     Repository[RefreshToken]
	tokenAge       time.Duration
	notBeforeShift time.Duration
}

func NewRefreshTokenService(repository Repository[RefreshToken]) *RefreshTokenService {
	return &RefreshTokenService{
		repository:     repository,
		notBeforeShift: 1 * time.Minute,
		tokenAge:       2 * time.Hour,
	}
}

func (service RefreshTokenService) Issue(
	username authentication.AuthenticatedUsername,
	clientId client.Id,
	scopes []string,
) (*RefreshToken, error) {

	issuedAt := time.Now()

	expiresAt := issuedAt.Add(service.tokenAge)
	notBefore := issuedAt.Add(-service.notBeforeShift)

	refreshToken, err := service.repository.Insert(RefreshToken{
		Value:     uuid.New(),
		Username:  username,
		ClientId:  clientId,
		Scopes:    scopes,
		IssuedAt:  issuedAt,
		ExpiresAt: expiresAt,
		NotBefore: notBefore,
	})

	return &refreshToken, err
}
