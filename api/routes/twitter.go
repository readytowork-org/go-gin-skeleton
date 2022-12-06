package routes

import (
	"boilerplate-api/api/controllers"
	"boilerplate-api/api/middlewares"
	"boilerplate-api/infrastructure"
)

// UserRoutes -> struct
type TwitterRoutes struct {
	logger            infrastructure.Logger
	router            infrastructure.Router
	twitterController controllers.TwitterController
	middleware        middlewares.FirebaseAuthMiddleware
}

// NewProductRoutes -> creates new user controller
func NewTwitterRoutes(
	logger infrastructure.Logger,
	router infrastructure.Router,
	tweetController controllers.TwitterController,
	middleware middlewares.FirebaseAuthMiddleware,
) TwitterRoutes {
	return TwitterRoutes{
		router:            router,
		logger:            logger,
		twitterController: tweetController,
		middleware:        middleware,
	}
}

// Setup user routes
func (i TwitterRoutes) Setup() {
	i.logger.Zap.Info(" Setting up user routes")
	products := i.router.Gin.Group("/twitter")
	{
		products.GET("/create", i.twitterController.CreateTweet)
	}
}
