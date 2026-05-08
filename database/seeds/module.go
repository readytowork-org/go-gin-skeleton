package seeds

import (
	"go.uber.org/fx"
)

// Seed db seed
type Seed interface {
	Run()
}

// Module exports seed module
var Module = fx.Module(
	"seeds",
	fx.Options(
		//fx.Provide(
		//	fx.Annotate(
		//		NewAdminSeed,
		//		fx.As(new(Seed)),
		//		fx.ResultTags(`group:"seeds"`),
		//	),
		//),
		//fx.Provide(
		//	fx.Annotate(
		//		NewProjectBudgetSeed,
		//		fx.As(new(Seed)),
		//		fx.ResultTags(`group:"seeds"`),
		//	),
		//),
	),
)
