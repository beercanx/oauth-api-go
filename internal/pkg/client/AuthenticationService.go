package client

import (
	"log"

	"github.com/alexedwards/argon2id"
)

type AuthenticationService struct {
	secretRepository    SecretRepository
	principalRepository PrincipalRepository
}

func (a AuthenticationService) AuthenticateAsPublic(clientId string) (Principal, bool) {
	principal, err := a.principalRepository.FindByClientId(clientId)
	switch {
	case err != nil:
		log.Printf("Failed to retrieve public client %s: %v", clientId, err)
		return principal, false
	case !principal.IsPublic():
		return Principal{}, false
	default:
		return principal, true
	}
}

func (a AuthenticationService) AuthenticateAsConfidential(clientId string, clientSecret string) (Principal, bool) {

	secrets, _ := a.secretRepository.FindByClientId(clientId)

	var secret Secret
	var matched = false
loop:
	for _, s := range secrets {
		match, err := argon2id.ComparePasswordAndHash(clientSecret, s.hashedSecret)
		switch {
		case err != nil:
			continue
		case match:
			secret = s
			matched = true
			break loop
		}
	}

	if !matched {
		return Principal{}, false
	}

	principal, err := a.principalRepository.FindById(secret.clientId)
	switch {
	case err != nil:
		log.Printf("Failed to retrieve confidential client %s: %v", secret.clientId, err)
		return Principal{}, false
	case !principal.IsConfidential():
		return Principal{}, false
	default:
		return principal, true
	}
}

// assert AuthenticationService implements Authenticator
var _ Authenticator = (*AuthenticationService)(nil)

func NewAuthenticationService(secretRepository SecretRepository, principalRepository PrincipalRepository) *AuthenticationService {
	return &AuthenticationService{secretRepository, principalRepository}
}
