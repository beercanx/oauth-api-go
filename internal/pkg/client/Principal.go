package client

import (
	"baconi.co.uk/oauth/internal/pkg/grant"
	"baconi.co.uk/oauth/internal/pkg/scope"
	"fmt"
	"slices"
)

type Principal struct {
	Id                Id
	Type              Type
	RedirectUris      []string
	AllowedScopes     []scope.Scope
	AllowedActions    []Action
	AllowedGrantTypes []grant.Type
}

// Verify that the Principal is configured correctly for its Type. TODO - Work out if we still want this and where it goes.
func (p Principal) verify() {

	require(p.Type == Public || p.Type == Confidential, fmt.Sprintf("[%s] type cannot be [%s]", p.Id, p.Type))

	if p.IsConfidential() {
		require(p.Type == Confidential, fmt.Sprintf("[%s] type cannot be [%s]", p.Id, p.Type))
	}

	if p.IsPublic() {
		require(p.Type == Public, fmt.Sprintf("[%s] type cannot be [%s]", p.Id, p.Type))
		require(!p.CanPerformAction(Introspect), fmt.Sprintf("public clients must not be allowed to introspect: %s", p.Id))
		require(!p.CanBeGranted(grant.Password), fmt.Sprintf("public clients must not use password grant: %s", p.Id))
		require(!(p.CanBeGranted(grant.AuthorisationCode) && !p.CanPerformAction(ProofKeyForCodeExchange)),
			fmt.Sprintf("public clients must not use authorisation code grant without PKCE: %s", p.Id),
		)
	}

	require(!(p.CanPerformAction(Authorise) && !p.CanBeGranted(grant.AuthorisationCode)), // TODO - Replace with implied action based on grant type?
		fmt.Sprintf("clients with 'Authorise' must have 'AuthorisationCode': %s", p.Id),
	)

	require(!(p.CanPerformAction(Authorise) && len(p.RedirectUris) <= 0),
		fmt.Sprintf("clients with 'Authorise' must have some 'RedirectUris': %s", p.Id),
	)
}

func (p Principal) IsPublic() bool {
	return p.Type == Public
}

func (p Principal) IsConfidential() bool {
	return p.Type == Confidential
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
