package user

import "errors"

var (
	ErrNoSuchCredential = errors.New("credential does not exist")
)

type CredentialRepository interface {
	Insert(new Credential) error
	FindByUsername(username string) (Credential, error)
}
