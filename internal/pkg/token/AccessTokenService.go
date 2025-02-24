package token

import (
	"baconi.co.uk/oauth/internal/pkg/authentication"
	"baconi.co.uk/oauth/internal/pkg/client"
	"github.com/google/uuid"
	"time"
)

type AccessTokenService struct {
	repository     Repository[AccessToken]
	tokenAge       time.Duration
	notBeforeShift time.Duration
}

func NewAccessTokenService(repository Repository[AccessToken]) *AccessTokenService {
	return &AccessTokenService{
		repository:     repository,
		notBeforeShift: 1 * time.Minute,
		tokenAge:       2 * time.Hour,
	}
}

func (service AccessTokenService) Issue(
	username authentication.AuthenticatedUsername,
	clientId client.Id,
	scopes []string,
) (*AccessToken, error) {

	issuedAt := time.Now()

	expiresAt := issuedAt.Add(service.tokenAge)
	notBefore := issuedAt.Add(-service.notBeforeShift)

	accessToken, err := service.repository.Insert(AccessToken{
		Value:     uuid.New(),
		Username:  username,
		ClientId:  clientId,
		Scopes:    scopes,
		IssuedAt:  issuedAt,
		ExpiresAt: expiresAt,
		NotBefore: notBefore,
	})

	return &accessToken, err
}

func (service AccessTokenService) Authenticate(token uuid.UUID) (*AccessToken, error) {

	accessToken, err := service.repository.FindById(token)
	switch {

	case err != nil:
		return nil, err

	case accessToken == nil:
		return nil, nil

	case accessToken.HasExpired():
		return nil, service.repository.DeleteByRecord(*accessToken)

	case accessToken.IsBefore():
		return nil, nil

	default:
		return accessToken, nil
	}
}
