package seeds

import "go.uber.org/fx"

// Module exports seed module
var Module = fx.Options(
	fx.Provide(NewAdminSeed),
	fx.Provide(NewSeeds),
)

// Seed db seed
type Seed interface {
	Run()
}

// Seeds listing of seeds
type Seeds []Seed

// Run run the seed data
func (s Seeds) Run() {
	for _, seed := range s {
		seed.Run()
	}
}

// NewSeeds creates new seeds
func NewSeeds(
	adminSeed AdminSeed,
) Seeds {
	return Seeds{
		adminSeed,
	}
}
