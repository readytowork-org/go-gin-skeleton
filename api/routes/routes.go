package routes

import "go.uber.org/fx"

// Module exports dependency to container
var Module = fx.Options(
	fx.Provide(NewRoutes),
	fx.Provide(NewUtilityRoutes),
	fx.Provide(NewDocsRoutes),
	fx.Provide(NewJwtAuthRoutes),
	fx.Provide(NewUserRoutes),
	fx.Provide(NewGCPBillingRoutes),
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
	gCPBillingRoutes GCPBillingRoutes,
	jwtAuthRoutes JwtAuthRoutes,
	docsRoutes DocsRoutes,
) Routes {
	return Routes{
		utilityRoutes,
		jwtAuthRoutes,
		userRoutes,
		gCPBillingRoutes,
		docsRoutes,
	}
}

// Setup all the route
func (r Routes) Setup() {
	for _, route := range r {
		route.Setup()
	}
}
