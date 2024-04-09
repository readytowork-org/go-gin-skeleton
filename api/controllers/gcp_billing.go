package controllers

import (
	"net/http"

	"boilerplate-api/external_services/gcp"
	"boilerplate-api/internal/api_errors"
	"boilerplate-api/internal/config"
	"boilerplate-api/internal/json_response"
	"cloud.google.com/go/billing/budgets/apiv1/budgetspb"
	"google.golang.org/api/cloudbilling/v1"

	"github.com/gin-gonic/gin"
)

// GCPBillingController -> struct
type GCPBillingController struct {
	logger  config.Logger
	env     config.Env
	service gcp.BillingService
}

// NewGCPBillingController -> constructor
func NewGCPBillingController(
	logger config.Logger,
	env config.Env,
	service gcp.BillingService,
) GCPBillingController {
	return GCPBillingController{
		logger:  logger,
		env:     env,
		service: service,
	}
}

// GetCost -> Get Cost
func (cc GCPBillingController) GetCost(c *gin.Context) {
	billingData, err := cc.service.GetBillingInfo()
	if err != nil {
		cc.logger.Error("Error fetching Billing Info records", err.Error())
		err := api_errors.InternalError.Wrap(err, "Failed To Find Billing info for GCP project")
		status, errM := api_errors.HandleError(err)
		c.JSON(status, json_response.Error{Error: errM})
		return
	}
	c.JSON(http.StatusOK, json_response.Data[*cloudbilling.ProjectBillingInfo]{
		Data: billingData,
	})
}

// GetBudgetInfo -> Get Cost
func (cc GCPBillingController) GetBudgetInfo(c *gin.Context) {
	billingData, err := cc.service.GetExistingBudgetList(c)
	if err != nil {
		cc.logger.Error("Error fetching Billing Info records", err.Error())
		err := api_errors.InternalError.Wrap(err, "Failed To Find Billing info for GCP project")
		status, errM := api_errors.HandleError(err)
		c.JSON(status, json_response.Error{Error: errM})
		return
	}
	c.JSON(http.StatusOK, json_response.Data[*budgetspb.Budget]{
		Data: billingData,
	})
}

// CreateUpdateBudget -> Get Cost
func (cc GCPBillingController) CreateUpdateBudget(c *gin.Context) {
	billingData, err := cc.service.CreateOrUpdateBudget(c)
	if err != nil {
		cc.logger.Error("Error fetching Billing Info records", err.Error())
		err := api_errors.InternalError.Wrap(err, "Failed To Find Billing info for GCP project")
		status, errM := api_errors.HandleError(err)
		c.JSON(status, json_response.Error{Error: errM})
		return
	}
	c.JSON(http.StatusOK, json_response.Data[*budgetspb.Budget]{
		Data: billingData,
	})
}
