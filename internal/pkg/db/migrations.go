package db

import (
	"database/sql"
	"errors"
	"fmt"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/sqlite3"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

func RunMigrations(database *sql.DB, migrations string) error {

	driver, driverError := sqlite3.WithInstance(database, &sqlite3.Config{})
	if driverError != nil {
		return fmt.Errorf("initializing sqlite3 driver: %w", driverError)
	}

	migrator, migratorError := migrate.NewWithDatabaseInstance(migrations, "sqlite3", driver)
	if migratorError != nil {
		return fmt.Errorf("initializing migration: %w", migratorError)
	}

	if migrationError := migrator.Up(); migrationError != nil && !errors.Is(migrationError, migrate.ErrNoChange) {
		return migrationError
	}

	return nil
}
