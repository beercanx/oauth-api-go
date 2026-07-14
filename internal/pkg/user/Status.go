package user

type Status struct {
	username string `db:"username"`
	locked   bool   `db:"locked"`
}
