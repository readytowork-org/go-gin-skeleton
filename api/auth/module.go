package auth

import (
	"go.uber.org/fx"
)

var Module = fx.Module("auth",
	fx.Options(
		fx.Provide(
			NewJwtAuthController,
			NewRefreshTokenRepository,
		),
		fx.Invoke(SetupRoutes),
	))
