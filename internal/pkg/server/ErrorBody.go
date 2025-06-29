package server

type ErrorBody struct {
	ErrorType   string `json:"error,omitempty,omitzero"`
	Description string `json:"description,omitempty,omitzero"`
}
