package user

import (
	"context"
	"database/sql"
	"errors"
)

var (
	ErrNoSuchUserStatus = errors.New("user status does not exist")
)

type StatusRepository interface {
	insert(status Status) error
	FindByUsername(username string) (Status, error)
}

type statusRepository struct {
	ctx      context.Context
	database *sql.DB
}

// assert statusRepository implements StatusRepository
var _ StatusRepository = (*statusRepository)(nil)

func NewStatusRepository(ctx context.Context, database *sql.DB) StatusRepository {
	return &statusRepository{ctx, database}
}

const insertUserStatus = `
INSERT INTO user_statuses (username, locked) 
VALUES (?, ?);
`

func (r statusRepository) insert(status Status) error {
	_, err := r.database.ExecContext(r.ctx, insertUserStatus,
		status.username,
		status.locked,
	)
	return err
}

const findUserStatusByUsername = `
SELECT username, locked 
FROM user_statuses 
WHERE username = ?;
`

func (r statusRepository) FindByUsername(username string) (Status, error) {
	status, err := r.findOne(findUserStatusByUsername, username)
	if errors.Is(err, sql.ErrNoRows) {
		return status, ErrNoSuchUserStatus
	}
	return status, err
}

func (r statusRepository) findOne(query string, args ...any) (Status, error) {
	row := r.database.QueryRowContext(r.ctx, query, args...)
	var status Status
	if scanError := row.Scan(&status.username, &status.locked); scanError != nil {
		return Status{}, scanError
	}
	return status, nil
}
