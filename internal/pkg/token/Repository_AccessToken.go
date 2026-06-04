package token

import (
	"context"
	"database/sql"
	_ "embed"
	"errors"

	"baconi.co.uk/oauth/internal/pkg/db"
	"github.com/google/uuid"
)

type accessTokenRepository struct {
	ctx     context.Context
	dbtx    db.DBTX
	queries *db.Queries
}

func (r accessTokenRepository) Insert(new db.CreateAccessTokenParams) error {
	return r.queries.CreateAccessToken(r.ctx, new)
}

func (r accessTokenRepository) FindById(id uuid.UUID) (db.AccessToken, error) {
	token, err := r.queries.GetAccessToken(r.ctx, id)
	if err != nil && errors.Is(err, sql.ErrNoRows) {
		return db.AccessToken{}, ErrNoSuchToken
	}
	return token, err
}

func (r accessTokenRepository) DeleteById(id uuid.UUID) error {
	return r.queries.DeleteAccessToken(r.ctx, id)
}

func (r accessTokenRepository) DeleteByRecord(record db.AccessToken) error {
	return r.DeleteById(record.ID)
}

func (r accessTokenRepository) DeletedExpired() error {
	return r.queries.DeleteExpiredAccessTokens(r.ctx)
}

var _ Repository[db.CreateAccessTokenParams, db.AccessToken] = (*accessTokenRepository)(nil)

func NewAccessTokenRepository(ctx context.Context, connection db.DBTX) Repository[db.CreateAccessTokenParams, db.AccessToken] {
	return &accessTokenRepository{ctx, connection, db.New(connection)}
}
