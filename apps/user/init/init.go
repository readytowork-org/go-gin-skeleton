package init

import (
	"boilerplate-api/apps/user"
	"boilerplate-api/apps/user/controllers"
	"boilerplate-api/apps/user/repository"
	"boilerplate-api/apps/user/routes"
	"boilerplate-api/apps/user/services"

	"go.uber.org/fx"
)

// Module exported for initializing application
var Module = fx.Options(
	fx.Provide(controllers.ControllerConstuctor),
	fx.Provide(services.UserServiceConstuctor),
	fx.Provide(repository.UserRepositoryConstuctor),
	fx.Provide(routes.UserRouteConstructor),
	fx.Provide(user.NewUserValidator),
)
