package db

import (
	"database/sql"

	_ "github.com/ncruces/go-sqlite3/driver"
)

func Connect(source string) (*sql.DB, error) {

	connection, err := sql.Open("sqlite3", source)
	if err != nil {
		return nil, err
	}

	return connection, nil
}
