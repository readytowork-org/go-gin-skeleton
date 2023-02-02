package routes

import (
	"boilerplate-api/app/global/infrastructure"
	"boilerplate-api/routes/api/v1"
	"boilerplate-api/routes/web"

	"github.com/gin-gonic/gin"
	"go.uber.org/fx"
)

// Module exports dependency to container
var Module = fx.Options(
	fx.Provide(NewRoutes),
	// api
	fx.Provide(api.NewUserRoutes),
	fx.Provide(api.NewUtilityRoutes),
	// web
	fx.Provide(web.NewHomeRoutes),
)

// Routes contains multiple routes
type Routes struct {
	apiRoutes []Route
	webRoutes []Route
	router    infrastructure.Router
}

// Route interface
type Route interface {
	Setup(*gin.RouterGroup)
}

// NewRoutes sets up routes
func NewRoutes(
	utilityRoutes api.UtilityRoutes,
	userRoutes api.UserRoutes,
	homeRoutes web.HomeRoutes,
	router infrastructure.Router,
) Routes {
	apiRoutes := []Route{
		utilityRoutes,
		userRoutes,
	}
	webRoutes := []Route{
		homeRoutes,
	}
	return Routes{
		apiRoutes: apiRoutes,
		webRoutes: webRoutes,
		router:    router,
	}
}

// Setup all the route
func (r Routes) Setup() {
	api := r.router.Gin.Group("/api/v1")
	for _, route := range r.apiRoutes {
		route.Setup(api)
	}
	web := r.router.Gin.Group("/")
	for _, route := range r.webRoutes {
		route.Setup(web)
	}
}
