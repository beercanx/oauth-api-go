package user

import (
	"encoding/json"
	"errors"
	"fmt"
)

type AuthenticatedUsername struct {
	Value string
}

var _ fmt.Stringer = (*AuthenticatedUsername)(nil)
var _ json.Marshaler = (*AuthenticatedUsername)(nil)
var _ json.Unmarshaler = (*AuthenticatedUsername)(nil)

func (u AuthenticatedUsername) String() string {
	return u.Value
}

func (u AuthenticatedUsername) MarshalJSON() ([]byte, error) {
	return json.Marshal(u.Value)
}

// UnmarshalJSON is purposely not supported to prevent deserialization of an AuthenticatedUsername from raw input.
func (u AuthenticatedUsername) UnmarshalJSON(_ []byte) error {
	return errors.ErrUnsupported
}
