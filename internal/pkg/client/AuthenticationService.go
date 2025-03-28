package client

import "github.com/alexedwards/argon2id"

type AuthenticationService struct {
	secretRepository    SecretRepository
	principalRepository PrincipalRepository
}

func (a AuthenticationService) AuthenticateAsPublic(clientId string) (Principal, bool) {
	principal, ok := a.principalRepository.FindByClientId(clientId)
	switch {
	case !ok:
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
	for _, s := range secrets {
		match, err := argon2id.ComparePasswordAndHash(clientSecret, s.hashedSecret)
		switch {
		case err != nil:
			continue
		case match:
			secret = s
			matched = true
			break
		}
	}

	if !matched {
		return Principal{}, false
	}

	principal, ok := a.principalRepository.FindById(secret.clientId)
	switch {
	case !ok:
		return principal, false
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
