package cli

import (
	"go.uber.org/fx"
)

// Module exports modules
var Module = fx.Options(
	fx.Provide(NewCreateSeedData),
	fx.Provide(NewMigrate),
	fx.Provide(NewApplication),
)
