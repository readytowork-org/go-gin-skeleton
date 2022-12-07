package services

import "go.uber.org/fx"

// Module exports services present
var Module = fx.Options(
	fx.Provide(NewFirebaseService),
	fx.Provide(NewStorageBucketService),
	fx.Provide(NewUserService),
	fx.Provide(NewGmailService),
	fx.Provide(NewS3BucketService),
	fx.Provide(NewProductService),
	fx.Provide(NewThirdPartyService),
)
