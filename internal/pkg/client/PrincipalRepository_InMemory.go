package client

import (
	"baconi.co.uk/oauth/internal/pkg/grant"
	"baconi.co.uk/oauth/internal/pkg/scope"
)

type InMemoryPrincipalRepository struct {
	byClientId map[Id]Principal
}

func (i InMemoryPrincipalRepository) insert(principal Principal) {
	i.byClientId[principal.Id] = principal
}

func (i InMemoryPrincipalRepository) FindById(id Id) (Principal, bool) {
	principal, ok := i.byClientId[id]
	if ok {
		principal.verify()
	}
	return principal, ok
}

func (i InMemoryPrincipalRepository) FindByClientId(clientId string) (Principal, bool) {
	principal, ok := i.byClientId[Id(clientId)]
	if ok {
		principal.verify()
	}
	return principal, ok
}

var _ PrincipalRepository = (*InMemoryPrincipalRepository)(nil)

func NewInMemoryPrincipalRepository() *InMemoryPrincipalRepository {
	repository := &InMemoryPrincipalRepository{make(map[Id]Principal)}

	repository.insert(Principal{
		Id:                "aardvark",
		Type:              Confidential,
		AllowedScopes:     []scope.Scope{"basic"},
		AllowedGrantTypes: []grant.Type{grant.Password},
		AllowedActions:    []Action{Introspect},
	})

	repository.insert(Principal{
		Id:                "cicada",
		Type:              Public,
		RedirectUris:      []string{"https://cicada.baconi.co.uk/callback"},
		AllowedScopes:     []scope.Scope{"basic"},
		AllowedGrantTypes: []grant.Type{grant.AuthorisationCode},
		AllowedActions:    []Action{Authorise, ProofKeyForCodeExchange},
	})

	repository.insert(Principal{
		Id:                "dodo",
		Type:              Confidential,
		RedirectUris:      []string{"https://dodo.baconi.co.uk/callback"},
		AllowedScopes:     []scope.Scope{"basic"},
		AllowedGrantTypes: []grant.Type{grant.AuthorisationCode},
		AllowedActions:    []Action{Authorise, ProofKeyForCodeExchange},
	})

	return repository
}
