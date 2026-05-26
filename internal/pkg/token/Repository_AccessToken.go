package token

import (
	"context"
	"database/sql"
	_ "embed"
	"errors"
	"fmt"
	"log"
	"os"
	"path/filepath"

	"baconi.co.uk/oauth/internal/pkg/client"
	"baconi.co.uk/oauth/internal/pkg/db"
	"baconi.co.uk/oauth/internal/pkg/user"
	"github.com/google/uuid"

	_ "github.com/ncruces/go-sqlite3/driver"
)

func NewAccessTokenRepository(ctx context.Context) (Repository[db.AccessToken], error) {
	return newAccessTokenRepository(ctx, "file:out/database.sqlite")
}

type RepositoryWithMigrations[T any] interface {
	Repository[T]
	Migrate() error
}

func (r accessTokenRepository) Migrate() error {

	files, globError := filepath.Glob("../../../sqlc/migrations/*.up.sql")
	if globError != nil {
		return fmt.Errorf("finding migration files: %w", globError)
	}

	log.Printf("Migrating database: %s", files)

	for _, file := range files {

		cleanFile := filepath.Clean(file)
		log.Printf("Reading migration: %s", cleanFile)

		content, readError := os.ReadFile(cleanFile)
		if readError != nil {
			return fmt.Errorf("reading migration file %s: %w", file, readError)
		}

		log.Printf("Running migration: %s", content)

		if _, execError := r.dbtx.ExecContext(r.ctx, string(content)); execError != nil {
			return fmt.Errorf("executing migration %s: %w", file, execError)
		}
	}

	return nil
}

// NewInMemoryAccessTokenRepository should only be used in unit/integration tests
func NewInMemoryAccessTokenRepository(ctx context.Context) RepositoryWithMigrations[db.AccessToken] {
	repo, err := newAccessTokenRepository(ctx, "file:access_tokens?mode=memory&cache=shared")
	if err != nil {
		panic(err)
	}
	return repo
}

func newAccessTokenRepository(ctx context.Context, source string) (RepositoryWithMigrations[db.AccessToken], error) {

	connection, err := sql.Open("sqlite3", source)
	if err != nil {
		return nil, err
	}

	queries := db.New(connection)

	return &accessTokenRepository{ctx, connection, queries}, nil
}

type accessTokenRepository struct {
	ctx     context.Context
	dbtx    db.DBTX
	queries *db.Queries
}

var _ Repository[db.AccessToken] = (*accessTokenRepository)(nil)

func (r accessTokenRepository) Insert(new db.AccessToken) error {
	return r.queries.CreateAccessToken(r.ctx, db.CreateAccessTokenParams(new))
}

func (r accessTokenRepository) FindById(id uuid.UUID) (db.AccessToken, error) {
	token, err := r.queries.GetAccessToken(r.ctx, id)
	if err != nil && errors.Is(err, sql.ErrNoRows) {
		return db.AccessToken{}, ErrNoSuchToken
	}
	return token, err
}

func (r accessTokenRepository) FindAllByUsername(username user.AuthenticatedUsername) ([]db.AccessToken, error) {
	panic("implement me")
}

func (r accessTokenRepository) FindAllByClientId(clientId client.Id) ([]db.AccessToken, error) {
	panic("implement me")
}

func (r accessTokenRepository) DeleteById(id uuid.UUID) error { // TODO - Add integration test for safe deleting
	return r.queries.DeleteAccessToken(r.ctx, id)
}

func (r accessTokenRepository) DeleteByRecord(record db.AccessToken) error { // TODO - Add integration test for safe deleting
	return r.DeleteById(record.ID)
}

func (r accessTokenRepository) DeletedExpired() error { // TODO - Add integration test for safe deleting
	panic("implement me")
}
