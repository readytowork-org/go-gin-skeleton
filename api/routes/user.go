package routes

import (
	"boilerplate-api/api/controllers"
	"boilerplate-api/api/middlewares"
	"boilerplate-api/constants"
	"boilerplate-api/infrastructure"
)

// UserRoutes -> struct
type UserRoutes struct {
	logger         infrastructure.Logger
	router         infrastructure.Router
	userController controllers.UserController
	middleware     middlewares.FirebaseAuthMiddleware
	jwtMiddleware  middlewares.JWTAuthMiddleWare
	trxMiddleware  middlewares.DBTransactionMiddleware
	rateLimitMiddleware middlewares.RateLimitMiddleware
}

// Setup user routes
func (i UserRoutes) Setup() {
	i.logger.Zap.Info(" Setting up user routes")
	users := i.router.Gin.Group("/users").Use(i.rateLimitMiddleware.HandleRateLimit(constants.BasicRateLimit,constants.BasicPeriod))
	{
		users.GET("", i.userController.GetAllUsers)
		users.POST("", i.trxMiddleware.DBTransactionHandle(), i.userController.CreateUser)
	}
	i.router.Gin.GET("/profile",i.jwtMiddleware.Handle(), i.userController.GetUserProfile)
}

// NewUserRoutes -> creates new user controller
func NewUserRoutes(
	logger infrastructure.Logger,
	router infrastructure.Router,
	userController controllers.UserController,
	middleware middlewares.FirebaseAuthMiddleware,
	jwtMiddleware middlewares.JWTAuthMiddleWare,
	trxMiddleware middlewares.DBTransactionMiddleware,
	rateLimitMiddleware middlewares.RateLimitMiddleware,
) UserRoutes {
	return UserRoutes{
		router:         router,
		logger:         logger,
		userController: userController,
		middleware:     middleware,
		jwtMiddleware:  jwtMiddleware,
		trxMiddleware:  trxMiddleware,
		rateLimitMiddleware: rateLimitMiddleware,
	}
}
