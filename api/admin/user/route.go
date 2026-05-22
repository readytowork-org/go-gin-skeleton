package user

import (
	"boilerplate-api/lib/config"
	"boilerplate-api/lib/constants"
	"boilerplate-api/lib/middlewares"
	"boilerplate-api/lib/router"
)

// SetupRoutes user routes
func SetupRoutes(
	logger config.Logger,
	router router.Router,
	userController Controller,
	jwtMiddleware middlewares.JWTAuthMiddleWare,
	trxMiddleware middlewares.DBTransactionMiddleware,
	rateLimitMiddleware middlewares.RateLimitMiddleware,
	idempotencyMiddleware middlewares.IdempotencyMiddleware,
) {
	logger.Info(" Setting up user routes")
	users := router.V1.Group("/users").
		Use(rateLimitMiddleware.HandleRateLimit(constants.BasicRateLimit, constants.BasicPeriod)).
		Use(jwtMiddleware.Handle())
	{
		users.GET("", userController.GetAllUsers)
		users.POST("",
			idempotencyMiddleware.Handle(),
			trxMiddleware.DBTransactionHandle(),
			userController.CreateUser,
		)
		users.GET("/:id", userController.GetOneUser)
	}
}
