package user

import "errors"

var (
	ErrNoSuchStatus = errors.New("status does not exist")
)

type StatusRepository interface {
	Insert(status Status) error
	FindByUsername(username string) (Status, error)
}
