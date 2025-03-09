package token_exchange

import (
	"baconi.co.uk/oauth/internal/pkg/scope"
	"baconi.co.uk/oauth/internal/pkg/token"
	"baconi.co.uk/oauth/internal/pkg/user"
	"math"
	"time"
)

type PasswordGrant struct {
	accessTokenIssuer  token.Issuer[token.AccessToken]
	refreshTokenIssuer token.Issuer[token.RefreshToken]
	userAuthenticator  user.Authenticator
}

func NewPasswordGrant(accessTokenIssuer token.Issuer[token.AccessToken], refreshTokenIssuer token.Issuer[token.RefreshToken], userAuthenticator user.Authenticator) *PasswordGrant {
	return &PasswordGrant{accessTokenIssuer, refreshTokenIssuer, userAuthenticator}
}

func (grant PasswordGrant) Exchange(request PasswordRequest) (Response, error) {

	success, failure, err := grant.userAuthenticator.Authenticate(request.Username, request.Password)
	if err != nil {
		return nil, err
	}
	if failure != nil {
		return Failed{Error: InvalidGrant, Description: string(failure.Reason)}, nil
	}

	accessToken := grant.accessTokenIssuer.Issue(success.Username, request.Principal.Id, request.Scopes)

	refreshToken := grant.refreshTokenIssuer.Issue(success.Username, request.Principal.Id, request.Scopes)

	scopes := scope.MarshalScopes(accessToken.Scopes)

	return Success{
		AccessToken:  accessToken.Value,
		TokenType:    token.Bearer,
		ExpiresIn:    secondsBetween(accessToken.ExpiresAt, accessToken.IssuedAt),
		RefreshToken: refreshToken.Value,
		Scope:        scopes,
	}, nil
}

func secondsBetween(end time.Time, start time.Time) int64 {
	return int64(math.Round(end.Sub(start).Seconds()))
}
