package client

import (
	"database/sql"
	"database/sql/driver"
	"encoding/json"
	"fmt"

	"github.com/google/uuid"
)

// SecretId exists because uuid.UUID is stored as a string, not as a binary blob in the database.
type SecretId uuid.UUID

var _ fmt.Stringer = (*SecretId)(nil)
var _ driver.Valuer = (*SecretId)(nil)
var _ sql.Scanner = (*SecretId)(nil)
var _ json.Marshaler = (*SecretId)(nil)

func (s SecretId) String() string {
	return uuid.UUID(s).String()
}

func (s SecretId) MarshalJSON() ([]byte, error) {
	return json.Marshal(uuid.UUID(s))
}

func (s SecretId) Value() (driver.Value, error) {
	return uuid.UUID(s).MarshalBinary()
}

//goland:noinspection GoMixedReceiverTypes
func (s *SecretId) Scan(src any) error {
	var rawUuid uuid.UUID
	if err := rawUuid.UnmarshalBinary(src.([]byte)); err != nil {
		return err
	}
	*s = SecretId(rawUuid)
	return nil
}

type Secret struct {
	id           SecretId `db:"id"`
	clientId     Id       `db:"client_id"`
	hashedSecret string   `db:"hash"`
	// TODO - make sure database has createdAt and updatedAt columns
}

var _ fmt.Stringer = (*Secret)(nil)

func (s Secret) String() string {
	return fmt.Sprintf("Secret{id: %s, clientId: %s}", s.id, s.clientId)
}
