package token_introspection

import (
	"errors"
	"log"

	"baconi.co.uk/oauth/internal/pkg/token"
)

type Introspector interface {
	introspect(request) (response, error)
}

func NewIntrospector(authenticator token.Authenticator[token.AccessToken]) Introspector {
	return &introspector{authenticator}
}

type introspector struct {
	authenticator token.Authenticator[token.AccessToken]
}

// assert introspector implements Introspector
var _ Introspector = (*introspector)(nil)

func (service introspector) introspect(r request) (response, error) {

	accessToken, err := service.authenticator.Authenticate(r.token)

	switch {

	case errors.Is(err, token.ErrNoSuchToken):
		log.Println("[TRACE][token_introspection.Introspector] No such token")
		return response{Active: false}, nil

	case errors.Is(err, token.ErrTokenHasExpired):
		log.Println("[TRACE][token_introspection.Introspector] Token has expired")
		return response{Active: false}, nil

	case errors.Is(err, token.ErrTokenIsBefore):
		log.Println("[TRACE][token_introspection.Introspector] Token is not yet valid")
		return response{Active: false}, nil

	case err != nil:
		log.Println("[ERROR][token_introspection.Introspector] Failed to authenticate token:", err)
		return response{}, err

	// TODO - Decide out if we want to block any Confident client from introspecting any token.

	default:
		log.Println("[TRACE][token_introspection.Introspector] Token found and is valid")
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
