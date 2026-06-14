package token_exchange

type Grant[Request any] interface {
	Exchange(request *Request) (Success, error)
}
