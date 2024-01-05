package config

import (
	"boilerplate-api/common/utils"
	"boilerplate-api/infrastructure"
	"boilerplate-api/middlewares"
	"context"

	"go.uber.org/fx"
)

var Module = fx.Options(
	InstalledApps,
	fx.Invoke(bootstrap),
)

func bootstrap(
	lifecycle fx.Lifecycle,
	handler infrastructure.Router,
	env infrastructure.Env,
	logger infrastructure.Logger,
	database infrastructure.Database,
	migrations infrastructure.Migrations,
	middlewares middlewares.Middlewares,
	routes Routes,
) {

	appStop := func(context.Context) error {
		logger.Zap.Info("Stopping Application")
		conn, _ := database.DB.DB()
		_ = conn.Close()
		return nil
	}

	if utils.IsCli() {
		lifecycle.Append(fx.Hook{
			OnStart: func(context.Context) error {
				logger.Zap.Info("Starting boilerplate cli Application")
				logger.Zap.Info("------ 🤖 Boilerplate 🤖 (CLI) ------")
				return nil
			},
			OnStop: appStop,
		})

		return
	}

	lifecycle.Append(fx.Hook{
		OnStart: func(context.Context) error {
			logger.Zap.Info("Starting Application")
			logger.Zap.Info("------------------------")
			logger.Zap.Info("------ Boilerplate 📺 ------")
			logger.Zap.Info("------------------------")

			go func() {
				if env.Environment == "production" {
					logger.Zap.Info("Migrating DB schema...")
					migrations.Migrate()
				}
				middlewares.Setup()
				routes.Setup()
				if env.ServerPort == "" {
					_ = handler.Gin.Run()
				} else {
					_ = handler.Gin.Run(":" + env.ServerPort)
				}
			}()
			return nil
		},
		OnStop: appStop,
	})
}
