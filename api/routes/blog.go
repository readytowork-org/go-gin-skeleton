package routes

import (
	"boilerplate-api/api/controllers"
	"boilerplate-api/api/middlewares"
	"boilerplate-api/infrastructure"
)

// UserRoutes -> struct
type BlogRoutes struct {
	logger         infrastructure.Logger
	router         infrastructure.Router
	blogController controllers.BlogController
	middleware     middlewares.FirebaseAuthMiddleware
	trxMiddleware  middlewares.DBTransactionMiddleware
}

// Setup user routes
func (i BlogRoutes) Setup() {
	i.logger.Zap.Info(" Setting up user routes")
	blog := i.router.Gin.Group("/blog")
	{
		blog.GET("", i.blogController.GetAllBlogs)
		// blog.GET("/:id", i.blogController.GetOneCategory)
		// blog.PUT("/:id", i.trxMiddleware.DBTransactionHandle(), i.middleware.Handle(), i.blogController.UpdateOneCategory)
		// // blog.PATCH("/:id", i.blogController.PartialUpdateUser)
		// blog.DELETE("/:id", i.trxMiddleware.DBTransactionHandle(), i.middleware.Handle(), i.blogController.DeleteOneCategory)
		// blog.POST("", i.trxMiddleware.DBTransactionHandle(), i.middleware.Handle(), i.blogController.CreateCategory)
	}
}

// NewUserRoutes -> creates new user controller
func NewBlogRoutes(
	logger infrastructure.Logger,
	router infrastructure.Router,
	blogController controllers.BlogController,
	middleware middlewares.FirebaseAuthMiddleware,
	trxMiddleware middlewares.DBTransactionMiddleware,
) BlogRoutes {
	return BlogRoutes{
		router:         router,
		logger:         logger,
		blogController: blogController,
		middleware:     middleware,
		trxMiddleware:  trxMiddleware,
	}
}
