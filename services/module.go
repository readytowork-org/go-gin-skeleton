package services

import (
	"go.uber.org/fx"
)

var Module = fx.Options(
	//aws.Module,
	// StripeService provider
	//fx.Provide(
	//	func(
	//		env config.Env,
	//		logger config.Logger,
	//	) StripeService {
	//		return NewStripeService(
	//			StripeConfig{
	//				stripeSecretKey: env.StripeSecretKey,
	//				stripeProductID: env.StripeProductID,
	//				logger:          logger.SugaredLogger,
	//			},
	//		)
	//	},
	//),
	// GmailService provider
	//fx.Provide(
	//	func(
	//		env config.Env,
	//		logger config.Logger,
	//	) *GmailService {
	//		return NewGmailService(
	//			GmailConfig{
	//				ClientID:     env.MailClientID,
	//				ClientSecret: env.MailClientSecret,
	//				AccessToken:  env.MailAccesstoken,
	//				RefreshToken: env.MailRefreshToken,
	//				HostURL:      env.HOST,
	//			},
	//			logger,
	//		)
	//	},
	//),
	// TwilioService provider
	//fx.Provide(
	//	func(
	//		env config.Env,
	//		logger config.Logger,
	//	) TwilioService {
	//		return NewTwilioService(
	//			TwilioService{
	//				baseURL:   env.TwilioBaseURL,
	//				smsFrom:   env.TwilioSMSFrom,
	//				sID:       env.TwilioSID,
	//				authToken: env.TwilioAuthToken,
	//				logger:    logger.SugaredLogger,
	//			},
	//		)
	//	},
	//),
)
