package router

import (
	"context"
	"fmt"
	"net/http"
	"strings"
	"time"

	"boilerplate-api/lib/config"
	"boilerplate-api/lib/middlewares"

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
func NewRouter(
	env config.Env,
	logger config.Logger,
	database config.Database,
	idempotency middlewares.IdempotencyMiddleware,
) Router {
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

	// RequestID must run before anything that might log or emit errors so the
	// request id is available in error envelopes and Sentry events.
	httpRouter.Use(middlewares.RequestID())
	httpRouter.Use(cors.New(buildCorsConfig(env)))
	httpRouter.Use(middlewares.ErrorHandler(logger))
	// Idempotency is mounted globally and self-gates: it only acts on POST
	// requests carrying an Idempotency-Key header, so non-POST and untagged
	// requests pay only a tiny method/header check.
	httpRouter.Use(idempotency.Handle())

	httpRouter.Use(
		sentrygin.New(
			sentrygin.Options{
				Repanic: true,
			},
		),
	)

	httpRouter.GET(
		"/livez", func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{"data": "alive"})
		},
	)

	healthCheck := func(c *gin.Context) {
		if database.DB == nil {
			c.JSON(http.StatusServiceUnavailable, gin.H{"status": "unavailable", "error": "db not initialized"})
			return
		}
		sqlDB, err := database.DB.DB()
		if err != nil {
			c.JSON(http.StatusServiceUnavailable, gin.H{"status": "unavailable", "error": err.Error()})
			return
		}
		ctx, cancel := context.WithTimeout(c.Request.Context(), 2*time.Second)
		defer cancel()
		if err := sqlDB.PingContext(ctx); err != nil {
			c.JSON(http.StatusServiceUnavailable, gin.H{"status": "unavailable", "error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{"data": " 📺 API Up and Running", "db": "ok"})
	}
	httpRouter.GET("/health-check", healthCheck)
	httpRouter.GET("/readyz", healthCheck)

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
