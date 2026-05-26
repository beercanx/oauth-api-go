package client

import (
	"encoding/json"
	"errors"
	"fmt"
)

type Id string

var _ fmt.Stringer = (*Id)(nil)
var _ json.Unmarshaler = (*Id)(nil)

func (id Id) String() string {
	return string(id)
}

// UnmarshalJSON is purposely not supported to prevent deserialization of an Id from raw input.
func (id Id) UnmarshalJSON(_ []byte) error {
	return errors.ErrUnsupported
}
