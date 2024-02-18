package middlewares

import (
	"boilerplate-api/api/services"
	"boilerplate-api/constants"
	"boilerplate-api/errors"
	"boilerplate-api/infrastructure"
	"boilerplate-api/responses"

	"github.com/gin-gonic/gin"
)

type OAuthMiddleWare struct {
	oAuthService services.OAuthService
	logger       infrastructure.Logger
	env          infrastructure.Env
	db           infrastructure.Database
}

func NewOAuthMiddleWare(
	oAuthService services.OAuthService,
	logger infrastructure.Logger,
	env infrastructure.Env,
	db infrastructure.Database,

) OAuthMiddleWare {
	return OAuthMiddleWare{
		oAuthService: oAuthService,
		logger:       logger,
		env:          env,
		db:           db,
	}
}

// Handle user with OAuth using this middleware
// If using postman, use OAuth2.0 token type in Authorization tab
func (m OAuthMiddleWare) Handle() gin.HandlerFunc {
	return func(c *gin.Context) {
		//Getting token from header
		user, err := m.oAuthService.GetHeaderTokenAndAuthorize(c)
		if err != nil {
			m.logger.Zap.Error("Access token header err: ", err.Error())
			err = errors.Unauthorized.Wrap(err, "Error getting token from header")
			responses.HandleError(c, err)
			c.Abort()
			// In Client side/FE, redirect user to OAuth Sign In page again.
			return
		}

		// Can set anything in the request context and passes the request to the next handler.
		c.Set(constants.UserID, user.ID)
		c.Next()
	}
}
