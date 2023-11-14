package config

import (
	auth "boilerplate-api/apps/auth/auth_router"
	user "boilerplate-api/apps/user/user_router"

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
	user user.Route,
	auth auth.Route,
) Routes {
	return Routes{
		user,
		auth,
	}
}

// Setup all the route
func (r Routes) Setup() {
	for _, route := range r {
		route.Setup()
	}
}
