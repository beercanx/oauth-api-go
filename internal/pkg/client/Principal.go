package client

import (
	"errors"
	"fmt"
	"slices"

	"baconi.co.uk/oauth/internal/pkg/grant"
	"baconi.co.uk/oauth/internal/pkg/scope"
)

type Principal struct {
	ClientID          Id
	ClientType        Type
	RedirectUris      RedirectUris
	AllowedScopes     scope.Scopes
	AllowedActions    Actions
	AllowedGrantTypes grant.Types
}

var ErrPrincipalIsInvalid = errors.New("principal is invalid")

// Validate that the Principal is configured correctly for its Type.
func (p Principal) Validate() error {

	if p.ClientType != Public && p.ClientType != Confidential {
		return fmt.Errorf("[%s] type cannot be [%s]: %w", p.ClientID, p.ClientType, ErrPrincipalIsInvalid)
	}

	if p.IsConfidential() {
		if p.ClientType != Confidential {
			return fmt.Errorf("[%s] type cannot be [%s]: %w", p.ClientID, p.ClientType, ErrPrincipalIsInvalid)
		}
	}

	if p.IsPublic() {
		if p.ClientType != Public {
			return fmt.Errorf("[%s] type cannot be [%s]: %w", p.ClientID, p.ClientType, ErrPrincipalIsInvalid)
		}
		if p.CanPerformAction(Introspect) {
			return fmt.Errorf("[%s] public clients must not be allowed to introspect: %w", p.ClientID, ErrPrincipalIsInvalid)
		}
		if p.CanBeGranted(grant.Password) {
			return fmt.Errorf("[%s] public clients must not use password grant: %w", p.ClientID, ErrPrincipalIsInvalid)
		}
		if p.CanBeGranted(grant.AuthorisationCode) && !p.CanPerformAction(ProofKeyForCodeExchange) {
			return fmt.Errorf("[%s] public clients must not use authorisation code grant without PKCE: %w", p.ClientID, ErrPrincipalIsInvalid)
		}
	}

	if p.CanPerformAction(Authorise) && !p.CanBeGranted(grant.AuthorisationCode) {
		return fmt.Errorf("[%s] clients with 'Authorise' must have 'AuthorisationCode': %w", p.ClientID, ErrPrincipalIsInvalid)
	}

	if p.CanPerformAction(Authorise) && len(p.RedirectUris) == 0 {
		return fmt.Errorf("[%s] clients with 'Authorise' must have some 'RedirectUris': %w", p.ClientID, ErrPrincipalIsInvalid)
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
