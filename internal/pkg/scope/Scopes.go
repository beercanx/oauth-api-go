package scope

import (
	"database/sql"
	"database/sql/driver"
	"encoding/json"
	"errors"
	"fmt"
)

type Scopes []Scope

var _ json.Marshaler = (*Scopes)(nil)
var _ json.Unmarshaler = (*Scopes)(nil)
var _ driver.Valuer = (*Scopes)(nil)
var _ sql.Scanner = (*Scopes)(nil)

func (scopes Scopes) MarshalJSON() ([]byte, error) {
	return json.Marshal(marshalSpaceDelimited(scopes, func(scope Scope) string {
		return string(scope)
	}))
}

func (scopes Scopes) UnmarshalJSON(_ []byte) error {
	return errors.ErrUnsupported
}

func (scopes Scopes) Value() (driver.Value, error) {
	marshaled, err := json.Marshal([]Scope(scopes))
	if err != nil {
		return nil, err
	}
	return string(marshaled), nil
}

//goland:noinspection GoMixedReceiverTypes
func (scopes *Scopes) Scan(src any) error {
	var source string
	switch v := src.(type) {
	case string:
		source = v
	case []byte:
		source = string(v)
	default:
		return fmt.Errorf("unsupported type for Scopes: %T", src)
	}
	var rawScopes []Scope
	if err := json.Unmarshal([]byte(source), &rawScopes); err != nil {
		return err
	}
	*scopes = rawScopes
	return nil
}
