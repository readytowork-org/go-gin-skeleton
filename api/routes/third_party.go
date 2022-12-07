package routes

import (
	"boilerplate-api/api/controllers"
	"boilerplate-api/api/middlewares"
	"boilerplate-api/infrastructure"
)

// UserRoutes -> struct
type ThirdPartyRoutes struct {
	logger               infrastructure.Logger
	router               infrastructure.Router
	thirdPartyController controllers.ThirdPartyController
	middleware           middlewares.FirebaseAuthMiddleware
}

// NewProductRoutes -> creates new user controller
func NewThirdPartyRoutes(
	logger infrastructure.Logger,
	router infrastructure.Router,
	thirdPartyController controllers.ThirdPartyController,
	middleware middlewares.FirebaseAuthMiddleware,
) ThirdPartyRoutes {
	return ThirdPartyRoutes{
		router:               router,
		logger:               logger,
		thirdPartyController: thirdPartyController,
		middleware:           middleware,
	}
}

// Setup user routes
func (i ThirdPartyRoutes) Setup() {
	i.logger.Zap.Info(" Setting up user routes")
	products := i.router.Gin.Group("/merchants")
	{
		products.POST("/create", i.thirdPartyController.MerchanntRegister)
	}
}
