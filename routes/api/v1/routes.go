package routes

import (
	"boilerplate-api/app/global/infrastructure"

	"github.com/gin-gonic/gin"
	"go.uber.org/fx"
)

// Module exports dependency to container
var Module = fx.Options(
	fx.Provide(NewRoutes),
	fx.Provide(NewUtilityRoutes),
	fx.Provide(NewUserRoutes),
)

// Routes contains multiple routes
type Routes struct {
	routes []Route
	router infrastructure.Router
}

// Route interface
type Route interface {
	Setup(*gin.RouterGroup)
}

// NewRoutes sets up routes
func NewRoutes(
	utilityRoutes UtilityRoutes,
	userRoutes UserRoutes,
	router infrastructure.Router,
) Routes {
	routes := []Route{
		utilityRoutes,
		userRoutes,
	}
	return Routes{
		routes: routes,
		router: router,
	}
}

// Setup all the route
func (r Routes) Setup() {
	v1 := r.router.Gin.Group("/api/v1")
	for _, route := range r.routes {
		route.Setup(v1)
	}
}
