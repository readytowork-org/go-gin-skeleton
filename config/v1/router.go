package v1

import (
	user "boilerplate-api/apps/v1/user/routes"

	"go.uber.org/fx"
)

// / Module exports dependency to container
var RouterModule = fx.Options(
	fx.Provide(RoutersConstructor),
)

// Routes contains multiple routes
type Routes []Route

// Route interface
type Route interface {
	Setup()
}

// NewRoutes sets up routes
func RoutersConstructor(
	UserRoutes user.UserRoute,
) Routes {
	return Routes{
		UserRoutes,
	}
}

// Setup all the route
func (r Routes) Setup() {
	for _, route := range r {
		route.Setup()
	}
}
