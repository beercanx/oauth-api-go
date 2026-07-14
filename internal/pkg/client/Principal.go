package client

import (
	"errors"
	"fmt"
	"slices"

	"baconi.co.uk/oauth/internal/pkg/grant"
	"baconi.co.uk/oauth/internal/pkg/scope"
)

type Principal struct {
	ClientId          Id           `db:"client_id"`
	ClientType        Type         `db:"client_type"`
	RedirectUris      RedirectUris `db:"redirect_uris"`
	AllowedScopes     scope.Scopes `db:"allowed_scopes"`
	AllowedActions    Actions      `db:"allowed_actions"`
	AllowedGrantTypes grant.Types  `db:"allowed_grant_types"`
	// TODO - make sure database has createdAt and updatedAt columns
}

var ErrPrincipalIsInvalid = errors.New("principal is invalid")

// Validate that the Principal is configured correctly for its Type.
func (p Principal) Validate() error {

	if p.ClientType != Public && p.ClientType != Confidential {
		return fmt.Errorf("[%s] type cannot be [%s]: %w", p.ClientId, p.ClientType, ErrPrincipalIsInvalid)
	}

	if p.IsConfidential() {
		if p.ClientType != Confidential {
			return fmt.Errorf("[%s] type cannot be [%s]: %w", p.ClientId, p.ClientType, ErrPrincipalIsInvalid)
		}
	}

	if p.IsPublic() {
		if p.ClientType != Public {
			return fmt.Errorf("[%s] type cannot be [%s]: %w", p.ClientId, p.ClientType, ErrPrincipalIsInvalid)
		}
		if p.CanPerformAction(Introspect) {
			return fmt.Errorf("[%s] public clients must not be allowed to introspect: %w", p.ClientId, ErrPrincipalIsInvalid)
		}
		if p.CanBeGranted(grant.Password) {
			return fmt.Errorf("[%s] public clients must not use password grant: %w", p.ClientId, ErrPrincipalIsInvalid)
		}
		if p.CanBeGranted(grant.AuthorisationCode) && !p.CanPerformAction(ProofKeyForCodeExchange) {
			return fmt.Errorf("[%s] public clients must not use authorisation code grant without PKCE: %w", p.ClientId, ErrPrincipalIsInvalid)
		}
	}

	if p.CanPerformAction(Authorise) && !p.CanBeGranted(grant.AuthorisationCode) {
		return fmt.Errorf("[%s] clients with 'Authorise' must have 'AuthorisationCode': %w", p.ClientId, ErrPrincipalIsInvalid)
	}

	if p.CanPerformAction(Authorise) && len(p.RedirectUris) == 0 {
		return fmt.Errorf("[%s] clients with 'Authorise' must have some 'RedirectUris': %w", p.ClientId, ErrPrincipalIsInvalid)
	}

	return nil
}

func (p Principal) IsPublic() bool {
	return p.ClientType == Public
}

func (p Principal) IsConfidential() bool {
	return p.ClientType == Confidential
}

func (p Principal) CanBeGranted(grantType grant.Type) bool {
	return slices.Contains(p.AllowedGrantTypes, grantType)
}

func (p Principal) CanPerformAction(action Action) bool {
	return slices.Contains(p.AllowedActions, action)
}

func (p Principal) CanBeIssued(scopes []scope.Scope) bool {
	result := true
	for _, s := range scopes {
		if result = slices.Contains(p.AllowedScopes, s); !result {
			break
		}
	}
	return result
}
