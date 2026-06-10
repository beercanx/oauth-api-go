package user

import (
	"errors"
	"fmt"

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

	credential, err := service.credentialRepository.FindByUsername(username)
	switch {
	case errors.Is(err, ErrNoSuchCredential):
		return Authenticated{}, AuthenticationFailure{Missing} // TODO - This is bad because of time-based attacks.
	case err != nil:
		return Authenticated{}, err
	}

	// TODO - How can this be made to check non existent hashes to reduce surface area of a time-based attack.
	match, err := argon2id.ComparePasswordAndHash(password, credential.hashedSecret)
	switch {
	case err != nil:
		return Authenticated{}, err
	case !match:
		return Authenticated{}, AuthenticationFailure{Mismatched}
	}

	status, err := service.statusRepository.FindByUsername(username)
	switch {
	case errors.Is(err, ErrNoSuchStatus):
		return Authenticated{}, AuthenticationFailure{Missing}
	case err != nil:
		return Authenticated{}, err
	case status.isLocked():
		return Authenticated{}, AuthenticationFailure{Locked}
	default:
		return Authenticated{AuthenticatedUsername(username)}, nil
	}
}
