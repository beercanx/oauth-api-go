package authentication

type UserAuthenticationService struct {
	// TODO userCredentialRepository UserCredentialRepository
	// TODO userStatusRepository UserStatusRepository
	// TODO inject argon2 support
}

func (service UserAuthenticationService) Authenticate(username string, password string) (*Success, *Failure, error) {
	// TODO implement me
	return &Success{AuthenticatedUsername{username}}, nil, nil
}
