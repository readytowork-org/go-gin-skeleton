package controllers

import (
	"boilerplate-api/api/services"
	"boilerplate-api/errors"
	"boilerplate-api/infrastructure"
	"boilerplate-api/responses"
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
	responses.InterfaceJson(c, http.StatusOK, billingData)
}

// GCPBilling -> Get Cost
func (cc GCPBillingController) GetBudgetInfo(c *gin.Context) {
	billingData, err := cc.service.GetExistingBudgetList(c)
	if err != nil {
		cc.logger.Zap.Error("Error fetching Billing Info records", err.Error())
		err := errors.InternalError.Wrap(err, "Failed To Find Billing info for GCP project")
		responses.HandleError(c, err)
		return
	}
	responses.InterfaceJson(c, http.StatusOK, billingData)
}

// GCPBilling -> Get Cost
func (cc GCPBillingController) CreateUpdateBudget(c *gin.Context) {
	billingData, err := cc.service.CreateOrUpdateBudget(c)
	if err != nil {
		cc.logger.Zap.Error("Error fetching Billing Info records", err.Error())
		err := errors.InternalError.Wrap(err, "Failed To Find Billing info for GCP project")
		responses.HandleError(c, err)
		return
	}
	responses.InterfaceJson(c, http.StatusOK, billingData)
}
