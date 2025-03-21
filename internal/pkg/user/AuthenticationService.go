package user

import (
	"errors"
	"github.com/alexedwards/argon2id"
)

type AuthenticationService struct {
	credentialRepository CredentialRepository
	statusRepository     StatusRepository
}

func (service *AuthenticationService) Authenticate(username string, password string) (*Success, *Failure) {

	credential, err := service.credentialRepository.FindByUsername(username)
	switch {
	case errors.Is(err, ErrNoSuchCredential):
		return nil, &Failure{Missing} // TODO - This is bad because of time-based attacks.
	case err != nil:
		panic(err)
	}

	// TODO - How can this be made to check non existent hashes to reduce surface area of a time-based attack.
	match, err := argon2id.ComparePasswordAndHash(password, credential.hashedSecret)
	switch {
	case err != nil:
		panic(err)
	case !match:
		return nil, &Failure{Mismatched}
	}

	status, err := service.statusRepository.FindByUsername(username)
	switch {
	case errors.Is(err, ErrNoSuchStatus):
		return nil, &Failure{Missing}
	case err != nil:
		panic(err)
	case status.isLocked():
		return nil, &Failure{Locked}
	default:
		return &Success{AuthenticatedUsername{username}}, nil
	}
}

// assert AuthenticationService implements Authenticator
var _ Authenticator = &AuthenticationService{}

func NewAuthenticationService(credentialRepository CredentialRepository, statusRepository StatusRepository) *AuthenticationService {
	return &AuthenticationService{credentialRepository: credentialRepository, statusRepository: statusRepository}
}
