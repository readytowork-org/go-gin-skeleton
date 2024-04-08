package internal

import (
	"boilerplate-api/internal/config"
	"boilerplate-api/internal/middlewares"
	"boilerplate-api/internal/request_validator"
	"boilerplate-api/internal/router"
	"go.uber.org/fx"
)

var Module = fx.Options(
	config.Module,
	middlewares.Module,
	fx.Provide(router.NewRouter),
	fx.Provide(request_validator.NewValidator),
)
