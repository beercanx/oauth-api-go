package token_introspection

import (
	"errors"

	"baconi.co.uk/oauth/internal/pkg/token"
)

type Introspector interface {
	introspect(request) (response, error)
}

func NewIntrospector(accessTokenRepository token.Repository[token.AccessToken]) Introspector {
	return &introspector{accessTokenRepository: accessTokenRepository}
}

type introspector struct {
	accessTokenRepository token.Repository[token.AccessToken]
}

// assert introspector implements Introspector
var _ Introspector = (*introspector)(nil)

func (service introspector) introspect(r request) (response, error) {

	accessToken, err := service.accessTokenRepository.FindById(r.token)

	switch {

	case err != nil && errors.Is(err, token.ErrNoSuchToken):
		return response{Active: false}, nil

	case err != nil:
		return response{}, err

	case token.HasExpired(accessToken):
		return response{Active: false}, nil

	case token.IsBefore(accessToken):
		return response{Active: false}, nil

	// TODO - Decide out if we want to block any Confident client from introspecting any token.

	default:
		return response{
			Active:         true,
			Scope:          accessToken.GetScopes(),
			Subject:        accessToken.GetUsername(),
			Username:       accessToken.GetUsername(),
			ClientId:       accessToken.GetClientId(),
			TokenType:      token.Bearer,
			IssuedAt:       accessToken.GetIssuedAt().Unix(),
			NotBefore:      accessToken.GetNotBefore().Unix(),
			ExpirationTime: accessToken.GetExpiresAt().Unix(),
		}, nil
	}
}
