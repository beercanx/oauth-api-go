package token

// ErrorResponse https://www.rfc-editor.org/rfc/rfc6749#section-5.2
type ErrorResponse struct {

	// Error A single ASCII error code from the defined list.
	Error ErrorType `json:"error"`

	// Description Human-readable ASCII text providing additional information, used
	// to assist the client developer in understanding the error that occurred.
	Description string `json:"error_description"`
}
