package user_router

import (
	"boilerplate-api/apps/user"
	"boilerplate-api/infrastructure"
	"boilerplate-api/middlewares"
)

type Route struct {
	logger     infrastructure.Logger
	router     infrastructure.Router
	controller user.Controller
	trx        middlewares.DBTransaction
}

func RouteConstructor(
	logger infrastructure.Logger,
	router infrastructure.Router,
	controller user.Controller,
	trx middlewares.DBTransaction,
) Route {
	return Route{
		router:     router,
		logger:     logger,
		controller: controller,
		trx:        trx,
	}
}

func (i Route) Setup() {
	i.logger.Zap.Info("->Setting up user routes<-")

	users := i.router.Gin.Group("/users")
	{
		users.GET("", i.controller.GetAllUsers)
		users.POST("", i.trx.DBTransactionHandle(), i.controller.CreateUser)
	}
	i.router.Gin.GET("/profile", i.controller.GetUserProfile)
}
