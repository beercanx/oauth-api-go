package client

import (
	"errors"
	"log/slog"

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
		slog.Debug("No such public client", "clientId", clientId)
		return Principal{}, false
	case err != nil:
		slog.Error("Failed to retrieve public client", "clientId", clientId, "error", err)
		return principal, false
	case !principal.IsPublic():
		slog.Debug("Client is not public", "clientId", clientId)
		return Principal{}, false
	default:
		slog.Debug("Public client found", "clientId", clientId)
		return principal, true
	}
}

func (a authenticator) AuthenticateAsConfidential(clientId string, clientSecret string) (Principal, bool) {

	secrets, secretError := a.secretRepository.FindByClientId(clientId)
	if secretError != nil {
		slog.Error("Failed to retrieve confidential client secrets", "clientId", clientId, "error", secretError)
		return Principal{}, false
	}

	var secret Secret
	var matched = false
loop:
	for _, s := range secrets {
		match, matchError := argon2id.ComparePasswordAndHash(clientSecret, s.hashedSecret)
		switch {
		case matchError != nil:
			slog.Error("Failed to compare client secret", "error", matchError)
			continue
		case match:
			slog.Debug("Confidential client secret matched", "clientId", clientId, "secretId", s.id)
			secret = s
			matched = true
			break loop
		}
	}

	if !matched {
		slog.Debug("No confidential client secret matched", "clientId", clientId)
		return Principal{}, false
	}

	principal, principalError := a.principalRepository.FindById(secret.clientId)
	switch {
	case errors.Is(principalError, ErrNoSuchClientPrincipal):
		slog.Debug("No confidential client principal found", "clientId", secret.clientId)
		return Principal{}, false
	case principalError != nil:
		slog.Error("Failed to retrieve confidential client", "clientId", secret.clientId, "error", principalError)
		return Principal{}, false
	case !principal.IsConfidential():
		slog.Debug("Client is not confidential", "clientId", secret.clientId)
		return Principal{}, false
	default:
		slog.Debug("Confidential client authenticated", "clientId", secret.clientId)
		return principal, true
	}
}
