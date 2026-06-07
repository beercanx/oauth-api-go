package client

import "errors"

var (
	ErrNoSuchClient = errors.New("client does not exist")
)

type PrincipalRepository interface {
	FindById(id Id) (Principal, error)
	FindByClientId(clientId string) (Principal, error)
}
