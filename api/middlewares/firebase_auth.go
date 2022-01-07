package middlewares

import (
	"net/http"
	"strconv"
	"strings"
	"boilerplate-api/api/responses"
	"boilerplate-api/api/services"
	"boilerplate-api/constants"

	"firebase.google.com/go/auth"
	"github.com/getsentry/sentry-go"
	"github.com/gin-gonic/gin"
)

// FirebaseAuthMiddleware structure
type FirebaseAuthMiddleware struct {
	service     services.FirebaseService
	userservice services.UserService
}

// NewFirebaseAuthMiddleware creates new firebase authentication
func NewFirebaseAuthMiddleware(
	service services.FirebaseService,
	userservice services.UserService,
) FirebaseAuthMiddleware {
	return FirebaseAuthMiddleware{
		service:     service,
		userservice: userservice,
	}
}

// Handle handles auth requests
func (m FirebaseAuthMiddleware) Handle() gin.HandlerFunc {
	return func(c *gin.Context) {
		token, err := m.getTokenFromHeader(c)
		if err != nil {
			responses.ErrorJSON(c, http.StatusUnauthorized, err.Error())
			c.Abort()
			return
		}
		id := int64(0) 
		if token.Claims[constants.DBUSERID] != nil{
			getIDFromClaims, err:= strconv.ParseInt(token.Claims[constants.DBUSERID].(string), 10, 64)
			if err != nil || getIDFromClaims ==0{
				responses.ErrorJSON(c, http.StatusUnauthorized, "Unauthorized user")
				c.Abort()
				return
			}
			id =getIDFromClaims
		}else{
			user, err := m.userservice.GetOneUserByFireBaseUID(token.UID)
			if err != nil {
				responses.ErrorJSON(c, http.StatusUnauthorized, "UnAuthorized User")
				c.Abort()
				return
			}
			id =user.ID
		}

		sentry.ConfigureScope(func(scope *sentry.Scope) {
			scope.SetUser(sentry.User{ID: token.UID})
		})

		c.Set(constants.Claims, token.Claims)
		c.Set(constants.UID,id)
		c.Next()
	}
}

// HandleAdminOnly handles middleware for admin role only
func (m FirebaseAuthMiddleware) HandleAdminOnly() gin.HandlerFunc {
	return func(c *gin.Context) {
		token, err := m.getTokenFromHeader(c)

		if err != nil {
			responses.ErrorJSON(c, http.StatusUnauthorized, err.Error())
			c.Abort()
			return
		}

		if !m.isAdmin(token.Claims) {
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
func (m FirebaseAuthMiddleware) getTokenFromHeader(c *gin.Context) (*auth.Token, error) {
	header := c.GetHeader("Authorization")
	idToken := strings.TrimSpace(strings.Replace(header, "Bearer", "", 1))

	token, err := m.service.VerifyToken(idToken)
	if err != nil {
		return nil, err
	}

	return token, nil
}

// isAdmin check if cliams is admin
func (M FirebaseAuthMiddleware) isAdmin(claims map[string]interface{}) bool {

	role := claims["role"]
	isAdmin := false
	if role != nil {
		isAdmin = role.(string) == "admin"
	}

	return isAdmin

}

// Handle handles auth requests if token provided
func (m FirebaseAuthMiddleware) HandleUser() gin.HandlerFunc {
	return func(c *gin.Context) {
		var UID int64
		c.Set(constants.UID, UID)
		token, _ := m.getTokenFromHeader(c)

		if token != nil {			
			if token.Claims[constants.DBUSERID]!= nil{
				getIDFromClaims, err:= strconv.ParseInt(token.Claims[constants.DBUSERID].(string), 10, 64)
				if err !=nil || getIDFromClaims== 0{
					responses.ErrorJSON(c, http.StatusUnauthorized, "UnAuthorized User")
					c.Abort()
					return
				}
				c.Set(constants.UID,getIDFromClaims)
				c.Next()
			}
			user, err := m.userservice.GetOneUserByFireBaseUID(token.UID)
			if err != nil {
				responses.ErrorJSON(c, http.StatusUnauthorized, "UnAuthorized User")
				c.Abort()
				return
			}
			c.Set(constants.UID,user.ID)
		}
		c.Next()
	}
}

