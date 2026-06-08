package sql

import (
	"context"
	"database/sql"
	"errors"
	"log"

	"baconi.co.uk/oauth/internal/pkg/client"
	"baconi.co.uk/oauth/internal/pkg/db"
)

type sqlPrincipalRepository struct {
	ctx     context.Context
	queries *db.Queries
}

func (r sqlPrincipalRepository) FindById(id client.Id) (client.Principal, error) {

	clientConfiguration, queryError := r.queries.GetClientConfiguration(r.ctx, id)
	if queryError != nil && errors.Is(queryError, sql.ErrNoRows) {
		return client.Principal{}, client.ErrNoSuchClient
	}
	if queryError != nil {
		log.Printf("Failed to retrieve client configuration for id %s: %v", id, queryError)
		return client.Principal{}, queryError
	}

	principal := client.Principal{
		Id:                clientConfiguration.ClientID,
		Type:              clientConfiguration.ClientType,
		RedirectUris:      clientConfiguration.RedirectUris,
		AllowedScopes:     clientConfiguration.AllowedScopes,
		AllowedActions:    clientConfiguration.AllowedActions,
		AllowedGrantTypes: clientConfiguration.AllowedGrantTypes,
	}

	if validateError := principal.Validate(); validateError != nil {
		return client.Principal{}, validateError
	}

	return principal, nil
}

func (r sqlPrincipalRepository) FindByClientId(clientId string) (client.Principal, error) {
	return r.FindById(client.Id(clientId))
}

var _ client.PrincipalRepository = (*sqlPrincipalRepository)(nil)

func NewSqlPrincipalRepository(ctx context.Context, connection db.DBTX) client.PrincipalRepository {
	return &sqlPrincipalRepository{ctx, db.New(connection)}
}
