package routes

import (
	_ "boilerplate-api/docs"
	"boilerplate-api/infrastructure"

	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

type DocsRoutes struct {
	logger infrastructure.Logger
	router infrastructure.Router
}

func NewDocsRoutes(
	logger infrastructure.Logger,
	router infrastructure.Router,
) DocsRoutes {
	return DocsRoutes{
		router: router,
		logger: logger,
	}
}

func (c DocsRoutes) Setup() {

	c.logger.Zap.Info(" Setting up Docs routes")

	swagger := c.router.Gin.Group("/docs")
	{
		swagger.GET("/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	}
}
