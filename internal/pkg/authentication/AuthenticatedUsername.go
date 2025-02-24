package authentication

type AuthenticatedUsername struct {
	value string
}

func (username AuthenticatedUsername) Value() string {
	return username.value
}
