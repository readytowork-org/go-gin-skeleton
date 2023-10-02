package infrastructure

import (
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

// NewOAuthClient
func NewOAuthClient(env Env) *oauth2.Config {

	OAuthClient := &oauth2.Config{
		ClientID:     env.OAuthClientId,
		ClientSecret: env.OAuthClientSecret,
		Endpoint:     google.Endpoint,
		RedirectURL:  "http://localhost",
		Scopes:       []string{"https://www.googleapis.com/auth/userinfo.email"},
	}

	return OAuthClient
}
