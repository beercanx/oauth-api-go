package user

import "errors"

var (
	ErrNoSuchStatus = errors.New("status does not exist")
)

type StatusRepository interface {
	Insert(status Status) error                     // TODO - Verify if Go DB libraries panic or return errors
	FindByUsername(username string) (Status, error) // TODO - Verify if Go DB libraries panic or return errors
}
