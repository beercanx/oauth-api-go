package client

import (
	"errors"
	"log"

	"github.com/alexedwards/argon2id"
)

// Authenticator TODO - Decide if we would continue returning bool, or switch to nil Principal and error.
type Authenticator interface {
	AuthenticateAsPublic(clientId string) (Principal, bool)
	AuthenticateAsConfidential(clientId string, clientSecret string) (Principal, bool)
}

type authenticator struct {
	secretRepository    SecretRepository
	principalRepository PrincipalRepository
}

var _ Authenticator = (*authenticator)(nil)

func NewAuthenticator(secretRepository SecretRepository, principalRepository PrincipalRepository) Authenticator {
	return &authenticator{secretRepository, principalRepository}
}

func (a authenticator) AuthenticateAsPublic(clientId string) (Principal, bool) {
	principal, err := a.principalRepository.FindByClientId(clientId)
	switch {
	case errors.Is(err, ErrNoSuchClientPrincipal):
		log.Printf("[DEBUG][client.Authenticator] No such public client: %s", clientId)
		return Principal{}, false
	case err != nil:
		log.Printf("[ERROR][client.Authenticator] Failed to retrieve public client %s: %v", clientId, err)
		return principal, false
	case !principal.IsPublic():
		log.Printf("[DEBUG][client.Authenticator] Client is not public: %s", clientId)
		return Principal{}, false
	default:
		log.Printf("[TRACE][client.Authenticator] Public client found for: %s", clientId)
		return principal, true
	}
}

func (a authenticator) AuthenticateAsConfidential(clientId string, clientSecret string) (Principal, bool) {

	secrets, _ := a.secretRepository.FindByClientId(clientId)

	var secret Secret
	var matched = false
loop:
	for _, s := range secrets {
		match, matchError := argon2id.ComparePasswordAndHash(clientSecret, s.hashedSecret)
		switch {
		case matchError != nil:
			log.Printf("[ERROR][client.Authenticator] Failed to compare client secret: %v", matchError)
			continue
		case match:
			log.Printf("[TRACE][client.Authenticator] Confidential client secret matched: %s - %s", clientId, s.id)
			secret = s
			matched = true
			break loop
		}
	}

	if !matched {
		log.Printf("[DEBUG][client.Authenticator] No confidential client secret matched: %s", clientId)
		return Principal{}, false
	}

	principal, repositoryError := a.principalRepository.FindById(secret.clientId)
	switch {
	case errors.Is(repositoryError, ErrNoSuchClientPrincipal):
		log.Printf("[DEBUG][client.Authenticator] No confidential client principal found: %s", secret.clientId)
		return Principal{}, false
	case repositoryError != nil:
		log.Printf("[ERROR][client.Authenticator] Failed to retrieve confidential client %s: %v", secret.clientId, repositoryError)
		return Principal{}, false
	case !principal.IsConfidential():
		log.Printf("[DEBUG][client.Authenticator] Client is not confidential: %s", secret.clientId)
		return Principal{}, false
	default:
		log.Printf("[TRACE][client.Authenticator] Confidential client authenticated for: %s", secret.clientId)
		return principal, true
	}
}
