package routes

import (
	"boilerplate-api/api/controllers"
	"boilerplate-api/api/middlewares"
	"boilerplate-api/infrastructure"
)

// UserRoutes -> struct
type CategoryRoutes struct {
	logger             infrastructure.Logger
	router             infrastructure.Router
	categoryController controllers.CategoryController
	middleware         middlewares.FirebaseAuthMiddleware
	trxMiddleware      middlewares.DBTransactionMiddleware
}

// Setup user routes
func (i CategoryRoutes) Setup() {
	i.logger.Zap.Info(" Setting up user routes")
	category := i.router.Gin.Group("/category")
	{
		category.GET("", i.categoryController.GetAllCategory)
		category.GET("/:id", i.categoryController.GetOneCategory)
		category.PUT("/:id", i.trxMiddleware.DBTransactionHandle(), i.middleware.Handle(), i.categoryController.UpdateOneCategory)
		// category.PATCH("/:id", i.categoryController.PartialUpdateUser)
		category.DELETE("/:id", i.trxMiddleware.DBTransactionHandle(), i.middleware.Handle(), i.categoryController.DeleteOneCategory)
		category.POST("", i.trxMiddleware.DBTransactionHandle(), i.middleware.Handle(), i.categoryController.CreateCategory)
	}
}

// NewUserRoutes -> creates new user controller
func NewCategoryRoutes(
	logger infrastructure.Logger,
	router infrastructure.Router,
	categoryController controllers.CategoryController,
	middleware middlewares.FirebaseAuthMiddleware,
	trxMiddleware middlewares.DBTransactionMiddleware,
) CategoryRoutes {
	return CategoryRoutes{
		router:             router,
		logger:             logger,
		categoryController: categoryController,
		middleware:         middleware,
		trxMiddleware:      trxMiddleware,
	}
}
