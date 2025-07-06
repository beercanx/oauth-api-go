package user

import (
	"errors"
	"github.com/alexedwards/argon2id"
)

type AuthenticationService struct {
	credentialRepository CredentialRepository
	statusRepository     StatusRepository
}

func (service *AuthenticationService) Authenticate(username string, password string) (Authenticated, error) {

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
		return Authenticated{AuthenticatedUsername{username}}, nil
	}
}

// assert AuthenticationService implements Authenticator
var _ Authenticator = (*AuthenticationService)(nil)

func NewAuthenticationService(credentialRepository CredentialRepository, statusRepository StatusRepository) *AuthenticationService {
	return &AuthenticationService{credentialRepository: credentialRepository, statusRepository: statusRepository}
}
