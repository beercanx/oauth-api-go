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

func (i InMemoryPrincipalRepository) FindById(id Id) (Principal, error) {
	principal, ok := i.byClientId[id]
	if !ok {
		return Principal{}, ErrNoSuchClient
	}
	if validateError := principal.Validate(); validateError != nil {
		return Principal{}, validateError
	}
	return principal, nil
}

func (i InMemoryPrincipalRepository) FindByClientId(clientId string) (Principal, error) {
	return i.FindById(Id(clientId))
}

var _ PrincipalRepository = (*InMemoryPrincipalRepository)(nil)

func NewInMemoryPrincipalRepository() PrincipalRepository {
	repository := &InMemoryPrincipalRepository{make(map[Id]Principal)}

	repository.insert(Principal{
		Id:                "aardvark",
		Type:              Confidential,
		AllowedScopes:     scope.Scopes{"basic", "read", "write"},
		AllowedGrantTypes: grant.Types{grant.Password},
		AllowedActions:    Actions{Introspect},
	})

	repository.insert(Principal{
		Id:                "cicada",
		Type:              Public,
		RedirectUris:      RedirectUris{"https://cicada.baconi.co.uk/callback"},
		AllowedScopes:     scope.Scopes{"basic"},
		AllowedGrantTypes: grant.Types{grant.AuthorisationCode},
		AllowedActions:    Actions{Authorise, ProofKeyForCodeExchange},
	})

	repository.insert(Principal{
		Id:                "dodo",
		Type:              Confidential,
		RedirectUris:      RedirectUris{"https://dodo.baconi.co.uk/callback"},
		AllowedScopes:     scope.Scopes{"basic"},
		AllowedGrantTypes: grant.Types{grant.AuthorisationCode},
		AllowedActions:    Actions{Authorise, ProofKeyForCodeExchange},
	})

	return repository
}
