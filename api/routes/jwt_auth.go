package routes

import (
	"boilerplate-api/api/controllers"
	"boilerplate-api/internal/config"
	"boilerplate-api/internal/constants"
	"boilerplate-api/internal/middlewares"
	"boilerplate-api/internal/router"
)

// JwtAuthRoutes struct
type JwtAuthRoutes struct {
	logger              config.Logger
	router              router.Router
	jwtController       controllers.JwtAuthController
	rateLimitMiddleware middlewares.RateLimitMiddleware
}

// NewJwtAuthRoutes creates new jwt controller
func NewJwtAuthRoutes(
	logger config.Logger,
	router router.Router,
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
	i.logger.Info(" Setting up jwt routes")
	jwt := i.router.Group("/login").Use(i.rateLimitMiddleware.HandleRateLimit(constants.LoginRateLimit, constants.LoginPeriod))
	{
		jwt.POST("", i.jwtController.LoginUserWithJWT)
		jwt.POST("/refresh", i.jwtController.RefreshJwtToken)
	}
}
