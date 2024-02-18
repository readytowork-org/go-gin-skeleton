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
	oauthMiddleware     middlewares.OAuthMiddleWare
	trxMiddleware       middlewares.DBTransactionMiddleware
	rateLimitMiddleware middlewares.RateLimitMiddleware
}

// NewUserRoutes creates new user controller
func NewUserRoutes(
	logger infrastructure.Logger,
	router infrastructure.Router,
	userController controllers.UserController,
	middleware middlewares.FirebaseAuthMiddleware,
	jwtMiddleware middlewares.JWTAuthMiddleWare,
	oauthMiddleware middlewares.OAuthMiddleWare,
	trxMiddleware middlewares.DBTransactionMiddleware,
	rateLimitMiddleware middlewares.RateLimitMiddleware,
) UserRoutes {
	return UserRoutes{
		router:              router,
		logger:              logger,
		userController:      userController,
		middleware:          middleware,
		jwtMiddleware:       jwtMiddleware,
		oauthMiddleware:     oauthMiddleware,
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
	i.router.Gin.GET("/profile", i.jwtMiddleware.Handle(), i.userController.GetUserProfile)
	// OAuth middleware implemented
	i.router.Gin.GET("/profile/check-oAuth-middleware", i.oauthMiddleware.Handle(), i.userController.GetUserProfile)
	oAuth := i.router.Gin.Group("/oauth")
	{
		oAuth.POST("/sign-in", i.userController.OAuthSignIn)
		oAuth.GET("/callback", i.trxMiddleware.DBTransactionHandle(), i.userController.OAuthCallback)
	}
}
