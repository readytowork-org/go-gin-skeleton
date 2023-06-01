package infrastructure

import (
	"context"

	budgets "cloud.google.com/go/billing/budgets/apiv1"
	"google.golang.org/api/cloudbilling/v1"
	"google.golang.org/api/option"
)

type GCPBilling struct {
	BillingClient *cloudbilling.APIService
	BudgetClient  *budgets.BudgetClient
}

// NewGCPBilling creates a new gcp billing api client
func NewGCPBilling(logger Logger, env Env) GCPBilling {
	ctx := context.Background()
	credentials := option.WithCredentialsFile("serviceAccountKey.json")

	billingClient, err := cloudbilling.NewService(ctx, credentials)

	if err != nil {
		logger.Zap.Error("Failed to create cloud billing api client: %v \n", err)
	}

	budgetClient, err := budgets.NewBudgetClient(ctx, option.WithCredentialsFile("serviceAccountKey.json"))

	if err != nil {
		logger.Zap.Error("Failed to create cloud budget api client: %v \n", err)
	}

	return GCPBilling{
		BillingClient: billingClient,
		BudgetClient:  budgetClient,
	}
}
