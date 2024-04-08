package routes

import (
	"boilerplate-api/internal/config"
	"boilerplate-api/internal/router"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

type DocsRoutes struct {
	logger config.Logger
	router router.Router
	env    config.Env
}

func NewDocsRoutes(
	logger config.Logger,
	router router.Router,
	env config.Env,
) DocsRoutes {
	return DocsRoutes{
		router: router,
		logger: logger,
		env:    env,
	}
}

func (c DocsRoutes) Setup() {
	if c.env.Environment != "production" {
		c.logger.Info(" Setting up Docs routes")

		swagger := c.router.Group("/docs")
		{
			swagger.GET("/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
		}
	}
}
