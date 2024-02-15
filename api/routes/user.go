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
	jwtMiddleware       middlewares.JWTAuthMiddleWare
	trxMiddleware       middlewares.DBTransactionMiddleware
	rateLimitMiddleware middlewares.RateLimitMiddleware
	redisMiddleware     middlewares.RedisMiddleware
}

// NewUserRoutes creates new user controller
func NewUserRoutes(
	logger infrastructure.Logger,
	router infrastructure.Router,
	userController controllers.UserController,
	middleware middlewares.FirebaseAuthMiddleware,
	jwtMiddleware middlewares.JWTAuthMiddleWare,
	trxMiddleware middlewares.DBTransactionMiddleware,
	rateLimitMiddleware middlewares.RateLimitMiddleware,
	redisMiddleware middlewares.RedisMiddleware,
) UserRoutes {
	return UserRoutes{
		router:              router,
		logger:              logger,
		userController:      userController,
		middleware:          middleware,
		jwtMiddleware:       jwtMiddleware,
		trxMiddleware:       trxMiddleware,
		rateLimitMiddleware: rateLimitMiddleware,
		redisMiddleware:     redisMiddleware,
	}
}

// Setup user routes
func (i UserRoutes) Setup() {
	i.logger.Zap.Info(" Setting up user routes")
	users := i.router.Gin.Group("/users").Use(i.rateLimitMiddleware.HandleRateLimit(constants.BasicRateLimit, constants.BasicPeriod))
	{
		// i.redisMiddleware.VerifyRedisCache()
		users.GET("", i.userController.GetAllUsers)
		users.POST("", i.trxMiddleware.DBTransactionHandle(), i.userController.CreateUser)
	}
	i.router.Gin.GET("/profile", i.jwtMiddleware.Handle(), i.redisMiddleware.VerifyRedisCache(), i.userController.GetUserProfile)
}
