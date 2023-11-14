package auth_router

import (
	"boilerplate-api/apps/auth"
	"boilerplate-api/infrastructure"
)

// JwtAuthRoutes struct
type Route struct {
	logger     infrastructure.Logger
	router     infrastructure.Router
	controller auth.Controller
}

// NewJwtAuthRoutes creates new jwt controller
func RouteConstructor(
	logger infrastructure.Logger,
	router infrastructure.Router,
	controller auth.Controller,
) Route {
	return Route{
		router:     router,
		logger:     logger,
		controller: controller,
	}
}

// Setup Obtain Jwt Token Routes
func (i Route) Setup() {
	i.logger.Zap.Info(" Setting up jwt routes")
	jwt := i.router.Gin.Group("/login")
	{
		jwt.POST("", i.controller.LoginUser)
		jwt.POST("/refresh", i.controller.RefreshToken)
	}
}
