package infrastructure

import "go.uber.org/fx"

var Module = fx.Options(
	fx.Provide(NewLogger),
	fx.Provide(NewRouter),
	fx.Provide(NewEnv),
	fx.Provide(NewDatabase),
	fx.Provide(NewMigrations),
)
