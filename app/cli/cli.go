package cli

import (
	"go.uber.org/fx"
)

// Module exports modules
var Module = fx.Options(
	fx.Provide(NewCreateAdminUser),
	fx.Provide(NewCreateDummyAdminUser),
	fx.Provide(NewCreateSeedData),
	fx.Provide(NewApplication),
)
