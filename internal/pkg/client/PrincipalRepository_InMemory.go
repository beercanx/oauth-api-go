package client

import (
	"baconi.co.uk/oauth/internal/pkg/grant"
	"baconi.co.uk/oauth/internal/pkg/scope"
)

type InMemoryPrincipalRepository struct {
	byClientId map[string]Principal
}

func (i InMemoryPrincipalRepository) insert(principal Principal) {
	i.byClientId[principal.Id.Value] = principal
}

func (i InMemoryPrincipalRepository) FindById(id Id) (Principal, bool) {
	principal, ok := i.byClientId[id.Value]
	if ok {
		principal.verify()
	}
	return principal, ok
}

func (i InMemoryPrincipalRepository) FindByClientId(clientId string) (Principal, bool) {
	principal, ok := i.byClientId[clientId]
	if ok {
		principal.verify()
	}
	return principal, ok
}

var _ PrincipalRepository = (*InMemoryPrincipalRepository)(nil)

func NewInMemoryPrincipalRepository() *InMemoryPrincipalRepository {
	repository := &InMemoryPrincipalRepository{make(map[string]Principal)}

	repository.insert(Principal{
		Id:                Id{"aardvark"},
		Type:              Confidential,
		AllowedScopes:     []scope.Scope{{Value: "basic"}},
		AllowedGrantTypes: []grant.Type{grant.Password},
		AllowedActions:    []Action{Introspect},
	})

	repository.insert(Principal{
		Id:                Id{"cicada"},
		Type:              Public,
		RedirectUris:      []string{"https://cicada.baconi.co.uk/callback"},
		AllowedScopes:     []scope.Scope{{Value: "basic"}},
		AllowedGrantTypes: []grant.Type{grant.AuthorisationCode},
		AllowedActions:    []Action{Authorise, ProofKeyForCodeExchange},
	})

	repository.insert(Principal{
		Id:                Id{"dodo"},
		Type:              Confidential,
		RedirectUris:      []string{"https://dodo.baconi.co.uk/callback"},
		AllowedScopes:     []scope.Scope{{Value: "basic"}},
		AllowedGrantTypes: []grant.Type{grant.AuthorisationCode},
		AllowedActions:    []Action{Authorise, ProofKeyForCodeExchange},
	})

	return repository
}
