package middlewares

import (
	"boilerplate-api/constants"
	"boilerplate-api/infrastructure"
	"boilerplate-api/utils"
	"net/http"

	"github.com/gin-gonic/gin"
)

type DBTransaction struct {
	logger infrastructure.Logger
	db     infrastructure.Database
}

// NewDBTransaction new instance of transaction
func NewDBTransaction(
	logger infrastructure.Logger,
	db infrastructure.Database,
) DBTransaction {
	return DBTransaction{
		logger: logger,
		db:     db,
	}
}

// DBTransactionHandle It setup the database transaction middleware
func (m DBTransaction) DBTransactionHandle() gin.HandlerFunc {
	m.logger.Zap.Info("setting up database transaction middleware")

	return func(c *gin.Context) {
		txHandle := m.db.DB.Begin()
		m.logger.Zap.Info("beginning database transaction")

		defer func() {
			if r := recover(); r != nil {
				if err := txHandle.Error; err != nil {
					m.logger.Zap.Error("trx commit error: ", err)
				}
				txHandle.Rollback()
			}
		}()

		c.Set(constants.DBTransaction, txHandle)
		c.Next()

		if utils.StatusInList(c.Writer.Status(), []int{http.StatusOK, http.StatusCreated}) {
			m.logger.Zap.Info("committing transactions")
			if err := txHandle.Commit().Error; err != nil {
				m.logger.Zap.Error("trx commit error: ", err)
			}
		} else {
			m.logger.Zap.Info("rolling back transaction due to status code: ", c.Writer.Status())
			txHandle.Rollback()
		}
	}
}
