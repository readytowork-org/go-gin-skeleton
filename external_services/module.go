package external_services

import (
	"boilerplate-api/external_services/firebase"
	"go.uber.org/fx"
)

var Module = fx.Options(
	firebase.Module,
	fx.Provide(NewStorageBucketService),
	fx.Provide(NewStripeService),
)
