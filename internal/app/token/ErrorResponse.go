package token

type ErrorResponse struct {
	Error       string `json:"error"`
	Description string `json:"error_description"`
}
