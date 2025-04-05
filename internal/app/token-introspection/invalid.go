package token_introspection

import "fmt"

type invalid struct {
	ErrorType   errorType `json:"error"`
	Description string    `json:"description,omitempty,omitzero"`
}

// assert invalid implements error
var _ error = (*invalid)(nil)

func (i invalid) Error() string {
	return fmt.Sprintf("%s: %s", i.ErrorType, i.Description)
}
