package token

type Grant[Request Valid] interface {
	Exchange(request Request) (*Success, *Failed, error)
}
