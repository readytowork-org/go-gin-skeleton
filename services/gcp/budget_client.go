package gcp

import (
	"boilerplate-api/internal/config"
	"cloud.google.com/go/billing/budgets/apiv1"
	"context"
)

type BudgetClient struct {
	*budgets.BudgetClient
}

func NewGCPBudgetClient(logger config.Logger, clientOption config.GCPClientOption) BudgetClient {
	budgetClient, err := budgets.NewBudgetClient(context.Background(), clientOption)

	if err != nil {
		logger.Panic("Failed to create cloud budget api client: %v \n", err)
	}
	return BudgetClient{
		budgetClient,
	}
}
