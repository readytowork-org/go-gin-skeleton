package routes

import "go.uber.org/fx"

// Module exports dependency to container
var Module = fx.Options(
	fx.Provide(NewRoutes),
  fx.Provide(NewStudentInfoRoutes),
  fx.Provide(NewPostRoutes),
	fx.Provide(NewUtilityRoutes),
	fx.Provide(NewUserRoutes),
)

// Routes contains multiple routes
type Routes []Route

// Route interface
type Route interface {
	Setup()
}

// NewRoutes sets up routes
func NewRoutes(
	 studentInfoRoutes StudentInfoRoutes,
	 postRoutes PostRoutes,
	utilityRoutes UtilityRoutes,
	userRoutes UserRoutes,
) Routes {
	return Routes{
	 studentInfoRoutes,
	 postRoutes,
		utilityRoutes,
		userRoutes,
	}
}

// Setup all the route
func (r Routes) Setup() {
	for _, route := range r {
		route.Setup()
	}
}
