package token_introspection

import (
	"errors"

	"baconi.co.uk/oauth/internal/pkg/db"
	"baconi.co.uk/oauth/internal/pkg/token"
)

type Introspector interface {
	introspect(request) (response, error)
}

func NewIntrospector(accessTokenRepository token.Repository[db.AccessToken]) Introspector {
	return &introspector{accessTokenRepository: accessTokenRepository}
}

type introspector struct {
	accessTokenRepository token.Repository[db.AccessToken]
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

	case accessToken.HasExpired():
		return response{Active: false}, nil

	case accessToken.IsBefore():
		return response{Active: false}, nil

	// TODO - Decide out if we want to block any Confident client from introspecting any token.

	default:
		return response{
			Active:         true,
			Scope:          accessToken.Scopes,
			Subject:        accessToken.Username,
			Username:       accessToken.Username,
			ClientId:       accessToken.ClientID,
			TokenType:      token.Bearer,
			IssuedAt:       accessToken.IssuedAt.Unix(),
			NotBefore:      accessToken.NotBefore.Unix(),
			ExpirationTime: accessToken.ExpiresAt.Unix(),
		}, nil
	}
}
