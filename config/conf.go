package config

import (
	"boilerplate-api/apps/auth"
	"boilerplate-api/apps/auth/auth_router"
	"boilerplate-api/apps/user"
	"boilerplate-api/apps/user/user_router"
	"boilerplate-api/infrastructure"
	"boilerplate-api/middlewares"

	"go.uber.org/fx"
)

var InstalledRoutes = fx.Options(
	fx.Provide(user_router.RouteConstructor),
	fx.Provide(auth_router.RouteConstructor),
)

var InstalledApps = fx.Options(
	user.Module,
	auth.Module,
	infrastructure.Module,
	middlewares.Module,
	RouterModule,
	InstalledRoutes,
)
