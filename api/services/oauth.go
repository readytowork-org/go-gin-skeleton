package services

import (
	"boilerplate-api/infrastructure"
	"context"

	"golang.org/x/oauth2"
)

type OAuthService struct {
	oAuthService *oauth2.Config
}

func NewOAuthService(oAuthService *oauth2.Config, env infrastructure.Env) OAuthService {
	return OAuthService{
		oAuthService: oAuthService,
	}
}

func (c OAuthService) GetURL(randomString string) string {
	url := c.oAuthService.AuthCodeURL(randomString)
	return url
}

func (c OAuthService) GetToken(code string) (*oauth2.Token, error) {
	token, err := c.oAuthService.Exchange(context.Background(), code)
	if err != nil {
		return nil, err
	}
	return token, err
}
