package routes

import (
	"boilerplate-api/api/controllers"
	"boilerplate-api/api/middlewares"
	"boilerplate-api/infrastructure"
)

// UserRoutes -> struct
type ProductRoutes struct {
	logger            infrastructure.Logger
	router            infrastructure.Router
	productController controllers.ProductController
	middleware        middlewares.FirebaseAuthMiddleware
	trxMiddleware     middlewares.DBTransactionMiddleware
}

// Setup user routes
func (i ProductRoutes) Setup() {
	i.logger.Zap.Info(" Setting up user routes")
	products := i.router.Gin.Group("/products")
	{
		products.GET("", i.productController.GetAllProducts)
		products.POST("", i.productController.AddProducts)
		products.GET("/filter/user/:id", i.productController.FilterUserProducts)
	}
}

// NewProductRoutes -> creates new user controller
func NewProductRoutes(
	logger infrastructure.Logger,
	router infrastructure.Router,
	productController controllers.ProductController,
	middleware middlewares.FirebaseAuthMiddleware,
	trxMiddleware middlewares.DBTransactionMiddleware,
) ProductRoutes {
	return ProductRoutes{
		router:            router,
		logger:            logger,
		productController: productController,
		middleware:        middleware,
		trxMiddleware:     trxMiddleware,
	}
}
