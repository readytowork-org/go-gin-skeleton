package routes

import (
	"boilerplate-api/api/controllers"
	"boilerplate-api/infrastructure"
)

type ProductRoutes struct {
	logger            infrastructure.Logger
	router            infrastructure.Router
	productController controllers.ProductController
}

func (i ProductRoutes) Setup() {
	i.logger.Zap.Info("setting up product routes")
	products := i.router.Gin.Group("/products")
	{
		products.POST("", i.productController.CreateProduct)
		products.GET("", i.productController.GetAllProducts)
		products.GET("/:id", i.productController.GetProduct)
	}
}

func NewProductRoutes(
	logger infrastructure.Logger,
	router infrastructure.Router,
	productController controllers.ProductController,

) ProductRoutes {
	return ProductRoutes{
		router:            router,
		logger:            logger,
		productController: productController,
	}
}
