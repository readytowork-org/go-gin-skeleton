package config

import (
	"errors"
	"fmt"
	"strings"
	"time"

	"gorm.io/gorm"
)

type DBDialect struct {
	gorm.Dialector
	DSN    string
	DBName string
}

// Database modal
type Database struct {
	*gorm.DB
	connectionError error
}

// NewDatabase creates a new database instance
func NewDatabase(logger Logger, dbDialect DBDialect, env Env) Database {
	database := Database{}
	db, err := gorm.Open(dbDialect, &gorm.Config{Logger: logger.GetGormLogger()})
	if err != nil {
		_err := errors.New(
			fmt.Sprintf(
				"Database connection failed\n Please check dsn:: %+v\n+%v", dbDialect.DSN, err.Error(),
			),
		)
		logger.Error(_err)
		database.connectionError = _err
		return database
	}
	database.DB = db

	// not needed in sqlite
	if !strings.Contains(dbDialect.Dialector.Name(), "sqlite") {
		logger.Info("using given database")
		if err = db.Exec(fmt.Sprintf("USE %s", dbDialect.DBName)).Error; err != nil {
			logger.Error("Cannot use the given database")
			database.connectionError = err
			return database
		}

		if sqlDB, sqlErr := db.DB(); sqlErr == nil {
			applyPoolSettings(sqlDB, env, logger)
		} else {
			logger.Warn("could not access underlying *sql.DB to tune pool: ", sqlErr.Error())
		}
	}

	logger.Infof("Database connection established : %s", db.Migrator().CurrentDatabase())

	return database
}

// applyPoolSettings tunes the connection pool from Env, with safe defaults.
// GORM's default pool is effectively unbounded, which behind Cloud SQL proxy
// or a small DB can exhaust connections quickly under load.
func applyPoolSettings(sqlDB interface {
	SetMaxOpenConns(int)
	SetMaxIdleConns(int)
	SetConnMaxLifetime(time.Duration)
}, env Env, logger Logger) {
	maxOpen := env.DBMaxOpenConns
	if maxOpen <= 0 {
		maxOpen = 25
	}
	maxIdle := env.DBMaxIdleConns
	if maxIdle <= 0 {
		maxIdle = 5
	}
	lifetime := env.DBConnMaxLifetime
	if lifetime <= 0 {
		lifetime = 30 * time.Minute
	}

	sqlDB.SetMaxOpenConns(maxOpen)
	sqlDB.SetMaxIdleConns(maxIdle)
	sqlDB.SetConnMaxLifetime(lifetime)
	logger.Infof("DB pool tuned: maxOpen=%d maxIdle=%d connMaxLifetime=%s", maxOpen, maxIdle, lifetime)
}

func (d *Database) ConnectionError() error {
	return d.connectionError
}
