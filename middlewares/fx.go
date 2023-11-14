package middlewares

import "go.uber.org/fx"

var Module = fx.Options(
	fx.Provide(NewDBTransaction),
	fx.Provide(NewJWTAuthMiddleWare),
)
