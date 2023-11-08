package middlewares

import "go.uber.org/fx"

var Module = fx.Provide(
	AuthMiddlewareConstrctor,
)
