package config

import (
	"fmt"
	"time"

	goMySql "github.com/go-sql-driver/mysql"
	"gorm.io/driver/mysql"
)

func NewSqlDialect(env Env) DBDialect {
	net := "tcp"
	address := fmt.Sprintf("%s:%s", env.DBHost, env.DBPort)
	if env.Environment == "development" || env.Environment == "production" {
		net = "unix"
		address = fmt.Sprintf("/cloudsql/%s", env.DBHost)
	}

	location, _ := time.LoadLocation(env.TimeZone)

	mysqlDSNConfig := goMySql.Config{
		User:                 env.DBUsername,
		Passwd:               env.DBPassword,
		Net:                  net,
		Addr:                 address,
		DBName:               env.DBName,
		ParseTime:            true,
		Loc:                  location,
		AllowNativePasswords: true,
		CheckConnLiveness:    true,
	}

	return DBDialect{
		Dialector: mysql.New(
			mysql.Config{
				DSNConfig: &mysqlDSNConfig,
			},
		),
		DBName: env.DBName,
		DSN:    mysqlDSNConfig.FormatDSN(),
	}
}
