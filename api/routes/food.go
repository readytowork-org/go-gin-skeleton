package routes

import (
	"boilerplate-api/api/controllers"
	"boilerplate-api/api/middlewares"
	"boilerplate-api/infrastructure"
)

// FoodRoutes -> struct
type FoodRoutes struct {
	logger                    infrastructure.Logger
	router                    infrastructure.Router
	foodController controllers.FoodController
	middleware                middlewares.FirebaseAuthMiddleware
}

// Setup food routes
func (c FoodRoutes) Setup() {
	c.logger.Zap.Info(" Setting up Food routes")
	food := c.router.Gin.Group("/food")
	{
		food.POST("", c.foodController.CreateFood)
		food.GET("", c.foodController.GetAllFood)
		food.GET("/:id", c.foodController.GetOneFood)
		food.PUT("/:id", c.foodController.UpdateOneFood)
		food.DELETE("/:id", c.foodController.DeleteOneFood)
	}
}

// NewFoodRoutes -> creates new Food controller
func NewFoodRoutes(
	logger infrastructure.Logger,
	router infrastructure.Router,
	foodController controllers.FoodController,
	middleware middlewares.FirebaseAuthMiddleware,
) FoodRoutes {
	return FoodRoutes{
		router:                    router,
		logger:                    logger,
		foodController: foodController,
		middleware:                middleware,
	}
}
