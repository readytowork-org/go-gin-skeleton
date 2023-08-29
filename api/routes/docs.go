package routes

import (
	"boilerplate-api/infrastructure"

	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

type DocsRoutes struct {
	logger infrastructure.Logger
	router infrastructure.Router
	env    infrastructure.Env
}

func NewDocsRoutes(
	logger infrastructure.Logger,
	router infrastructure.Router,
	env infrastructure.Env,
) DocsRoutes {
	return DocsRoutes{
		router: router,
		logger: logger,
		env:    env,
	}
}

func (c DocsRoutes) Setup() {
	if c.env.Environment != "production" {
		c.logger.Zap.Info(" Setting up Docs routes")

		swagger := c.router.Gin.Group("/docs")
		{
			swagger.GET("/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
		}
	}
}
