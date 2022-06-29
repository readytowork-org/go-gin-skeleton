package routes

import (
	"boilerplate-api/api/controllers"
	"boilerplate-api/api/middlewares"
	"boilerplate-api/infrastructure"
)

// PostRoutes -> struct
type PostRoutes struct {
	logger                    infrastructure.Logger
	router                    infrastructure.Router
	postController controllers.PostController
	middleware                middlewares.FirebaseAuthMiddleware
}

// NewPostRoutes -> creates new Post controller
func NewPostRoutes(
	logger infrastructure.Logger,
	router infrastructure.Router,
	postController controllers.PostController,
	middleware middlewares.FirebaseAuthMiddleware,
) PostRoutes {
	return PostRoutes{
		router:                    router,
		logger:                    logger,
		postController: postController,
		middleware:                middleware,
	}
}

// Setup post routes
func (c PostRoutes) Setup() {
	c.logger.Zap.Info(" Setting up Post routes")
	post := c.router.Gin.Group("/post")
	{
		post.POST("", c.postController.CreatePost)
		post.GET("", c.postController.GetAllPost)
		post.GET("/:id", c.postController.GetOnePost)
		post.PUT("/:id", c.postController.UpdateOnePost)
		post.DELETE("/:id", c.postController.DeleteOnePost)
	}
}
