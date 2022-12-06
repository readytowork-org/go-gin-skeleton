package controllers

import "go.uber.org/fx"

// Module exported for initializing application
var Module = fx.Options(
	fx.Provide(NewUserController),
	fx.Provide(NewUtilityController),
	fx.Provide(NewProductController),
	fx.Provide(NewTwitterController),
)
