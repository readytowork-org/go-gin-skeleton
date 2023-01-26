package web

import (
	"boilerplate-api/app/global/infrastructure"
	"net/http"

	"github.com/gin-gonic/gin"
)

// HomeRoutes -> struct
type HomeRoutes struct {
	logger infrastructure.Logger
}

// Setup user routes
func (i HomeRoutes) Setup(route *gin.RouterGroup) {
	i.logger.Zap.Info(" Setting up user routes")
	route.GET("/", func(c *gin.Context) {
		c.Data(http.StatusOK, "text/html; charset=utf-8", []byte(`<html>working!</html>`))
		return
	})
}

// NewHomeRoutes -> creates new user controller
func NewHomeRoutes(
	logger infrastructure.Logger,
) HomeRoutes {
	return HomeRoutes{
		logger: logger,
	}
}
