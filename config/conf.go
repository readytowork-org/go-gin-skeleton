package config

import (
	user "boilerplate-api/apps/user/init"
	"boilerplate-api/infrastructure"
	"boilerplate-api/middlewares"

	"go.uber.org/fx"
)

var InstalledApps = fx.Options(
	user.Module,
	infrastructure.Module,
	middlewares.Module,
	RouterModule,
)
