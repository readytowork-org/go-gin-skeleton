package routes

import (
	"boilerplate-api/api/controllers"
	"boilerplate-api/api/middlewares"
	"boilerplate-api/infrastructure"
)

// FruitRoutes -> struct
type FruitRoutes struct {
	logger                    infrastructure.Logger
	router                    infrastructure.Router
	fruitController controllers.FruitController
	middleware                middlewares.FirebaseAuthMiddleware
}

// NewFruitRoutes -> creates new Fruit controller
func NewFruitRoutes(
	logger infrastructure.Logger,
	router infrastructure.Router,
	fruitController controllers.FruitController,
	middleware middlewares.FirebaseAuthMiddleware,
) FruitRoutes {
	return FruitRoutes{
		router:                    router,
		logger:                    logger,
		fruitController: fruitController,
		middleware:                middleware,
	}
}

// Setup fruit routes
func (c FruitRoutes) Setup() {
	c.logger.Zap.Info(" Setting up Fruit routes")
	fruit := c.router.Gin.Group("/fruit")
	{
		fruit.POST("", c.fruitController.CreateFruit)
		fruit.GET("", c.fruitController.GetAllFruit)
		fruit.GET("/:id", c.fruitController.GetOneFruit)
		fruit.PUT("/:id", c.fruitController.UpdateOneFruit)
		fruit.DELETE("/:id", c.fruitController.DeleteOneFruit)
	}
}
