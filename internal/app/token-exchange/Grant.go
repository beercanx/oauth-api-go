package token_exchange

type Grant[Request Valid] interface {
	Exchange(request Request) Response
}
