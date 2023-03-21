package routes

import (
	"boilerplate-api/api/controllers"
	"boilerplate-api/infrastructure"
)

// ObtainJwtTokenRoutes -> struct
type ObtainJwtTokenRoutes struct {
	logger        infrastructure.Logger
	router        infrastructure.Router
	jwtController controllers.JwtAuthController
}

// Setup Obtain Jwt Token Routes
func (i ObtainJwtTokenRoutes) Setup() {
	i.logger.Zap.Info(" Setting up jwt routes")
	jwt := i.router.Gin.Group("/login")
	{
		jwt.POST("", i.jwtController.LoginUserWithJWT)
		jwt.POST("/refresh", i.jwtController.RefreshJwtToken)
	}
}

// NewObtainJwtTokenRoutes -> creates new jwt controller
func NewObtainJwtTokenRoutes(
	logger infrastructure.Logger,
	router infrastructure.Router,
	jwtController controllers.JwtAuthController,

) ObtainJwtTokenRoutes {
	return ObtainJwtTokenRoutes{
		router:        router,
		logger:        logger,
		jwtController: jwtController,
	}
}
