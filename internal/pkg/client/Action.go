package client

type Action string

const (
	Authorise               Action = "authorise"
	Introspect              Action = "introspect"
	ProofKeyForCodeExchange Action = "pkce"
)
