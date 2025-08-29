package token_exchange

import (
	"errors"
	"math"
	"time"

	"baconi.co.uk/oauth/internal/pkg/token"
	"baconi.co.uk/oauth/internal/pkg/user"
)

type PasswordGrant struct {
	accessTokenIssuer  token.Issuer[token.AccessToken]
	refreshTokenIssuer token.Issuer[token.RefreshToken]
	userAuthenticator  user.Authenticator
}

func NewPasswordGrant(accessTokenIssuer token.Issuer[token.AccessToken], refreshTokenIssuer token.Issuer[token.RefreshToken], userAuthenticator user.Authenticator) *PasswordGrant {
	return &PasswordGrant{accessTokenIssuer, refreshTokenIssuer, userAuthenticator}
}

func (grant PasswordGrant) Exchange(request PasswordRequest) (Success, error) {

	success, err := grant.userAuthenticator.Authenticate(request.Username, request.Password)
	var failure user.AuthenticationFailure
	switch {
	case errors.As(err, &failure):
		return Success{}, Failed{Err: InvalidGrant, Description: string(failure.Reason)}
	case err != nil:
		return Success{}, err
	}

	accessToken, err := grant.accessTokenIssuer.Issue(success.Username, request.Principal.Id, request.Scopes)
	if err != nil {
		return Success{}, err
	}

	refreshToken, err := grant.refreshTokenIssuer.Issue(success.Username, request.Principal.Id, request.Scopes)
	if err != nil {
		return Success{}, err
	}

	return Success{
		AccessToken:  accessToken.GetValue(),
		TokenType:    token.Bearer,
		ExpiresIn:    secondsBetween(accessToken.GetExpiresAt(), accessToken.GetIssuedAt()),
		RefreshToken: refreshToken.GetValue(),
		Scope:        accessToken.GetScopes(),
		State:        request.State,
	}, nil
}

func secondsBetween(end time.Time, start time.Time) int64 {
	return int64(math.Round(end.Sub(start).Seconds()))
}
