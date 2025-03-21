package scope

import "errors"

var (
	ErrNoSuchScope = errors.New("scope does not exist")
)

type Repository interface {
	FindById(id string) (Scope, error)
}
