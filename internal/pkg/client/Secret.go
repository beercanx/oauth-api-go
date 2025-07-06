package client

import (
	"fmt"
	"github.com/google/uuid"
)

type Secret struct {
	id           uuid.UUID
	clientId     Id
	hashedSecret string
}

var _ fmt.Stringer = (*Secret)(nil)

func (s Secret) String() string {
	return fmt.Sprintf("Secret{id=%s, clientId=%s, hashedSecret='REDACTED'}", s.id, s.clientId)
}
