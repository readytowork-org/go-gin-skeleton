package api

import (
	"boilerplate-api/app/global/infrastructure"
	"boilerplate-api/app/global/middlewares"
	"boilerplate-api/app/packages/user"

	"github.com/gin-gonic/gin"
)

// UserRoutes -> struct
type UserRoutes struct {
	logger         infrastructure.Logger
	userController user.Controller
	middleware     middlewares.FirebaseAuthMiddleware
	trxMiddleware  middlewares.DBTransactionMiddleware
}

// Setup user routes
func (i UserRoutes) Setup(v1 *gin.RouterGroup) {
	i.logger.Zap.Info(" Setting up user routes")
	users := v1.Group("/users")
	{
		users.GET("", i.userController.GetAllUsers)
		users.POST("", i.trxMiddleware.DBTransactionHandle(), i.userController.CreateUser)
	}
}

// NewUserRoutes -> creates new user controller
func NewUserRoutes(
	logger infrastructure.Logger,
	userController user.Controller,
	middleware middlewares.FirebaseAuthMiddleware,
	trxMiddleware middlewares.DBTransactionMiddleware,
) UserRoutes {
	return UserRoutes{
		logger:         logger,
		userController: userController,
		middleware:     middleware,
		trxMiddleware:  trxMiddleware,
	}
}
