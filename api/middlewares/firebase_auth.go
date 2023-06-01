package middlewares

import (
	"boilerplate-api/api/services"
	"boilerplate-api/constants"
	"boilerplate-api/responses"
	"net/http"
	"strings"

	"firebase.google.com/go/auth"
	"github.com/getsentry/sentry-go"
	"github.com/gin-gonic/gin"
)

// FirebaseAuthMiddleware structure
type FirebaseAuthMiddleware struct {
	service     services.FirebaseService
	userService services.UserService
}

// NewFirebaseAuthMiddleware creates new firebase authentication
func NewFirebaseAuthMiddleware(
	service services.FirebaseService,
	userservice services.UserService,
) FirebaseAuthMiddleware {
	return FirebaseAuthMiddleware{
		service:     service,
		userService: userservice,
	}
}

// Handle handles auth requests
func (f FirebaseAuthMiddleware) Handle() gin.HandlerFunc {
	return func(c *gin.Context) {
		token, err := f.getTokenFromHeader(c)

		if err != nil {
			responses.ErrorJSON(c, http.StatusUnauthorized, err.Error())
			c.Abort()
			return
		}

		sentry.ConfigureScope(func(scope *sentry.Scope) {
			scope.SetUser(sentry.User{ID: token.UID})
		})

		c.Set(constants.Claims, token.Claims)
		c.Set(constants.UID, token.UID)

		c.Next()
	}
}

// HandleAdminOnly handles middleware for admin role only
func (f FirebaseAuthMiddleware) HandleAdminOnly() gin.HandlerFunc {
	return func(c *gin.Context) {
		token, err := f.getTokenFromHeader(c)

		if err != nil {
			responses.ErrorJSON(c, http.StatusUnauthorized, err.Error())
			c.Abort()
			return
		}

		if !f.isAdmin(token.Claims) {
			responses.ErrorJSON(c, http.StatusUnauthorized, "un-authorized request")
			c.Abort()
			return
		}

		sentry.ConfigureScope(func(scope *sentry.Scope) {
			scope.SetUser(sentry.User{ID: token.UID})
		})

		c.Set(constants.Claims, token.Claims)
		c.Set(constants.UID, token.UID)

		c.Next()
	}
}

// getTokenFromHeader gets token from header
func (f FirebaseAuthMiddleware) getTokenFromHeader(c *gin.Context) (*auth.Token, error) {
	header := c.GetHeader("Authorization")
	idToken := strings.TrimSpace(strings.Replace(header, "Bearer", "", 1))

	token, err := f.service.VerifyToken(idToken)
	if err != nil {
		return nil, err
	}

	return token, nil
}

// isAdmin check if claims is admin
func (f FirebaseAuthMiddleware) isAdmin(claims map[string]interface{}) bool {

	role := claims["role"]
	isAdmin := false
	if role != nil {
		isAdmin = role.(string) == "admin"
	}

	return isAdmin

}
