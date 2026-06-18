package user

import (
	"errors"
	"fmt"
	"log/slog"

	"github.com/alexedwards/argon2id"
)

type Authenticated struct {
	Username AuthenticatedUsername
}

type AuthenticationFailure struct {
	Reason AuthenticationFailureReason
}

func (f AuthenticationFailure) Error() string {
	return fmt.Sprintf("AuthenticationFailure: %s", f.Reason)
}

// assert AuthenticationFailure implements error
var _ error = (*AuthenticationFailure)(nil)

type AuthenticationFailureReason string

const (
	Missing    AuthenticationFailureReason = "missing"
	Mismatched AuthenticationFailureReason = "mismatched"
	Locked     AuthenticationFailureReason = "locked"
)

type Authenticator interface {
	Authenticate(username string, password string) (Authenticated, error)
}

type authenticator struct {
	credentialRepository CredentialRepository
	statusRepository     StatusRepository
}

// assert authenticator implements Authenticator
var _ Authenticator = (*authenticator)(nil)

func NewAuthenticator(credentialRepository CredentialRepository, statusRepository StatusRepository) Authenticator {
	return &authenticator{credentialRepository: credentialRepository, statusRepository: statusRepository}
}

func (service *authenticator) Authenticate(username string, password string) (Authenticated, error) {

	credential, credentialError := service.credentialRepository.FindByUsername(username)
	switch {
	case errors.Is(credentialError, ErrNoSuchCredential):
		slog.Debug("No such credential", "username", username)
		return Authenticated{}, AuthenticationFailure{Missing}
	case credentialError != nil:
		slog.Error("Failed to find user credential", "username", username, "error", credentialError)
		return Authenticated{}, credentialError
	}

	match, matchError := argon2id.ComparePasswordAndHash(password, credential.hashedSecret)
	switch {
	case matchError != nil:
		slog.Error("Failed to compare password and hash", "username", username, "error", matchError)
		return Authenticated{}, matchError
	case !match:
		slog.Debug("Password mismatched", "username", username)
		return Authenticated{}, AuthenticationFailure{Mismatched}
	}

	status, statusError := service.statusRepository.FindByUsername(username)
	switch {
	case errors.Is(statusError, ErrNoSuchStatus):
		slog.Warn("No such status", "username", username)
		return Authenticated{}, AuthenticationFailure{Missing}
	case statusError != nil:
		slog.Error("Failed to find user status", "username", username, "error", statusError)
		return Authenticated{}, statusError
	case status.locked:
		slog.Warn("User is locked", "username", username)
		return Authenticated{}, AuthenticationFailure{Locked}
	default:
		slog.Debug("User authenticated", "username", username)
		return Authenticated{AuthenticatedUsername(username)}, nil
	}
}
