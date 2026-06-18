package user

import (
	"errors"
	"fmt"
	"log"

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
		log.Printf("[TRACE][user.Authenticator] No such credential for %s", username)
		return Authenticated{}, AuthenticationFailure{Missing}
	case credentialError != nil:
		log.Printf("[ERROR][user.Authenticator] Failed to find user credential for %s: %v", username, credentialError)
		return Authenticated{}, credentialError
	}

	match, matchError := argon2id.ComparePasswordAndHash(password, credential.hashedSecret)
	switch {
	case matchError != nil:
		log.Printf("[ERROR][user.Authenticator] Failed to compare password and hash for %s: %v", username, matchError)
		return Authenticated{}, matchError
	case !match:
		log.Printf("[DEBUG][user.Authenticator] Password mismatched for %s", username)
		return Authenticated{}, AuthenticationFailure{Mismatched}
	}

	status, statusError := service.statusRepository.FindByUsername(username)
	switch {
	case errors.Is(statusError, ErrNoSuchStatus):
		log.Printf("[WARN][user.Authenticator] No such status for %s", username)
		return Authenticated{}, AuthenticationFailure{Missing}
	case statusError != nil:
		log.Printf("[ERROR][user.Authenticator] Failed to find user status for %s: %v", username, statusError)
		return Authenticated{}, statusError
	case status.locked:
		log.Printf("[WARN][user.Authenticator] User is locked for %s", username)
		return Authenticated{}, AuthenticationFailure{Locked}
	default:
		log.Printf("[TRACE][user.Authenticator] User authenticated for %s", username)
		return Authenticated{AuthenticatedUsername(username)}, nil
	}
}
