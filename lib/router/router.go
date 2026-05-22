package router

import (
	"context"
	"fmt"
	"net/http"
	"strings"
	"time"

	"boilerplate-api/lib/config"

	"github.com/getsentry/sentry-go"
	sentrygin "github.com/getsentry/sentry-go/gin"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

// Router Gin Router
type Router struct {
	*gin.Engine
	V1 *gin.RouterGroup
}

// NewRouter : all the routes are defined here
func NewRouter(env config.Env, logger config.Logger) Router {
	appEnv := env.Environment

	if appEnv != "local" {
		if err := sentry.Init(
			sentry.ClientOptions{
				Dsn:              env.SentryDSN,
				Environment:      `Demo-backend-` + env.Environment,
				AttachStacktrace: true,
			},
		); err != nil {
			fmt.Printf("Sentry initialization failed: %v\n", err)
		}
	}

	// TODO :: after cli config
	// gin.DefaultWriter = logger.GetGinLogger()

	if appEnv == "production" {
		gin.SetMode(gin.ReleaseMode)
	} else {
		gin.SetMode(gin.DebugMode)
	}

	httpRouter := gin.Default()

	httpRouter.Use(cors.New(buildCorsConfig(env)))

	httpRouter.Use(
		sentrygin.New(
			sentrygin.Options{
				Repanic: true,
			},
		),
	)

	httpRouter.GET(
		"/health-check", func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{"data": " 📺 API Up and Running"})
		},
	)

	api := httpRouter.Group("/api")
	v1 := api.Group("/v1")

	return Router{
		Engine: httpRouter,
		V1:     v1,
	}
}

// buildCorsConfig produces a CORS policy from CORS_ALLOWED_ORIGINS.
// A wildcard cannot be combined with credentials per the CORS spec, so when no
// origins are configured we fall back to AllowAllOrigins without credentials.
func buildCorsConfig(env config.Env) cors.Config {
	cfg := cors.Config{
		AllowMethods: []string{"PUT", "PATCH", "GET", "POST", "OPTIONS", "DELETE"},
		AllowHeaders: []string{"Origin", "Content-Type", "Accept", "Authorization", "X-Requested-With"},
	}

	origins := splitAndTrim(env.CorsAllowedOrigins)
	if len(origins) == 0 {
		cfg.AllowAllOrigins = true
		cfg.AllowCredentials = false
		return cfg
	}

	cfg.AllowOrigins = origins
	cfg.AllowCredentials = true
	return cfg
}

func splitAndTrim(s string) []string {
	if s == "" {
		return nil
	}
	parts := strings.Split(s, ",")
	out := make([]string, 0, len(parts))
	for _, p := range parts {
		if t := strings.TrimSpace(p); t != "" {
			out = append(out, t)
		}
	}
	return out
}
