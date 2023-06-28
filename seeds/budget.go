package seeds

import (
	"boilerplate-api/api/services"
	"boilerplate-api/infrastructure"
	"context"
	"time"
)

// ProjectBudgetSeed  Budget setup seed
type ProjectBudgetSeed struct {
	logger        infrastructure.Logger
	budgetService services.GCPBillingService
	env           infrastructure.Env
}

// NewAdminSeed creates budget if set on environment variable
func NewProjectBudgetSeed(
	logger infrastructure.Logger,
	budgetService services.GCPBillingService,
	env infrastructure.Env,
) ProjectBudgetSeed {
	return ProjectBudgetSeed{
		logger:        logger,
		budgetService: budgetService,
		env:           env,
	}
}

func (c ProjectBudgetSeed) getContext() context.Context {

	// Create a context.Background()
	ctx := context.Background()

	// Create a context.WithCancel() to create a cancellable context
	context.WithCancel(ctx)

	// Create a context.WithTimeout() to create a context with a timeout
	timeout := 5 * time.Second
	context.WithTimeout(ctx, timeout)

	// Create a context.WithDeadline() to create a context with a deadline
	deadline := time.Now().Add(10 * time.Second)
	context.WithDeadline(ctx, deadline)

	return ctx
}

// Run the seed data
func (c ProjectBudgetSeed) Run() {
	c.logger.Zap.Info("🌱 seeding  budget alert related setup...")

	if c.env.SetBudget == 1 {
		ctx := c.getContext()
		_, err := c.budgetService.CreateOrUpdateBudget(ctx)

		if err != nil {
			c.logger.Zap.Info("There is an error setting up budget alert ")
		} else {
			c.logger.Zap.Info("budget alert setup successfully")

		}

	}

}
