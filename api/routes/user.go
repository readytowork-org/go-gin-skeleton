package routes

import (
	"boilerplate-api/api/controllers"
	"boilerplate-api/api/middlewares"
	"boilerplate-api/constants"
	"boilerplate-api/infrastructure"
)

// UserRoutes struct
type UserRoutes struct {
	logger              infrastructure.Logger
	router              infrastructure.Router
	userController      controllers.UserController
	middleware          middlewares.FirebaseAuthMiddleware
	trxMiddleware       middlewares.DBTransactionMiddleware
	rateLimitMiddleware middlewares.RateLimitMiddleware
}

// NewUserRoutes creates new user controller
func NewUserRoutes(
	logger infrastructure.Logger,
	router infrastructure.Router,
	userController controllers.UserController,
	middleware middlewares.FirebaseAuthMiddleware,
	trxMiddleware middlewares.DBTransactionMiddleware,
	rateLimitMiddleware middlewares.RateLimitMiddleware,
) UserRoutes {
	return UserRoutes{
		router:              router,
		logger:              logger,
		userController:      userController,
		middleware:          middleware,
		trxMiddleware:       trxMiddleware,
		rateLimitMiddleware: rateLimitMiddleware,
	}
}

// Setup user routes
func (i UserRoutes) Setup() {
	i.logger.Zap.Info(" Setting up user routes")
	users := i.router.Gin.Group("/users").Use(i.rateLimitMiddleware.HandleRateLimit(constants.BasicRateLimit, constants.BasicPeriod))
	{
		users.GET("", i.userController.GetAllUsers)
		users.POST("", i.trxMiddleware.DBTransactionHandle(), i.userController.CreateUser)
	}
}
