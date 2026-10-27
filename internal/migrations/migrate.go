package migrations

import (
	"errors"
	"fmt"
	"net/url"
	"os"
	"path/filepath"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

func Up(databaseURL, migrationsDir string) error {
	m, err := newMigrator(databaseURL, migrationsDir)
	if err != nil {
		return err
	}
	defer closeMigrator(m)

	if err := m.Up(); err != nil && !errors.Is(err, migrate.ErrNoChange) {
		return fmt.Errorf("apply migrations: %w", err)
	}
	return nil
}

// Down rolls back exactly one migration step.
func Down(databaseURL, migrationsDir string) error {
	m, err := newMigrator(databaseURL, migrationsDir)
	if err != nil {
		return err
	}
	defer closeMigrator(m)

	if err := m.Steps(-1); err != nil && !errors.Is(err, migrate.ErrNoChange) {
		return fmt.Errorf("roll back one migration: %w", err)
	}
	return nil
}

// newMigrator centralizes migration source and database initialization for all
// migration operations.
func newMigrator(databaseURL, migrationsDir string) (*migrate.Migrate, error) {
	if databaseURL == "" {
		return nil, fmt.Errorf("database URL is empty")
	}

	sourceURL, err := migrationSourceURL(migrationsDir)
	if err != nil {
		return nil, err
	}

	m, err := migrate.New(sourceURL, databaseURL)
	if err != nil {
		return nil, fmt.Errorf("initialize migrations: %w", err)
	}
	return m, nil
}

// migrationSourceURL resolves and validates the migrations directory in one place.
func migrationSourceURL(migrationsDir string) (string, error) {
	if migrationsDir == "" {
		migrationsDir = "migrations"
	}

	absDir, err := filepath.Abs(migrationsDir)
	if err != nil {
		return "", fmt.Errorf("resolve migrations directory: %w", err)
	}
	info, err := os.Stat(absDir)
	if err != nil {
		return "", fmt.Errorf("access migrations directory %q: %w", absDir, err)
	}
	if !info.IsDir() {
		return "", fmt.Errorf("migrations path %q is not a directory", absDir)
	}

	return (&url.URL{
		Scheme: "file",
		Path:   filepath.ToSlash(absDir),
	}).String(), nil
}

func closeMigrator(m *migrate.Migrate) {
	_, _ = m.Close()
}
