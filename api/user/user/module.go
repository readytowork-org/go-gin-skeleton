package user

import (
	"go.uber.org/fx"
)

var Module = fx.Module("user",
	fx.Options(
		fx.Provide(
			NewRepository,
			NewService,
			NewController,
		),
		fx.Invoke(SetupRoutes),
	))
