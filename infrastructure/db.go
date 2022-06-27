package infrastructure

import (
	"fmt"
	"log"
	"os"
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// Database modal
type Database struct {
	DB *gorm.DB
}

// NewDatabase creates a new database instance
func NewDatabase(Zaplogger Logger, env Env) Database {

	newLogger := logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags), // io writer
		logger.Config{
			SlowThreshold: time.Second, // Slow SQL threshold
			LogLevel:      logger.Info, // Log level
			Colorful:      true,        // Disable color
		},
	)

	Zaplogger.Zap.Info(env)

	url := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local", env.DBUsername, env.DBPassword, env.DBHost, env.DBPort, env.DBName)

	if env.Environment != "local" {
		url = fmt.Sprintf(
			"%s:%s@unix(/cloudsql/%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
			env.DBUsername,
			env.DBPassword,
			env.DBHost,
			env.DBName,
		)
	}

	db, err := gorm.Open(mysql.Open(url), &gorm.Config{Logger: newLogger})
	_ = db.Exec("CREATE DATABASE IF NOT EXISTS " + env.DBName + ";")
	if err != nil {
		Zaplogger.Zap.Info("Url: ", url)
		Zaplogger.Zap.Panic(err)
	}

	Zaplogger.Zap.Info("using given database")
	if err := db.Exec(fmt.Sprintf("USE %s", env.DBName)).Error; err != nil {
		Zaplogger.Zap.Info("cannot use the given database")
		Zaplogger.Zap.Panic(err)
	}
	database := Database{DB: db}
	Zaplogger.Zap.Info("Database connection established")
	if err := RunMigration(Zaplogger, database); err != nil {
		Zaplogger.Zap.Info("migration failed.")
		Zaplogger.Zap.Panic(err)
	}
	return Database{
		DB: db,
	}
}
