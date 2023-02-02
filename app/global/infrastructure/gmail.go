package infrastructure

import (
	"context"
	"time"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/gmail/v1"
	"google.golang.org/api/option"
)

// NewGmailService -> receive gmail service client
func NewGmailService(logger Logger, env Env) *gmail.Service {
	ctx := context.Background()

	config := oauth2.Config{
		ClientID:     env.MailClientID,
		ClientSecret: env.MailClientSecret,
		Endpoint:     google.Endpoint,
		RedirectURL:  "http://localhost",
		Scopes:       []string{"https://www.googleapis.com/auth/gmail.send"},
	}
	token := oauth2.Token{
		AccessToken:  env.MailAccesstoken,
		RefreshToken: env.MailRefreshToken,
		TokenType:    "Bearer",
		Expiry:       time.Now(),
	}
	var tokenSource = config.TokenSource(ctx, &token)
	srv, err := gmail.NewService(ctx, option.WithTokenSource(tokenSource))
	if err != nil {
		logger.Zap.Fatal("failed to receive gmail client", err.Error())
	}
	return srv

}
