package db

import (
	"context"
	"fmt"
	"log"
	"os"
	"path/filepath"
)

func RunMigrations(ctx context.Context, database DBTX) error {

	files, globError := filepath.Glob("../../../sqlc/migrations/*.up.sql")
	if globError != nil {
		return fmt.Errorf("finding migration files: %w", globError)
	}

	log.Printf("Migrating database: %s", files)

	for _, dirtyFile := range files {

		cleanFile := filepath.Clean(dirtyFile)
		log.Printf("Reading migration: %s", cleanFile)

		content, readError := os.ReadFile(cleanFile)
		if readError != nil {
			return fmt.Errorf("reading migration file %s: %w", cleanFile, readError)
		}

		log.Printf("Running migration: %s", content)

		if _, execError := database.ExecContext(ctx, string(content)); execError != nil {
			return fmt.Errorf("executing migration %s: %w", cleanFile, execError)
		}
	}

	return nil
}
