package routes

import (
	"boilerplate-api/api/controllers"
	"boilerplate-api/infrastructure"
)

type ProductRoutes struct {
	logger            infrastructure.Logger
	productController controllers.ProductController
	router            infrastructure.Router
}

func NewProductRoutes(
	logger infrastructure.Logger,
	productController controllers.ProductController,
	router infrastructure.Router,
) ProductRoutes {
	return ProductRoutes{
		logger:            logger,
		productController: productController,
		router:            router,
	}
}

func (i ProductRoutes) Setup() {
	i.logger.Zap.Info("Setting Up the Product Routes")
	product := i.router.Gin.Group("/products")
	{
		product.POST("/add", i.productController.AddProducts)
		product.GET("", i.productController.GetAllProducts)
	}
}
