package client

import (
	"context"
	"database/sql"
	"errors"
	"log/slog"
)

var (
	ErrNoSuchClientSecret = errors.New("client secret does not exist")
)

type SecretRepository interface {
	FindById(id SecretId) (Secret, error)
	FindByClient(client Id) ([]Secret, error)
	FindByClientId(clientId string) ([]Secret, error)
}

type secretRepository struct {
	ctx      context.Context
	database *sql.DB
}

// assert secretRepository implements SecretRepository
var _ SecretRepository = (*secretRepository)(nil)

func NewSecretRepository(ctx context.Context, database *sql.DB) SecretRepository {
	return &secretRepository{ctx, database}
}

// #nosec G101 -- SQL query for a credential, not a credential itself
const findClientSecretById = `
SELECT id, client_id, hash
FROM client_secrets
WHERE id = ?
LIMIT 1;
`

func (i secretRepository) FindById(id SecretId) (Secret, error) {
	secret, err := i.findOne(findClientSecretById, id)
	if errors.Is(err, sql.ErrNoRows) {
		return secret, ErrNoSuchClientSecret
	}
	return secret, err
}

// #nosec G101 -- SQL query for a credential, not a credential itself
const findClientSecretsByClientId = `
SELECT id, client_id, hash
FROM client_secrets
WHERE client_id = ?;
`

func (i secretRepository) FindByClient(client Id) ([]Secret, error) {
	return i.findMany(findClientSecretsByClientId, client)
}

func (i secretRepository) FindByClientId(clientId string) ([]Secret, error) {
	return i.findMany(findClientSecretsByClientId, clientId)
}

func (i secretRepository) findOne(query string, args ...any) (Secret, error) {
	row := i.database.QueryRowContext(i.ctx, query, args...)
	var secret Secret
	if scanError := row.Scan(&secret.id, &secret.clientId, &secret.hashedSecret); scanError != nil {
		return Secret{}, scanError
	}
	return secret, nil
}

func (i secretRepository) findMany(query string, args ...any) ([]Secret, error) {

	rows, queryError := i.database.QueryContext(i.ctx, query, args...)
	if queryError != nil {
		return []Secret{}, queryError
	}

	defer func() {
		if err := rows.Close(); err != nil {
			slog.Error("Failed to close rows", slog.String("query", query), slog.Any("error", err))
		}
	}()

	var secrets []Secret

	for rows.Next() {
		var secret Secret
		if scanError := rows.Scan(&secret.id, &secret.clientId, &secret.hashedSecret); scanError != nil {
			return []Secret{}, scanError
		}
		secrets = append(secrets, secret)
	}

	if rowsError := rows.Err(); rowsError != nil {
		return []Secret{}, rowsError
	}

	return secrets, nil
}
