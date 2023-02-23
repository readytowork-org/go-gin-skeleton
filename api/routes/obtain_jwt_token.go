package routes

import (
	"boilerplate-api/api/controllers"
	"boilerplate-api/api/middlewares"
	"boilerplate-api/constants"
	"boilerplate-api/infrastructure"
)

// ObtainJwtTokenRoutes -> struct
type ObtainJwtTokenRoutes struct {
	logger        infrastructure.Logger
	router        infrastructure.Router
	jwtController controllers.JwtAuthController
	rateLimitMiddleware middlewares.RateLimitMiddleware
}

// Setup Obtain Jwt Token Routes
func (i ObtainJwtTokenRoutes) Setup() {
	i.logger.Zap.Info(" Setting up jwt routes")
	jwt := i.router.Gin.Group("/login").Use(i.rateLimitMiddleware.HandleRateLimit(constants.LoginRateLimit,constants.LoginPeriod))
	{
		jwt.POST("", i.jwtController.ObtainJwtToken)
		jwt.POST("/refresh", i.jwtController.RefreshJwtToken)
	}
}

// NewObtainJwtTokenRoutes -> creates new jwt controller
func NewObtainJwtTokenRoutes(
	logger infrastructure.Logger,
	router infrastructure.Router,
	jwtController controllers.JwtAuthController,
	rateLimitMiddleware middlewares.RateLimitMiddleware,
) ObtainJwtTokenRoutes {
	return ObtainJwtTokenRoutes{
		router:        router,
		logger:        logger,
		jwtController: jwtController,
		rateLimitMiddleware: rateLimitMiddleware,
	}
}
