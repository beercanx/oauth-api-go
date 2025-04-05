package scope

import (
	"encoding/json"
	"errors"
)

type Scopes struct {
	Value []Scope
}

var _ json.Marshaler = (*Scopes)(nil)
var _ json.Unmarshaler = (*Scopes)(nil)

func (scopes Scopes) MarshalJSON() ([]byte, error) {
	return json.Marshal(marshalSpaceDelimited(scopes.Value, func(scope Scope) string {
		return scope.Value
	}))
}

func (scopes Scopes) UnmarshalJSON(_ []byte) error {
	return errors.ErrUnsupported // TODO - Implement
}
