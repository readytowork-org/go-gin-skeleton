package infrastructure

import (
	"fmt"
	"net/http"

	"github.com/getsentry/sentry-go"
	sentrygin "github.com/getsentry/sentry-go/gin"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

// Router -> Gin Router
type Router struct {
	Gin *gin.Engine
	Env Env
}

//NewRouter : all the routes are defined here
func NewRouter(env Env) Router {

	if env.Environment != "local" {
		if err := sentry.Init(sentry.ClientOptions{
			Dsn:              env.SentryDSN,
			Environment:      `boilerplate-backend-` + env.Environment,
			AttachStacktrace: true,
		}); err != nil {
			fmt.Printf("Sentry initialization failed: %v\n", err)
		}
	}
	if env.Environment == "production" {
		gin.SetMode(gin.ReleaseMode)
	}

	if env.Environment == "production" {
		gin.SetMode(gin.ReleaseMode)
	}

	httpRouter := gin.Default()

	httpRouter.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"PUT", "PATCH", "GET", "POST", "OPTIONS", "DELETE"},
		AllowHeaders:     []string{"*"},
		AllowCredentials: true,
	}))

	httpRouter.Use(sentrygin.New(sentrygin.Options{
		Repanic: true,
	}))

	httpRouter.GET("/health-check", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"data": "Boilerplate 📺 API Up and Running"})
	})

	return Router{
		Gin: httpRouter,
	}
}
