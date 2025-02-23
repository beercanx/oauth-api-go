package token

import (
	"baconi.co.uk/oauth/internal/pkg/authentication"
	"baconi.co.uk/oauth/internal/pkg/token"
	"math"
	"time"
)

type PasswordGrant struct {
	accessTokenIssuer  token.Issuer[token.AccessToken]
	refreshTokenIssuer token.Issuer[token.RefreshToken]
	userAuthenticator  authentication.UserAuthenticator
}

func NewPasswordGrant(accessTokenIssuer token.Issuer[token.AccessToken], refreshTokenIssuer token.Issuer[token.RefreshToken], userAuthenticator authentication.UserAuthenticator) *PasswordGrant {
	return &PasswordGrant{accessTokenIssuer, refreshTokenIssuer, userAuthenticator}
}

func (grant PasswordGrant) Exchange(request PasswordRequest) (*Success, *Failed, error) {

	success, failure, err := grant.userAuthenticator.Authenticate(request.Username, request.Password)
	if err != nil {
		return nil, nil, err
	}

	if failure != nil {
		return nil, &Failed{Error: InvalidGrant, Description: string(failure.Reason)}, nil
	}

	accessToken, err := grant.accessTokenIssuer.Issue(success.Username, request.Principal.Id, request.Scopes)
	if err != nil {
		return nil, nil, err
	}

	refreshToken, err := grant.refreshTokenIssuer.Issue(success.Username, request.Principal.Id, request.Scopes)
	if err != nil {
		return nil, nil, err
	}

	return &Success{
		AccessToken:  accessToken.Value,
		TokenType:    token.Bearer,
		ExpiresIn:    secondsBetween(accessToken.ExpiresAt, accessToken.IssuedAt),
		RefreshToken: &refreshToken.Value,
		Scope:        &accessToken.Scopes,
		State:        nil,
	}, nil, nil
}

func secondsBetween(end time.Time, start time.Time) int64 {
	return int64(math.Round(end.Sub(start).Seconds()))
}
