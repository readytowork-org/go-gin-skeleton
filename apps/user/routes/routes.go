package routes

import (
	"boilerplate-api/apps/user/controllers"
	"boilerplate-api/infrastructure"
	"boilerplate-api/middlewares"
)

type UserRoute struct {
	logger     infrastructure.Logger
	router     infrastructure.Router
	controller controllers.UserController
	trx        middlewares.DBTransaction
}

func UserRouteConstructor(
	logger infrastructure.Logger,
	router infrastructure.Router,
	controller controllers.UserController,
	trx middlewares.DBTransaction,
) UserRoute {
	return UserRoute{
		router:     router,
		logger:     logger,
		controller: controller,
		trx:        trx,
	}
}

func (i UserRoute) Setup() {
	i.logger.Zap.Info("->Setting up user routes<-")

	users := i.router.Gin.Group("/users")
	{
		users.GET("", i.controller.GetAllUsers)
		users.POST("", i.trx.DBTransactionHandle(), i.controller.CreateUser)
	}
	i.router.Gin.GET("/profile", i.controller.GetUserProfile)
}
