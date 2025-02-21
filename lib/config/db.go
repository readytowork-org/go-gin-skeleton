package config

import (
	"errors"
	"fmt"
	"strings"

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
func NewDatabase(logger Logger, dbDialect DBDialect) *Database {
	database := &Database{}
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
	}

	logger.Infof("Database connection established : %s", db.Migrator().CurrentDatabase())

	return database
}

func (d *Database) ConnectionError() error {
	return d.connectionError
}
