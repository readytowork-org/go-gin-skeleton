package external_services

import (
	"boilerplate-api/external_services/aws"
	"boilerplate-api/external_services/firebase"
	"boilerplate-api/external_services/gcp"
	"go.uber.org/fx"
)

var Module = fx.Options(
	firebase.Module,
	aws.Module,
	gcp.Module,
	fx.Provide(NewStripeService),
	fx.Provide(NewGmailService),
)
