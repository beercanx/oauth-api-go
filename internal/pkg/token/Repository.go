package token

import (
	"errors"

	"github.com/google/uuid"
)

var (
	ErrNoSuchToken = errors.New("token does not exist")
)

type RepositoryCreate[R any] interface {
	Insert(new R) error
}

type RepositoryRead[R any] interface {
	FindById(id uuid.UUID) (R, error)
}

type RepositoryDelete[R any] interface {
	DeleteById(id uuid.UUID) error
	DeletedExpired() error
}

type RepositoryReadDelete[R any] interface {
	RepositoryRead[R]
	RepositoryDelete[R]
}

type Repository[R any] interface {
	RepositoryCreate[R]
	RepositoryRead[R]
	RepositoryDelete[R]
}
