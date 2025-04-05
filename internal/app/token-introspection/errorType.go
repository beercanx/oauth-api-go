package token_introspection

type errorType string

const (
	InvalidRequest     errorType = "invalid_request"
	UnauthorizedClient errorType = "unauthorized_client"
)
