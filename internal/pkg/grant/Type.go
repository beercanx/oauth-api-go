package grant

type Type string

const (
	AuthorisationCode Type = "authorization_code"
	Password          Type = "password"
	RefreshToken      Type = "refresh_token"
	Assertion         Type = "urn:ietf:params:oauth:grant-type:jwt-bearer"
)
