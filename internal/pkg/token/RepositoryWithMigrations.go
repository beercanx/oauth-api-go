package token

type RepositoryWithMigrations[T any] interface {
	Repository[T]
	Migrate() error
}
