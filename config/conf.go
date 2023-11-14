package config

import (
	"boilerplate-api/apps/user"
	"boilerplate-api/apps/user/user_router"
	"boilerplate-api/infrastructure"
	"boilerplate-api/middlewares"

	"go.uber.org/fx"
)

var InstalledRoutes = fx.Options(
	fx.Provide(user_router.RouteConstructor),
)

var InstalledApps = fx.Options(
	user.Module,
	infrastructure.Module,
	middlewares.Module,
	RouterModule,
	InstalledRoutes,
)
