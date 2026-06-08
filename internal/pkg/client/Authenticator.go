package client

// Authenticator TODO - Decide if we would continue returning bool, or switch to nil Principal and error.
type Authenticator interface {
	AuthenticateAsPublic(clientId string) (Principal, bool)
	AuthenticateAsConfidential(clientId string, clientSecret string) (Principal, bool)
}
