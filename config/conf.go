package config

import (
	"boilerplate-api/apps/auth"
	authrouter "boilerplate-api/apps/auth/auth_router"
	"boilerplate-api/apps/user"
	userrouter "boilerplate-api/apps/user/user_router"
	"boilerplate-api/infrastructure"
	"boilerplate-api/middlewares"

	"go.uber.org/fx"
)

var InstalledRoutes = fx.Options(
	fx.Provide(userrouter.RouteConstructor),
	fx.Provide(authrouter.RouteConstructor),
)

var InstalledApps = fx.Options(
	user.Module,
	auth.Module,
	infrastructure.Module,
	middlewares.Module,
	RouterModule,
	InstalledRoutes,
)
