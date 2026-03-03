package config

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strconv"
	"strings"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/mysql"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"gorm.io/gorm"
)

// Migrations Migration Struct
// migratorIface abstracts the methods from github.com/golang-migrate/migrate
// that we rely on; using an interface makes it easier to substitute a fake
// implementation during unit tests.
type migratorIface interface {
	Up() error
	Down() error
	// the second return value is the "dirty" flag
	Version() (uint, bool, error)
}

// Migrations Migration Struct
// previously the migrator field was concrete *migrate.Migrate but that
// tied tests to an actual database driver.  Using the interface gives us
// flexibility without impacting production behaviour.
type Migrations struct {
	logger          Logger
	migrator        migratorIface
	migrationFolder string   // absolute path without scheme, used for manual gap filling
	db              *gorm.DB // underlying connection for executing missing SQL
}

// NewMigrations returns a new Migrations helper.  It requires the
// usual logger and environment settings plus the database dialect used by
// the migrate library and a *Database instance which is needed for
// manual gap-filling when migrations are applied out of order.
func NewMigrations(
	logger Logger,
	envPath EnvPath,
	dbDialect DBDialect,
	database *Database,
) *Migrations {
	folder := getMigrationFolder(envPath.ToString())
	path := fmt.Sprintf("file://%s/", folder)

	migrator, err := migrate.New(path, fmt.Sprintf("%v://%v", dbDialect.Name(), dbDialect.DSN))
	if err != nil {
		logger.Panic("Error in migration: ", err)
	}

	return &Migrations{
		logger:          logger,
		migrator:        migrator,
		migrationFolder: folder,
		db:              database.DB,
	}
}

// MigrateUp migrates all migrations in order and then applies any
// earlier‑versioned files that were skipped due to being added after the
// database was already advanced past their version.  Gopher teams often
// create migrations in parallel; if a later migration is merged first the
// earlier one will not execute.  After the normal migrate.Up we scan the
// migration directory and execute any missing versions manually.
func (m Migrations) MigrateUp() error {
	m.logger.Info("--- Running Migration Up ---")
	err := m.migrator.Up()
	if err != nil && !errors.Is(err, migrate.ErrNoChange) {
		return err
	}

	// attempt to fill any gaps left behind by out‑of‑order files
	if err := m.fillGapMigrations(); err != nil {
		return err
	}

	// ignore dirty flag (second return value)
	version, _, err := m.migrator.Version()
	if err != nil {
		return err
	}
	m.logger.Infof("--- Migration Success; Current Version: %v ---\n", version)
	return nil
}

// MigrateDown rolls back all applied migrations to version 0.
// Unlike MigrateUp which uses the underlying migrator, Down manually
// applies .down.sql files in reverse version order.  This is necessary
// because when migrations are applied out-of-order via gap filling, the
// underlying migrator's Down() method cannot handle non-sequential versions
// in the schema_migrations table.  We instead scan the table, sort versions
// descending, apply each .down.sql file, and remove the record.
func (m Migrations) MigrateDown() error {
	m.logger.Info("--- Running Migration Down ---")
	if m.db == nil {
		// delegate to standard migrator if no DB connection
		err := m.migrator.Down()
		if err != nil && !errors.Is(err, migrate.ErrNoChange) {
			return err
		}
		return nil
	}

	// fetch all applied versions from schema_migrations
	var versions []int64
	rows, err := m.db.Raw("SELECT version FROM schema_migrations ORDER BY version DESC").Rows()
	if err != nil {
		return err
	}
	defer rows.Close()
	for rows.Next() {
		var v int64
		if err := rows.Scan(&v); err != nil {
			return err
		}
		versions = append(versions, v)
	}

	if len(versions) == 0 {
		m.logger.Info("No migrations to roll back")
		return nil
	}

	// roll back each version in descending order
	for _, ver := range versions {
		downName := fmt.Sprintf("%d_*.down.sql", ver)
		pattern := filepath.Join(m.migrationFolder, downName)

		// use glob to find matching .down.sql file
		matches, err := filepath.Glob(pattern)
		if err != nil {
			return err
		}

		if len(matches) == 0 {
			m.logger.Warnf("No .down.sql file found for version %d", ver)
			// still remove from table even if no down file exists
			if err := m.db.Exec("DELETE FROM schema_migrations WHERE version = ?", ver).Error; err != nil {
				return err
			}
			continue
		}

		downPath := matches[0]
		m.logger.Infof("Rolling back migration %d (%s)", ver, downPath)
		bytes, err := os.ReadFile(downPath)
		if err != nil {
			return err
		}
		if err := m.db.Exec(string(bytes)).Error; err != nil {
			return err
		}
		if err := m.db.Exec("DELETE FROM schema_migrations WHERE version = ?", ver).Error; err != nil {
			return err
		}
	}

	m.logger.Info("--- Migration rollback completed successfully ---")
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

// fillGapMigrations applies migrations that exist on disk but are not
// recorded in the schema_migrations table.  It only operates in the
// forward direction and is safe to call multiple times; once a version is
// recorded it will be skipped on subsequent runs.  For downward
// migrations we rely on the underlying migrator's behavior.
func (m Migrations) fillGapMigrations() error {
	if m.db == nil {
		// nothing to do if no connection (e.g. during tests without a
		// database object)
		return nil
	}

	// read applied versions into a set
	applied := make(map[int64]struct{})
	rows, err := m.db.Raw("SELECT version FROM schema_migrations").Rows()
	if err != nil {
		return err
	}
	defer rows.Close()
	for rows.Next() {
		var v int64
		if err := rows.Scan(&v); err != nil {
			return err
		}
		applied[v] = struct{}{}
	}

	// scan migration directory for .up.sql files
	files, err := os.ReadDir(m.migrationFolder)
	if err != nil {
		return err
	}

	type mig struct {
		version int64
		path    string
	}
	var missing []mig

	for _, f := range files {
		if f.IsDir() {
			continue
		}
		name := f.Name()
		if !strings.HasSuffix(name, ".up.sql") {
			continue
		}
		parts := strings.SplitN(name, "_", 2)
		if len(parts) == 0 {
			continue
		}
		v, err := strconv.ParseInt(parts[0], 10, 64)
		if err != nil {
			continue
		}
		if _, ok := applied[v]; !ok {
			missing = append(missing, mig{version: v, path: filepath.Join(m.migrationFolder, name)})
		}
	}

	if len(missing) == 0 {
		return nil
	}

	// apply in version order
	sort.Slice(missing, func(i, j int) bool { return missing[i].version < missing[j].version })
	for _, mfile := range missing {
		m.logger.Infof("applying skipped migration %d (%s)", mfile.version, mfile.path)
		bytes, err := os.ReadFile(mfile.path)
		if err != nil {
			return err
		}
		if err := m.db.Exec(string(bytes)).Error; err != nil {
			return err
		}
		if err := m.db.Exec("INSERT INTO schema_migrations (version, dirty) VALUES (?, 0)", mfile.version).Error; err != nil {
			return err
		}
	}

	return nil
}
