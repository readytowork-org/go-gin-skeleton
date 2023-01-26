package packages

import (
	"boilerplate-api/app/packages/user"
	"boilerplate-api/app/packages/utility"

	"go.uber.org/fx"
)

// Module exported for initializing application
var Module = fx.Options(
	fx.Provide(user.UserController),
	fx.Provide(user.UserRepository),
	fx.Provide(user.UserService),
	fx.Provide(utility.UtilityController),
)
