package token_introspection

type invalid struct {
	ErrorType   errorType `json:"error"`
	Description string    `json:"description,omitempty,omitzero"`
}
