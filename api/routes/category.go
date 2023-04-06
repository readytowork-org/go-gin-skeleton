package routes

import (
	"boilerplate-api/api/controllers"
	"boilerplate-api/infrastructure"
)

type CategoryRoutes struct {
	logger             infrastructure.Logger
	router             infrastructure.Router
	categoryController controllers.CategoryController
}

func (i CategoryRoutes) Setup() {
	i.logger.Zap.Info("settting up category routes")
	category := i.router.Gin.Group("/category")
	{
		category.GET("", i.categoryController.GetAllCategories)
		category.GET(":id", i.categoryController.GetCategory)
		category.POST("", i.categoryController.CreateCategory)
		category.DELETE(":id", i.categoryController.DeleteCategory)
		category.PUT(":id", i.categoryController.UpdateCategory)
	}
}

func NewCategoryRoutes(
	logger infrastructure.Logger,
	router infrastructure.Router,
	categoryController controllers.CategoryController,

) CategoryRoutes {
	return CategoryRoutes{
		logger:             logger,
		router:             router,
		categoryController: categoryController,
	}
}
