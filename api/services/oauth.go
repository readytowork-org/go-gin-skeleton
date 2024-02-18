package services

import (
	"boilerplate-api/api/repository"
	"boilerplate-api/errors"
	"boilerplate-api/infrastructure"
	"boilerplate-api/models"
	"context"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"golang.org/x/oauth2"
)

type OAuthService struct {
	oAuthService   *oauth2.Config
	logger         infrastructure.Logger
	userRepository repository.UserRepository
}

func NewOAuthService(
	logger infrastructure.Logger,
	oAuthService *oauth2.Config,
	env infrastructure.Env,
	userRepository repository.UserRepository,
) OAuthService {
	return OAuthService{
		oAuthService:   oAuthService,
		userRepository: userRepository,
	}
}

func (c OAuthService) GetURL(randomString string) string {
	// Retrieves data even if user is in offline state
	url := c.oAuthService.AuthCodeURL(randomString, oauth2.AccessTypeOffline)
	return url
}

func (c OAuthService) GetToken(code string) (*oauth2.Token, error) {
	token, err := c.oAuthService.Exchange(context.Background(), code)
	if err != nil {
		return nil, err
	}
	return token, err
}

func (c OAuthService) GetHeaderTokenAndAuthorize(ctx *gin.Context) (*models.User, error) {
	// Get the token from the request header and check if user is in our DB or not
	header := ctx.GetHeader("Authorization")
	if header == "" {
		err := errors.BadRequest.New("Authorization token is required in header")
		err = errors.SetCustomMessage(err, "Authorization token is required in header")
		c.logger.Zap.Error("[GetHeader]: ", err.Error())
		return nil, err
	}

	if !strings.Contains(header, "Bearer") {
		err := errors.BadRequest.New("Token type is required")
		c.logger.Zap.Error("Missing token type: ", err.Error())
		return nil, err
	}
	tokenString := strings.TrimSpace(strings.Replace(header, "Bearer", "", 1))

	// Check with database
	getUser, err := c.userRepository.GetOneUserWithToken(tokenString)

	if err != nil || getUser == nil || getUser.TokenExpiryTime.Before(time.Now()) {
		return nil, err
	}
	return getUser, nil
}
