package client

import (
	"fmt"
	"testing"

	"baconi.co.uk/oauth/internal/pkg/db"
	"baconi.co.uk/oauth/internal/pkg/grant"
	"baconi.co.uk/oauth/internal/pkg/scope"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestPrincipalRepository(t *testing.T) {
	t.Parallel()

	database, databaseError := db.Connect("file:principal_repository_integration_tests?mode=memory&cache=shared")
	require.NoError(t, databaseError)
	require.NoError(t, db.RunMigrations(database, "file:../../../sql/migrations"))

	_, insertError := database.ExecContext(t.Context(), `
		INSERT INTO client_configurations 
		    (client_id, client_type, redirect_uris, allowed_scopes, allowed_actions, allowed_grant_types) 
		VALUES 
		    ('public-introspect', 'public', '[]', '[]', '["introspect"]', '[]'),
        ('public-password', 'public', '[]', '[]', '["authorise", "pkce"]', '["password"]'),
        ('public-without-pkce', 'public', '[]', '[]', '[]', '["authorization_code"]'),
        ('public-authorise-without-grant-type', 'public', '[]', '[]', '["authorise", "pkce"]', '[]'),
        ('public-authorise-without-redirect-uris', 'public', '[]', '[]', '["authorise", "pkce"]', '["authorization_code"]'),
        ('confidential-authorise-without-grant-type', 'confidential', '[]', '[]', '["authorise"]', '[]'),
        ('confidential-authorise-without-redirect-uris', 'confidential', '[]', '[]', '["authorise"]', '["authorization_code"]');
	`)
	require.NoError(t, insertError)

	underTest := NewPrincipalRepository(t.Context(), database)

	t.Run("FindById", func(t *testing.T) {
		testPrincipalRepository(t, func(id string) (Principal, error) {
			return underTest.FindById(Id(id))
		})
	})

	t.Run("FindByClientId", func(t *testing.T) {
		testPrincipalRepository(t, underTest.FindByClientId)
	})
}

func testPrincipalRepository(t *testing.T, underTest func(id string) (Principal, error)) {

	for _, clientId := range []string{"", " ", "no-such-client"} {
		t.Run(fmt.Sprintf("should return error on no such client: %s", clientId), func(t *testing.T) {
			principal, err := underTest(clientId)
			require.ErrorIs(t, err, ErrNoSuchClientPrincipal)
			assert.Zero(t, principal)
		})
	}

	for name, expectedError := range map[string]string{
		"public-introspect":                            "public clients must not be allowed to introspect",
		"public-password":                              "public clients must not use password grant",
		"public-without-pkce":                          "public clients must not use authorisation code grant without PKCE",
		"public-authorise-without-grant-type":          "clients with 'Authorise' must have 'AuthorisationCode'",
		"public-authorise-without-redirect-uris":       "clients with 'Authorise' must have some 'RedirectUris'",
		"confidential-authorise-without-grant-type":    "clients with 'Authorise' must have 'AuthorisationCode'",
		"confidential-authorise-without-redirect-uris": "clients with 'Authorise' must have some 'RedirectUris'",
	} {
		t.Run(fmt.Sprintf("should return error on invalid client configuration: %s", name), func(t *testing.T) {
			principal, err := underTest(name)
			require.ErrorIs(t, err, ErrPrincipalIsInvalid)
			require.ErrorContains(t, err, expectedError)
			assert.Zero(t, principal)
		})
	}

	t.Run("should be able to return a valid public client", func(t *testing.T) {
		principal, err := underTest("cicada")
		require.NoError(t, err)
		assert.Equal(t, Principal{
			ClientId:          "cicada",
			ClientType:        Public,
			RedirectUris:      RedirectUris{"https://cicada.baconi.co.uk/callback"},
			AllowedScopes:     scope.Scopes{"basic"},
			AllowedActions:    Actions{Authorise, ProofKeyForCodeExchange},
			AllowedGrantTypes: grant.Types{grant.AuthorisationCode},
		}, principal)
	})

	t.Run("should be able to return a valid confidential client", func(t *testing.T) {
		principal, err := underTest("aardvark")
		require.NoError(t, err)
		assert.Equal(t, Principal{
			ClientId:          "aardvark",
			ClientType:        Confidential,
			RedirectUris:      RedirectUris{},
			AllowedScopes:     scope.Scopes{"basic", "read", "write"},
			AllowedActions:    Actions{Introspect},
			AllowedGrantTypes: grant.Types{grant.Password},
		}, principal)
	})
}
