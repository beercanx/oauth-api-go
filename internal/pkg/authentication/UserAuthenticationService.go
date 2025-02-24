package authentication

import "github.com/alexedwards/argon2id"

type UserAuthenticationService struct {
	userCredentialRepository UserCredentialRepository
	userStatusRepository     UserStatusRepository
}

func NewUserAuthenticationService(userCredentialRepository UserCredentialRepository, userStatusRepository UserStatusRepository) *UserAuthenticationService {
	return &UserAuthenticationService{userCredentialRepository: userCredentialRepository, userStatusRepository: userStatusRepository}
}

func (service UserAuthenticationService) Authenticate(username string, password string) (*Success, *Failure, error) {

	credential, err := service.userCredentialRepository.FindByUsername(username)
	switch {
	case err != nil:
		return nil, nil, err
	case credential == nil:
		return nil, &Failure{Missing}, nil // TODO - This is bad because of time-based attacks.
	}

	// TODO - How can this be made to check non existent hashes to reduce surface area of a time-based attack.
	match, err := argon2id.ComparePasswordAndHash(password, credential.hashedSecret)
	switch {
	case err != nil:
		return nil, nil, err
	case !match:
		return nil, &Failure{Mismatched}, nil
	}

	status, err := service.userStatusRepository.FindByUsername(username)
	switch {
	case err != nil:
		return nil, nil, err
	case status == nil:
		return nil, &Failure{Missing}, nil
	case status.isLocked():
		return nil, &Failure{Locked}, nil
	default:
		return &Success{AuthenticatedUsername{username}}, nil, nil
	}
}
