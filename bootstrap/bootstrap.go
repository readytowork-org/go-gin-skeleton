package bootstrap

import (
	"context"

	"boilerplate-api/api"
	"boilerplate-api/cli"
	"boilerplate-api/database/seeds"
	"boilerplate-api/lib"
	"boilerplate-api/lib/config"
	"boilerplate-api/lib/router"
	"boilerplate-api/lib/utils"
	"boilerplate-api/services"
	"boilerplate-api/swagger"

	"go.uber.org/fx"
)

// Module exported for initializing application
var Module = fx.Options(
	lib.Module,
	seeds.Module,
	cli.Module,
	services.Module,
	api.Module,
	fx.Supply(config.EnvPath(".env")),
	fx.Invoke(fx.Annotate(bootstrap, fx.ParamTags(`group:"seeds"`))),
)

func bootstrap(
	seeds []seeds.Seed,
	lifecycle fx.Lifecycle,
	router router.Router,
	env config.Env,
	logger config.Logger,
	database config.Database,
	cliApp cli.Application,
	migrations *config.Migrations,
) {

	appStop := func(context.Context) error {
		logger.Info("Stopping Application")
		conn, _ := database.DB.DB()
		_ = conn.Close()
		return nil
	}

	if utils.IsCli() {
		lifecycle.Append(
			fx.Hook{
				OnStart: func(context.Context) error {
					logger.Info("Starting cli Application")
					logger.Info("------- (CLI) ------")
					go cliApp.Start()
					return nil
				},
				OnStop: appStop,
			},
		)

		return
	}

	lifecycle.Append(
		fx.Hook{
			OnStart: func(context.Context) error {
				logger.Info("Starting Application")
				logger.Info("------------------------")
				logger.Info("------ Gin Skeleton 📺 ------")
				logger.Info("------------------------")

				go func() {
					if env.Environment != "production" && env.HOST != "" {
						logger.Info("Setting Swagger Host...")
						swagger.SwaggerInfo.Host = env.HOST
					}

					if database.ConnectionError() != nil {
						logger.Error(database.ConnectionError())
						return
					}

					if env.Environment == "development" || env.Environment == "production" {
						logger.Info("Migrating DB schema...")
						err := migrations.MigrateUp()
						if err != nil {
							logger.Error("Error in migration steps: ", err.Error())
							return
						}
					}

					logger.Info("🌱 seeding data...")
					for _, seed := range seeds {
						seed.Run()
					}

					if env.ServerPort == "" {
						_ = router.Run()
					} else {
						_ = router.Run(":" + env.ServerPort)
					}
				}()
				return nil
			},
			OnStop: appStop,
		},
	)
}
