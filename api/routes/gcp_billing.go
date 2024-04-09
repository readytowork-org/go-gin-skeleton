package routes

import (
	"boilerplate-api/api/controllers"
	"boilerplate-api/internal/config"
	"boilerplate-api/internal/router"
)

// GCPBillingRoutes -> gcp billing routes struct
type GCPBillingRoutes struct {
	router               router.Router
	Logger               config.Logger
	GCPBillingController controllers.GCPBillingController
}

// NewGCPBillingRoutes -> returns new gcp billing route
func NewGCPBillingRoutes(
	logger config.Logger,
	router router.Router,
	GCPBillingController controllers.GCPBillingController,
) GCPBillingRoutes {
	return GCPBillingRoutes{
		Logger:               logger,
		router:               router,
		GCPBillingController: GCPBillingController,
	}
}

// Setup -> sets up route for util entities
func (u GCPBillingRoutes) Setup() {
	gcpBilling := u.router.V1.Group("/gcp-billing")
	{
		gcpBilling.GET("", u.GCPBillingController.GetCost)
		gcpBilling.GET("budget", u.GCPBillingController.GetBudgetInfo)
		gcpBilling.POST("budget", u.GCPBillingController.CreateUpdateBudget)
	}
}
