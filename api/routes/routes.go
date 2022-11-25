package routes

import "go.uber.org/fx"

// Module exports dependency to container
var Module = fx.Options(
	fx.Provide(NewRoutes),
	fx.Provide(NewUtilityRoutes),
	fx.Provide(NewUserRoutes),
	fx.Provide(NewCategoryRoutes),
	fx.Provide(NewBlogRoutes),
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
	categoryRoutes CategoryRoutes,
	blogRoutes BlogRoutes,
) Routes {
	return Routes{
		utilityRoutes,
		userRoutes,
		categoryRoutes,
		blogRoutes,
	}
}

// Setup all the route
func (r Routes) Setup() {
	for _, route := range r {
		route.Setup()
	}
}
