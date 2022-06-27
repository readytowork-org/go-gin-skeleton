package infrastructure

import (
	_ "github.com/golang-migrate/migrate/v4/database/mysql"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	migrate "github.com/rubenv/sql-migrate"
)

//Migrations -> Migration Struct
type Migrations struct {
	logger Logger
	db     Database
}

//NewMigrations -> return new Migrations struct
func NewMigrations(
	logger Logger,
	db Database,
) Migrations {
	return Migrations{
		logger: logger,
		db:     db,
	}
}

//Migrate -> migrates all table
func (m Migrations) Migrate() error {
	m.logger.Zap.Info("Migrating schemas...")

	migrations := &migrate.FileMigrationSource{
		Dir: "migration/",
	}

	sqlDB, err := m.db.DB.DB()
	if err != nil {
		return err
	}

	m.logger.Zap.Info("running migration.")
	_, err = migrate.Exec(sqlDB, "mysql", migrations, migrate.Up)
	if err != nil {
		return err
	}
	m.logger.Zap.Info("migration completed.")
	return nil

}

// RunMigration runs the migration provided logger and database instance
func RunMigration(logger Logger, db Database) error {
	m := Migrations{logger, db}
	return m.Migrate()
}
