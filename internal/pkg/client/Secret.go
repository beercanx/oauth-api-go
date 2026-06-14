package client

import (
	"fmt"

	"github.com/google/uuid"
)

type Secret struct {
	id           uuid.UUID `db:"id"`
	clientId     Id        `db:"client_id"`
	hashedSecret string    `db:"hash"`
	// TODO - make sure database has createdAt and updatedAt columns
}

var _ fmt.Stringer = (*Secret)(nil)

func (s Secret) String() string {
	return fmt.Sprintf("Secret{id: %s, clientId: %s}", s.id, s.clientId)
}
