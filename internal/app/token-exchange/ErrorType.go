package token_exchange

type ErrorType string

const (
	// InvalidRequest The request is missing a required parameter, includes an
	// unsupported parameter value (other than a grant type), repeats a parameter,
	// includes multiple credentials, uses more than one mechanism for
	// authenticating the client, or is otherwise malformed.
	InvalidRequest ErrorType = "invalid_request"

	// InvalidClient Client authentication failed (e.g., unknown client, no client
	// authentication included, or unsupported authentication method). The
	// authorization server MAY return an HTTP 401 (Unauthorized) status code to
	// indicate which HTTP authentication schemes are supported. If the client
	// attempted to authenticate via the "Authorization" request header field, the
	// authorization server MUST respond with an HTTP 401 (Unauthorized) status code
	// and include the "WWW-Authenticate" response header field matching the
	// authentication scheme used by the client.
	InvalidClient ErrorType = "invalid_client"

	// InvalidGrant The provided authorization grant (e.g., authorization code,
	// resource owner credentials) or refresh token is invalid, expired, revoked,
	// does not match the redirection URI used in the authorization request, or was
	// issued to another client.
	InvalidGrant ErrorType = "invalid_grant"

	// InvalidScope The requested scope is invalid, unknown, malformed, or exceeds
	// the scope granted by the resource owner.
	InvalidScope ErrorType = "invalid_scope"

	// UnauthorizedClient The authenticated client is not authorized to use this
	// authorization grant type.
	UnauthorizedClient ErrorType = "unauthorized_client"

	// UnsupportedGrantType The authorization grant type is not supported by the
	// authorization server.
	UnsupportedGrantType ErrorType = "unsupported_grant_type"
)
