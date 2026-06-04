package db

import "database/sql"

func Connect(source string) (*sql.DB, error) {

	connection, err := sql.Open("sqlite3", source)
	if err != nil {
		return nil, err
	}

	return connection, nil
}
