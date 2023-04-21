package services

import (
	"boilerplate-api/infrastructure"
	"context"

	"google.golang.org/api/cloudbilling/v1"
)

// GCPBillingService -> handles the gcp billing related functions
type GCPBillingService struct {
	logger     infrastructure.Logger
	gcpBilling infrastructure.GCPBilling
	env        infrastructure.Env
}

// NewGCPBillingService -> initilization for the GCPBilling struct
func NewGCPBillingService(
	logger infrastructure.Logger,
	gcpBilling infrastructure.GCPBilling,
	env infrastructure.Env,
) GCPBillingService {
	return GCPBillingService{
		logger:     logger,
		gcpBilling: gcpBilling,
		env:        env,
	}
}

// Get Billing info for certain date
func (s GCPBillingService) GetBillingInfo(
	ctx context.Context,
) (*cloudbilling.ProjectBillingInfo, error) {
	billingInfo, err := s.gcpBilling.Client.Projects.GetBillingInfo("projects/easyy-dev").Do()
	return billingInfo, err
}
