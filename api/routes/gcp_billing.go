package routes

import (
	"boilerplate-api/api/controllers"
	"boilerplate-api/infrastructure"
)

// GCPBillingRoutes -> gcp billing routes struct
type GCPBillingRoutes struct {
	router               infrastructure.Router
	Logger               infrastructure.Logger
	GCPBillingController controllers.GCPBillingController
}

// NewGCPBillingRoute -> returns new gcp billing route
func NewGCPBillingRoutes(
	logger infrastructure.Logger,
	router infrastructure.Router,
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
	gcpBilling := u.router.Gin.Group("/gcp-billing")
	{
		gcpBilling.GET("", u.GCPBillingController.GetCost)
	}
}
