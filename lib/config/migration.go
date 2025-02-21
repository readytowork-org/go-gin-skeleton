package config

import (
	"errors"
	"fmt"
	"regexp"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/mysql"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

// Migrations Migration Struct
type Migrations struct {
	logger   Logger
	migrator *migrate.Migrate
}

// NewMigrations return new Migrations struct
func NewMigrations(
	logger Logger,
	envPath EnvPath,
	db DBDialect,
) *Migrations {
	path := getMigrationFolder(envPath.ToString())
	path = fmt.Sprintf("file://%s/", path)

	migrator, err := migrate.New(path, fmt.Sprintf("%v://%v", db.Name(), db.DSN))
	if err != nil {
		logger.Panic("Error in migration: ", err)
	}

	return &Migrations{
		logger:   logger,
		migrator: migrator,
	}
}

// MigrateUp migrates all table
func (m Migrations) MigrateUp() error {
	m.logger.Info("--- Running Migration Up ---")
	err := m.migrator.Up()
	if err != nil && !errors.Is(err, migrate.ErrNoChange) {
		return err
	}
	version, _, err := m.migrator.Version()
	if err != nil {
		return err
	}
	m.logger.Infof("--- Migration Success; Current Version: %v ---\n", version)
	return nil
}

/*
getMigrationFolder path from env path.

e.g:

	../../<.test.env/.env> => ../../migration
	<.test.env/.env> => migration
*/
func getMigrationFolder(envPath string) string {
	m1 := regexp.MustCompile(`(\.(\w+))+`)
	return m1.ReplaceAllString(envPath, "database/migration")
}
