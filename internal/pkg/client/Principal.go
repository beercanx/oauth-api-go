package client

import (
	"baconi.co.uk/oauth/internal/pkg/grant"
	"slices"
)

type Principal struct {
	Id                Id
	Type              Type
	RedirectUris      []string
	AllowedScopes     []string
	AllowedActions    []string
	AllowedGrantTypes []grant.Type
}

func (principal Principal) IsPublic() bool {
	return principal.Type == Public
}

func (principal Principal) IsConfidential() bool {
	return principal.Type == Confidential
}

func (principal Principal) Can(grantType grant.Type) bool {
	return slices.Contains(principal.AllowedGrantTypes, grantType)
}

func (principal Principal) CanBeIssued(scopes []string) bool {
	result := true
	for _, scope := range scopes {
		if result = slices.Contains(principal.AllowedScopes, scope); !result {
			break
		}
	}
	return result
}
