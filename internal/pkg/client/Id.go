package client

import (
	"encoding/json"
	"errors"
	"fmt"
)

type Id struct {
	Value string
}

var _ fmt.Stringer = (*Id)(nil)
var _ json.Marshaler = (*Id)(nil)
var _ json.Unmarshaler = (*Id)(nil)

func (i Id) String() string {
	return i.Value
}

func (i Id) MarshalJSON() ([]byte, error) {
	return json.Marshal(i.Value)
}

// UnmarshalJSON is purposely not supported to prevent deserialization of an Id from raw input.
func (i Id) UnmarshalJSON(_ []byte) error {
	return errors.ErrUnsupported
}
