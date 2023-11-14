package userrouter

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
	jwt        middlewares.JWTAuthMiddleWare
}

func RouteConstructor(
	logger infrastructure.Logger,
	router infrastructure.Router,
	controller user.Controller,
	trx middlewares.DBTransaction,
	jwt middlewares.JWTAuthMiddleWare,
) Route {
	return Route{
		router:     router,
		logger:     logger,
		controller: controller,
		trx:        trx,
		jwt:        jwt,
	}
}

func (i Route) Setup() {
	i.logger.Zap.Info("->Setting up user routes<-")

	users := i.router.Gin.Group("/users")
	{
		users.GET("", i.controller.GetAllUsers)
		users.POST("", i.trx.DBTransactionHandle(), i.controller.CreateUser)
	}
	i.router.Gin.GET("/profile", i.jwt.Handle(), i.controller.GetUserProfile)
}
