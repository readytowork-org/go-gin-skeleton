package external_services

import (
	"boilerplate-api/external_services/aws"
	"boilerplate-api/external_services/firebase"
	"go.uber.org/fx"
)

var Module = fx.Options(
	firebase.Module,
	aws.Module,
	fx.Provide(NewStorageBucketService),
	fx.Provide(NewStripeService),
	fx.Provide(NewGmailService),
)
