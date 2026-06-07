package sql

import (
	"fmt"
	"testing"

	"baconi.co.uk/oauth/internal/pkg/client"
	"baconi.co.uk/oauth/internal/pkg/db"
	"baconi.co.uk/oauth/internal/pkg/grant"
	"baconi.co.uk/oauth/internal/pkg/scope"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestSqlPrincipalRepository_FindById(t *testing.T) {
	t.Parallel()

	database, databaseError := db.Connect("file:principal_repository_integration_tests?mode=memory&cache=shared")
	require.NoError(t, databaseError)
	require.NoError(t, db.RunMigrations(database, "file:../../../../sqlc/migrations"))

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

	underTest := NewSqlPrincipalRepository(t.Context(), database)

	for _, clientId := range []string{"", " ", "no-such-client"} {
		t.Run(fmt.Sprintf("should return error on no such client: %s", clientId), func(t *testing.T) {
			principal, err := underTest.FindById(client.Id(clientId))
			assert.ErrorIs(t, err, client.ErrNoSuchClient)
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
			principal, err := underTest.FindById(client.Id(name))
			assert.ErrorContains(t, err, expectedError)
			assert.Zero(t, principal)
		})
	}

	t.Run("should be able to return a valid public client", func(t *testing.T) {
		principal, err := underTest.FindById("cicada")
		assert.NoError(t, err)
		assert.Equal(t, client.Principal{
			Id:                "cicada",
			Type:              client.Public,
			RedirectUris:      client.RedirectUris{"https://cicada.baconi.co.uk/callback"},
			AllowedScopes:     scope.Scopes{"basic"},
			AllowedActions:    client.Actions{client.Authorise, client.ProofKeyForCodeExchange},
			AllowedGrantTypes: grant.Types{grant.AuthorisationCode},
		}, principal)
	})

	t.Run("should be able to return a valid confidential client", func(t *testing.T) {
		principal, err := underTest.FindById("aardvark")
		assert.NoError(t, err)
		assert.Equal(t, client.Principal{
			Id:                "aardvark",
			Type:              client.Confidential,
			RedirectUris:      client.RedirectUris{},
			AllowedScopes:     scope.Scopes{"basic", "read", "write"},
			AllowedActions:    client.Actions{client.Introspect},
			AllowedGrantTypes: grant.Types{grant.Password},
		}, principal)
	})
}
