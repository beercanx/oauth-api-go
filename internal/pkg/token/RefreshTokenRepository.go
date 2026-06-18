package token

import (
	"context"
	"database/sql"
	"errors"

	"github.com/google/uuid"
)

type refreshTokenRepository struct {
	ctx      context.Context
	database *sql.DB
}

// assert refreshTokenRepository implements Repository
var _ Repository[RefreshToken] = (*refreshTokenRepository)(nil)

func NewRefreshTokenRepository(ctx context.Context, database *sql.DB) Repository[RefreshToken] {
	return &refreshTokenRepository{ctx, database}
}

// #nosec G101 -- SQL query for a credential, not a credential itself
const insertRefreshToken = `
INSERT INTO refresh_tokens (id, username, client_id, scopes, issued_at, expires_at, not_before)
VALUES (?, ?, ?, ?, ?, ?, ?)
`

func (r refreshTokenRepository) Insert(new RefreshToken) error {
	_, err := r.database.ExecContext(r.ctx, insertRefreshToken,
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
const findRefreshTokenById = `
SELECT id, username, client_id, json(scopes), issued_at, expires_at, not_before
FROM refresh_tokens
WHERE id = ?
LIMIT 1;
`

func (r refreshTokenRepository) FindById(id uuid.UUID) (RefreshToken, error) {
	token, err := r.findOne(findRefreshTokenById, id)
	if errors.Is(err, sql.ErrNoRows) {
		return token, ErrNoSuchToken
	}
	return token, err
}

// #nosec G101 -- SQL query for a credential, not a credential itself
const deleteRefreshTokenById = `
DELETE
FROM refresh_tokens
WHERE id = ?;
`

func (r refreshTokenRepository) DeleteById(id uuid.UUID) error {
	_, err := r.database.ExecContext(r.ctx, deleteRefreshTokenById, id)
	return err
}

// #nosec G101 -- SQL query for a credential, not a credential itself
const deleteExpiredRefreshTokens = `
DELETE
FROM refresh_tokens
WHERE datetime(expires_at) <= CURRENT_TIMESTAMP;
`

func (r refreshTokenRepository) DeletedExpired() error {
	_, err := r.database.ExecContext(r.ctx, deleteExpiredRefreshTokens)
	return err
}

func (r refreshTokenRepository) findOne(query string, args ...any) (RefreshToken, error) {
	row := r.database.QueryRowContext(r.ctx, query, args...)
	var refreshToken RefreshToken
	queryError := row.Scan(
		&refreshToken.ID,
		&refreshToken.Username,
		&refreshToken.ClientID,
		&refreshToken.Scopes,
		&refreshToken.IssuedAt,
		&refreshToken.ExpiresAt,
		&refreshToken.NotBefore,
	)
	if queryError != nil {
		return RefreshToken{}, queryError
	}
	return refreshToken, nil
}