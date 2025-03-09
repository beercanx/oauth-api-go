package scope

type Repository interface {
	FindById(id string) (Scope, error)
}
