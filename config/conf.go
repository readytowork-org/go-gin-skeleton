package config

import (
	"boilerplate-api/infrastructure"

	"go.uber.org/fx"
)

var InstalledApps = fx.Options(
	infrastructure.Module,
	RouterModule,
)
