package token

import (
	"errors"

	"github.com/google/uuid"
)

var (
	ErrNoSuchToken = errors.New("token does not exist")
)

type RepositoryCreate[N any] interface {
	Insert(new N) error
}

type RepositoryRead[R any] interface {
	FindById(id uuid.UUID) (R, error)
}

type RepositoryDelete[R any] interface {
	DeleteById(id uuid.UUID) error
	DeleteByRecord(record R) error
	DeletedExpired() error
}

type RepositoryReadDelete[R any] interface {
	RepositoryRead[R]
	RepositoryDelete[R]
}

type Repository[N any, R any] interface {
	RepositoryCreate[N]
	RepositoryRead[R]
	RepositoryDelete[R]
}
