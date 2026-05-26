package scope

import (
	"fmt"
)

type Scope string

var _ fmt.Stringer = (*Scope)(nil)

func (s Scope) String() string {
	return string(s)
}
