package migrations

import (
	"errors"
	"fmt"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file" // Import for side effects
	"github.com/jmoiron/sqlx"

	"github.com/SurfShadow/surfshadow-server/pkg/logger"
)

const (
	errFailedToCreateMigrationDriver = "failed to create migration driver"
	errFailedToCreateMigrateInstance = "failed to create migrate instance"
	errFailedToApplyMigrations       = "failed to apply migrations"
	errFailedToRollbackMigration     = "failed to rollback migration"
	errFailedToDropAllMigrations     = "failed to drop all migrations"
)

func ApplyMigrations(db *sqlx.DB, migrationsPath string) error {
	logger.Instance.Debug("Starting to apply migrations")

	driver, err := postgres.WithInstance(db.DB, &postgres.Config{})
	if err != nil {
		return fmt.Errorf("%s: %w", errFailedToCreateMigrationDriver, err)
	}

	logger.Instance.Debug("Migration driver created successfully")

	m, err := migrate.NewWithDatabaseInstance(
		migrationsPath,
		"postgres",
		driver,
	)

	if err != nil {
		return fmt.Errorf("%s: %w", errFailedToCreateMigrateInstance, err)
	}

	logger.Instance.Debug("Migrate instance created successfully")

	if err = m.Up(); err != nil && !errors.Is(err, migrate.ErrNoChange) {
		return fmt.Errorf("%s: %w", errFailedToApplyMigrations, err)
	}

	logger.Instance.Info("Migrations applied successfully")

	return nil
}

func RollbackMigrations(db *sqlx.DB, migrationsPath string) error {
	logger.Instance.Debug("Starting to rollback migrations")

	driver, err := postgres.WithInstance(db.DB, &postgres.Config{})
	if err != nil {
		return fmt.Errorf("%s: %w", errFailedToCreateMigrationDriver, err)
	}

	logger.Instance.Debug("Migration driver created successfully")

	m, err := migrate.NewWithDatabaseInstance(
		migrationsPath,
		"postgres",
		driver,
	)

	if err != nil {
		return fmt.Errorf("%s: %w", errFailedToCreateMigrationDriver, err)
	}

	logger.Instance.Debug("Migrate instance created successfully")

	if err = m.Steps(-1); err != nil && !errors.Is(err, migrate.ErrNoChange) {
		return fmt.Errorf("%s: %w", errFailedToRollbackMigration, err)
	}

	logger.Instance.Info("Migrations rolled back successfully")

	return nil
}

func DropMigrations(db *sqlx.DB, migrationsPath string) error {
	logger.Instance.Debug("Starting to drop all migrations")

	driver, err := postgres.WithInstance(db.DB, &postgres.Config{})
	if err != nil {
		return fmt.Errorf("%s: %w", errFailedToCreateMigrationDriver, err)
	}

	logger.Instance.Debug("Migration driver created successfully")

	m, err := migrate.NewWithDatabaseInstance(
		migrationsPath,
		"postgres",
		driver,
	)

	if err != nil {
		return fmt.Errorf("%s: %w", errFailedToCreateMigrateInstance, err)
	}

	logger.Instance.Debug("Migrate instance created successfully")

	if err = m.Drop(); err != nil {
		return fmt.Errorf("%s: %w", errFailedToDropAllMigrations, err)
	}

	logger.Instance.Info("All migrations dropped successfully")

	return nil
}
