package user

import "errors"

var (
	ErrNoSuchCredential = errors.New("credential does not exist")
)

type CredentialRepository interface {
	Insert(new Credential) error                        // TODO - Verify if Go DB libraries panic or return errors
	FindByUsername(username string) (Credential, error) // TODO - Verify if Go DB libraries panic or return errors
}
