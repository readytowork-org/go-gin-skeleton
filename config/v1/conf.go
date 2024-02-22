package v1

import (
	user "boilerplate-api/apps/v1/user/init"

	"go.uber.org/fx"
)

var InstalledApps = fx.Options(
	user.Module,
	RouterModule,
)
