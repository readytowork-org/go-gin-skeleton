package controllers

import (
	"boilerplate-api/api/responses"
	"boilerplate-api/api/services"
	"boilerplate-api/errors"
	"boilerplate-api/infrastructure"
	"net/http"

	"github.com/gin-gonic/gin"
)

// GCPBilling -> struct
type GCPBillingController struct {
	logger  infrastructure.Logger
	env     infrastructure.Env
	service services.GCPBillingService
}

// NewGCPBillingController -> constructor
func NewGCPBillingController(
	logger infrastructure.Logger,
	env infrastructure.Env,
	service services.GCPBillingService,
) GCPBillingController {
	return GCPBillingController{
		logger:  logger,
		env:     env,
		service: service,
	}
}

// GCPBilling -> Get Cost
func (cc GCPBillingController) GetCost(c *gin.Context) {
	billingData, err := cc.service.GetBillingInfo(c)
	if err != nil {
		cc.logger.Zap.Error("Error fetching Billing Info records", err.Error())
		err := errors.InternalError.Wrap(err, "Failed To Find Billing info for GCP project")
		responses.HandleError(c, err)
		return
	}
	responses.SuccessJSON(c, http.StatusOK, billingData)
}
