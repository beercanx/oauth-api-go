package client

type PrincipalRepository interface {
	FindById(id Id) (Principal, bool)
	FindByClientId(clientId string) (Principal, bool)
}
