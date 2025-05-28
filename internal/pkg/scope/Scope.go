package scope

import (
	"encoding/json"
	"errors"
	"fmt"
)

type Scope struct {
	Value string
}

var _ fmt.Stringer = (*Scope)(nil)
var _ json.Marshaler = (*Scope)(nil)
var _ json.Unmarshaler = (*Scope)(nil)

func (s Scope) String() string {
	return s.Value
}

// MarshalJSON always returns unsupported because nothing should call this.
func (s Scope) MarshalJSON() ([]byte, error) {
	return nil, errors.ErrUnsupported
}

// UnmarshalJSON always returns unsupported because nothing should call this.
func (s Scope) UnmarshalJSON(_ []byte) error {
	return errors.ErrUnsupported
}

// Basic scope
//
// Deprecated: To be replaced by config sourced Clients
var Basic = Scope{"basic"}
