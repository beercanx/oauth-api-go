package client

import (
	"context"
	"database/sql"
	"errors"
	"log"
)

var (
	ErrNoSuchClientPrincipal = errors.New("client principal does not exist")
)

type PrincipalRepository interface {
	FindById(id Id) (Principal, error)
	FindByClientId(clientId string) (Principal, error)
}

type principalRepository struct {
	ctx      context.Context
	database *sql.DB
}

const findClientConfigurationByClientId = `
SELECT client_id, client_type, json(redirect_uris), json(allowed_scopes), json(allowed_actions), json(allowed_grant_types)
FROM client_configurations
WHERE client_id = ?
LIMIT 1;
`

func (r principalRepository) FindById(id Id) (Principal, error) {

	principal, queryError := r.findOne(findClientConfigurationByClientId, id)

	if errors.Is(queryError, sql.ErrNoRows) {
		return principal, ErrNoSuchClientPrincipal
	}

	if queryError != nil {
		log.Printf("Failed to retrieve client configuration for id %s: %v", id, queryError)
		return principal, queryError
	}

	if validateError := principal.Validate(); validateError != nil {
		return Principal{}, validateError
	}

	return principal, nil
}

func (r principalRepository) FindByClientId(clientId string) (Principal, error) {
	return r.FindById(Id(clientId))
}

func (r principalRepository) findOne(query string, args ...any) (Principal, error) {
	row := r.database.QueryRowContext(r.ctx, query, args...)
	var principal Principal
	queryError := row.Scan(
		&principal.ClientId,
		&principal.ClientType,
		&principal.RedirectUris,
		&principal.AllowedScopes,
		&principal.AllowedActions,
		&principal.AllowedGrantTypes,
	)
	if queryError != nil {
		return Principal{}, queryError
	}
	return principal, nil
}

var _ PrincipalRepository = (*principalRepository)(nil)

func NewPrincipalRepository(ctx context.Context, database *sql.DB) PrincipalRepository {
	return &principalRepository{ctx, database}
}
