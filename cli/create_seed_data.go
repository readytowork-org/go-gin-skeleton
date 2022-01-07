package cli

import (
	"boilerplate-api/infrastructure"
	"boilerplate-api/seeds"
)

// CreateSeedData command
type CreateSeedData struct {
	logger infrastructure.Logger
	seeds  seeds.Seeds
}

// NewCreateSeedData creates instance of admin user
func NewCreateSeedData(
	logger infrastructure.Logger,
	seeds seeds.Seeds,
) CreateSeedData {
	return CreateSeedData{
		logger: logger,
		seeds:  seeds,
	}
}

// Run runs command
func (c CreateSeedData) Run() {
	c.logger.Zap.Info("🌱 Creating seed data...")
	c.seeds.Run()
}

// Name return name of command
func (c CreateSeedData) Name() string {
	return "CREATE_SEED_DATA"
}
