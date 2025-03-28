package client

type Authenticator interface {
	AuthenticateAsPublic(clientId string) (Principal, bool)
	AuthenticateAsConfidential(clientId string, clientSecret string) (Principal, bool)
}
