package routes

import "go.uber.org/fx"

// Module exports dependency to container
var Module = fx.Options(
	fx.Provide(NewRoutes),
	fx.Provide(NewUtilityRoutes),
	fx.Provide(NewUserRoutes),
	fx.Provide(NewProductRoutes),
	fx.Provide(NewThirdPartyRoutes),
)

// Routes contains multiple routes
type Routes []Route

// Route interface
type Route interface {
	Setup()
}

// NewRoutes sets up routes
func NewRoutes(
	utilityRoutes UtilityRoutes,
	userRoutes UserRoutes,
	productRoutes ProductRoutes,
	thirdPartyRoutes ThirdPartyRoutes,
) Routes {
	return Routes{
		utilityRoutes,
		userRoutes,
		productRoutes,
		thirdPartyRoutes,
	}
}

// Setup all the route
func (r Routes) Setup() {
	for _, route := range r {
		route.Setup()
	}
}
