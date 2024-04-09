package routes

import (
	"boilerplate-api/api/controllers"
	"boilerplate-api/external_services/firebase"
	"boilerplate-api/internal/config"
	"boilerplate-api/internal/constants"
	"boilerplate-api/internal/middlewares"
	"boilerplate-api/internal/router"
)

// UserRoutes struct
type UserRoutes struct {
	logger              config.Logger
	router              router.Router
	userController      controllers.UserController
	middleware          firebase.AuthMiddleware
	jwtMiddleware       middlewares.JWTAuthMiddleWare
	trxMiddleware       middlewares.DBTransactionMiddleware
	rateLimitMiddleware middlewares.RateLimitMiddleware
}

// NewUserRoutes creates new user controller
func NewUserRoutes(
	logger config.Logger,
	router router.Router,
	userController controllers.UserController,
	middleware firebase.AuthMiddleware,
	jwtMiddleware middlewares.JWTAuthMiddleWare,
	trxMiddleware middlewares.DBTransactionMiddleware,
	rateLimitMiddleware middlewares.RateLimitMiddleware,
) UserRoutes {
	return UserRoutes{
		router:              router,
		logger:              logger,
		userController:      userController,
		middleware:          middleware,
		jwtMiddleware:       jwtMiddleware,
		trxMiddleware:       trxMiddleware,
		rateLimitMiddleware: rateLimitMiddleware,
	}
}

// Setup user routes
func (i UserRoutes) Setup() {
	i.logger.Info(" Setting up user routes")
	users := i.router.V1.Group("/users").Use(i.rateLimitMiddleware.HandleRateLimit(constants.BasicRateLimit, constants.BasicPeriod))
	{
		users.GET("", i.userController.GetAllUsers)
		users.POST("", i.trxMiddleware.DBTransactionHandle(), i.userController.CreateUser)
	}
	i.router.GET("/profile", i.jwtMiddleware.Handle(), i.userController.GetUserProfile)
}
