package scope

import (
	"encoding/json"
	"fmt"
)

type Scope struct { // TODO: Should this be a type string rather than struct?
	Value string
}

var _ fmt.Stringer = (*Scope)(nil)
var _ json.Marshaler = (*Scope)(nil)
var _ json.Unmarshaler = (*Scope)(nil)

func (s Scope) String() string {
	return s.Value
}

func (s Scope) MarshalJSON() ([]byte, error) {
	return json.Marshal(s.Value)
}

func (s Scope) UnmarshalJSON(data []byte) error {
	return json.Unmarshal(data, &s.Value)
}
