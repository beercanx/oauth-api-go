package user

import (
	"context"
	"database/sql"
	"errors"
)

var (
	ErrNoSuchUserCredential = errors.New("user credential does not exist")
)

type CredentialRepository interface {
	insert(new Credential) error
	FindByUsername(username string) (Credential, error)
}

type credentialRepository struct {
	ctx      context.Context
	database *sql.DB
}

// assert credentialRepository implements CredentialRepository
var _ CredentialRepository = (*credentialRepository)(nil)

func NewCredentialRepository(ctx context.Context, database *sql.DB) CredentialRepository {
	return &credentialRepository{ctx, database}
}

const insertUserCredential = `
INSERT INTO user_credentials (username, hash) 
VALUES (?, ?);
`

func (r credentialRepository) insert(credential Credential) error {
	_, err := r.database.ExecContext(r.ctx, insertUserCredential,
		credential.username,
		credential.hashedSecret,
	)
	return err
}

const findUserCredentialByUsername = `
SELECT username, hash 
FROM user_credentials 
WHERE username = ?;
`

func (r credentialRepository) FindByUsername(username string) (Credential, error) {
	credential, err := r.findOne(findUserCredentialByUsername, username)
	if errors.Is(err, sql.ErrNoRows) {
		return credential, ErrNoSuchUserCredential
	}
	return credential, err
}

func (r credentialRepository) findOne(query string, args ...any) (Credential, error) {
	row := r.database.QueryRowContext(r.ctx, query, args...)
	var credential Credential
	if scanError := row.Scan(&credential.username, &credential.hashedSecret); scanError != nil {
		return Credential{}, scanError
	}
	return credential, nil
}