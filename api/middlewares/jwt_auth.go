package middlewares

import (
	"boilerplate-api/api/responses"
	"boilerplate-api/api/services"
	"boilerplate-api/errors"
	"boilerplate-api/infrastructure"

	"github.com/getsentry/sentry-go"
	"github.com/gin-gonic/gin"
)

type JWTAuthMiddleWare struct {
	jwtService services.JWTAuthService
	logger     infrastructure.Logger
	env        infrastructure.Env
	db         infrastructure.Database
}

func NewJWTAuthMiddleWare(
	jwtService services.JWTAuthService,
	logger infrastructure.Logger,
	env infrastructure.Env,
	db infrastructure.Database,

) JWTAuthMiddleWare {
	return JWTAuthMiddleWare{
		jwtService: jwtService,
		logger:     logger,
		env:        env,
		db:         db,
	}
}

// Authenticate user with jwt using this middleware
func (m JWTAuthMiddleWare) Handle() gin.HandlerFunc {
	return func(c *gin.Context) {
		//Getting token from header
		tokenString, err := m.jwtService.GetTokenFromHeader(c)
		if err != nil {
			m.logger.Zap.Error("Error getting token from header: ", err.Error())
			err = errors.Unauthorized.Wrap(err, "Something went wrong")
			responses.HandleError(c, err)
			c.Abort()
			return
		}
		// Parsing token
		parsedToken, parseErr := m.jwtService.ParseToken(tokenString, m.env.JWT_ACCESS_SECRET)
		if parseErr != nil {
			m.logger.Zap.Error("Error parsing token: ", parseErr.Error())
			err = errors.Unauthorized.Wrap(parseErr, "Something went wrong")
			responses.HandleError(c, err)
			c.Abort()
			return
		}
		//verify token
		claims, verifyErr := m.jwtService.VerifyToken(parsedToken)
		if verifyErr != nil {
			m.logger.Zap.Error("Error veriefying token: ", verifyErr.Error())
			err = errors.Unauthorized.Wrap(verifyErr, "Something went wrong")
			responses.HandleError(c, err)
			c.Abort()
			return
		}
		// ser user to the scope
		sentry.ConfigureScope(func(scope *sentry.Scope) {
			scope.SetUser(sentry.User{ID: claims.Id})
		})
		// Can set anything in the request context and passes the request to the next handler.
		c.Set("user_id", claims.Id)
		c.Next()

	}
}
