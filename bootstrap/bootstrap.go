package bootstrap

import (
	"context"

	"boilerplate-api/api/controllers"
	"boilerplate-api/api/repository"
	"boilerplate-api/api/routes"
	"boilerplate-api/api/services"
	"boilerplate-api/cli"
	"boilerplate-api/database/seeds"
	"boilerplate-api/docs"
	"boilerplate-api/external_services"
	"boilerplate-api/internal"
	"boilerplate-api/internal/config"
	"boilerplate-api/internal/router"
	"boilerplate-api/internal/utils"
	"go.uber.org/fx"
)

// Module exported for initializing application
var Module = fx.Options(
	controllers.Module,
	routes.Module,
	services.Module,
	repository.Module,
	internal.Module,
	external_services.Module,
	cli.Module,
	seeds.Module,
	fx.Supply(config.EnvPath(".env")),
	fx.Invoke(bootstrap),
)

func bootstrap(
	lifecycle fx.Lifecycle,
	handler router.Router,
	routes routes.Routes,
	env config.Env,
	logger config.Logger,
	database config.Database,
	cliApp cli.Application,
	migrations config.Migrations,
	seeds seeds.Seeds,
) {

	appStop := func(context.Context) error {
		logger.Info("Stopping Application")
		conn, _ := database.DB.DB()
		_ = conn.Close()
		return nil
	}

	if utils.IsCli() {
		lifecycle.Append(fx.Hook{
			OnStart: func(context.Context) error {
				logger.Info("Starting Golf Simulation cli Application")
				logger.Info("------ 🤖 Golf Simulation 🤖 (CLI) ------")
				go cliApp.Start()
				return nil
			},
			OnStop: appStop,
		})

		return
	}

	lifecycle.Append(fx.Hook{
		OnStart: func(context.Context) error {
			logger.Info("Starting Application")
			logger.Info("------------------------")
			logger.Info("------ Golf Simulation 📺 ------")
			logger.Info("------------------------")

			go func() {
				if env.Environment != "production" && env.HOST != "" {
					logger.Info("Setting Swagger Host...")
					docs.SwaggerInfo.Host = env.HOST
				}

				if env.Environment == "development" || env.Environment == "production" {
					logger.Info("Migrating DB schema...")
					migrations.MigrateUp()
				}
				routes.Setup()
				logger.Info("🌱 seeding data...")
				seeds.Run()
				if env.ServerPort == "" {
					_ = handler.Run()
				} else {
					_ = handler.Run(":" + env.ServerPort)
				}
			}()
			return nil
		},
		OnStop: appStop,
	})
}
