package init

import (
	"boilerplate-api/apps/v1/user"
	"boilerplate-api/apps/v1/user/controllers"
	"boilerplate-api/apps/v1/user/repository"
	"boilerplate-api/apps/v1/user/routes"
	"boilerplate-api/apps/v1/user/services"

	"go.uber.org/fx"
)

// Module exported for initializing application
var Module = fx.Options(
	fx.Provide(controllers.UserControllerConstuctor),
	fx.Provide(services.UserServiceConstuctor),
	fx.Provide(repository.UserRepositoryConstuctor),
	fx.Provide(routes.UserRouteConstructor),
	fx.Provide(user.NewUserValidator),
)
