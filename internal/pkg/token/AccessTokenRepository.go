package token

import (
	"context"
	"database/sql"
	"errors"

	"github.com/google/uuid"
)

type accessTokenRepository struct {
	ctx      context.Context
	database *sql.DB
}

var _ Repository[AccessToken] = (*accessTokenRepository)(nil)

func NewAccessTokenRepository(ctx context.Context, database *sql.DB) Repository[AccessToken] {
	return &accessTokenRepository{ctx, database}
}

// #nosec G101 -- SQL query for a credential, not a credential itself
const insertAccessToken = `
INSERT INTO access_tokens (id, username, client_id, scopes, issued_at, expires_at, not_before)
VALUES (?, ?, ?, ?, ?, ?, ?)
`

func (r accessTokenRepository) Insert(new AccessToken) error {
	_, err := r.database.ExecContext(r.ctx, insertAccessToken,
		new.ID,
		new.Username,
		new.ClientID,
		new.Scopes,
		new.IssuedAt,
		new.ExpiresAt,
		new.NotBefore,
	)
	return err
}

// #nosec G101 -- SQL query for a credential, not a credential itself
const findAccessTokenById = `
SELECT id, username, client_id, json(scopes), issued_at, expires_at, not_before
FROM access_tokens
WHERE id = ?
LIMIT 1;
`

func (r accessTokenRepository) FindById(id uuid.UUID) (AccessToken, error) {
	token, err := r.findOne(findAccessTokenById, id)
	if errors.Is(err, sql.ErrNoRows) {
		return token, ErrNoSuchToken
	}
	return token, err
}

// #nosec G101 -- SQL query for a credential, not a credential itself
const deleteAccessTokenById = `
DELETE
FROM access_tokens
WHERE id = ?;
`

func (r accessTokenRepository) DeleteById(id uuid.UUID) error {
	_, err := r.database.ExecContext(r.ctx, deleteAccessTokenById, id)
	return err
}

// #nosec G101 -- SQL query for a credential, not a credential itself
const deleteExpiredAccessTokens = `
DELETE
FROM access_tokens
WHERE datetime(expires_at) <= CURRENT_TIMESTAMP;
`

func (r accessTokenRepository) DeletedExpired() error {
	_, err := r.database.ExecContext(r.ctx, deleteExpiredAccessTokens)
	return err
}

func (r accessTokenRepository) findOne(query string, args ...any) (AccessToken, error) {
	row := r.database.QueryRowContext(r.ctx, query, args...)
	var accessToken AccessToken
	queryError := row.Scan(
		&accessToken.ID,
		&accessToken.Username,
		&accessToken.ClientID,
		&accessToken.Scopes,
		&accessToken.IssuedAt,
		&accessToken.ExpiresAt,
		&accessToken.NotBefore,
	)
	if queryError != nil {
		return AccessToken{}, queryError
	}
	return accessToken, nil
}
