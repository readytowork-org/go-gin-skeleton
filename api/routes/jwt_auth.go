package routes

import (
	"boilerplate-api/api/controllers"
	"boilerplate-api/api/middlewares"
	"boilerplate-api/constants"
	"boilerplate-api/infrastructure"
)

// JwtAuthRoutes struct
type JwtAuthRoutes struct {
	logger              infrastructure.Logger
	router              infrastructure.Router
	jwtController       controllers.JwtAuthController
	rateLimitMiddleware middlewares.RateLimitMiddleware
}

// NewJwtAuthRoutes creates new jwt controller
func NewJwtAuthRoutes(
	logger infrastructure.Logger,
	router infrastructure.Router,
	jwtController controllers.JwtAuthController,
	rateLimitMiddleware middlewares.RateLimitMiddleware,
) JwtAuthRoutes {
	return JwtAuthRoutes{
		router:              router,
		logger:              logger,
		jwtController:       jwtController,
		rateLimitMiddleware: rateLimitMiddleware,
	}
}

// Setup Obtain Jwt Token Routes
func (i JwtAuthRoutes) Setup() {
	i.logger.Zap.Info(" Setting up jwt routes")
	jwt := i.router.Gin.Group("/login").Use(i.rateLimitMiddleware.HandleRateLimit(constants.LoginRateLimit, constants.LoginPeriod))
	{
		jwt.POST("", i.jwtController.LoginUserWithJWT)
		jwt.POST("/refresh", i.jwtController.RefreshJwtToken)
	}
}
