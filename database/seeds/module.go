package seeds

import (
	"boilerplate-api/internal/config"
	"go.uber.org/fx"
)

// Seed db seed
type Seed interface {
	Run()
}

// Module exports seed module
var Module = fx.Module("seeds",
	fx.Options(
		fx.Provide(
			NewAdminSeed,
			NewProjectBudgetSeed,
		),
		fx.Invoke(SetupSeeds),
	),
)

// SetupSeeds creates new seeds
func SetupSeeds(
	logger config.Logger,
	adminSeed AdminSeed,
	budgetSeed ProjectBudgetSeed,
) {
	logger.Info("🌱 seeding data...")
	adminSeed.Run()
	budgetSeed.Run()
}
