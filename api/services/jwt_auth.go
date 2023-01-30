package services

import (
	"boilerplate-api/errors"
	"boilerplate-api/infrastructure"
	"strings"

	"github.com/getsentry/sentry-go"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
)

type JWTAuthService struct {
	logger infrastructure.Logger
	env    infrastructure.Env
}

func NewJWTAuthService(
	logger infrastructure.Logger,
	env infrastructure.Env,
) JWTAuthService {
	return JWTAuthService{
		logger: logger,
		env:    env,
	}
}

type JWTClaims struct {
	jwt.StandardClaims
	// ...other claims
}

func (m JWTAuthService) ParseToken(tokenString string) (*jwt.Token, error) {
	// Parse the token using the secret key
	token, err := jwt.ParseWithClaims(tokenString, &JWTClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(m.env.JWT_SECRET), nil
	})
	if err != nil {
		if !strings.Contains(err.Error(), "expired") {
			m.logger.Zap.Error("Invalid token[ParseWithClaims] :", err.Error())
			err := errors.BadRequest.New("Invalid ID token")
			return nil, err
		}
		m.logger.Zap.Error("Invalid token[ParseWithClaims] :", err.Error())
		return nil, err
	}
	return token, nil

}

func (m JWTAuthService) VerifyToken(c *gin.Context) (bool, error) {
	// Get the token from the request header
	header := c.GetHeader("Authorization")
	if header == "" {
		err := errors.BadRequest.New("Authorization token is required in header")
		err = errors.SetCustomMessage(err, "Authorization token is required in header")
		m.logger.Zap.Error("[GetHeader]: ", err.Error())
		return false, err
	}

	if !strings.Contains(header, "Bearer") {
		err := errors.BadRequest.New("Token type is required")
		m.logger.Zap.Error("Missing token type: ", err.Error())
		return false, err
	}

	tokenString := strings.TrimSpace(strings.Replace(header, "Bearer", "", 1))
	token, err := m.ParseToken(tokenString)
	if err != nil {
		m.logger.Zap.Error("Error parsing token", err.Error())
		return false, err
	}

	claims, ok := token.Claims.(*JWTClaims)
	if !ok || !token.Valid {
		err := errors.BadRequest.New("Invalid token")
		err = errors.SetCustomMessage(err, "Invalid token")
		m.logger.Zap.Error("Invalid token [token.Valid]: ", err.Error())
		return false, err
	}
	sentry.ConfigureScope(func(scope *sentry.Scope) {
		scope.SetUser(sentry.User{ID: claims.Id})
	})
	// Can set anything in the request context and passes the request to the next handler.
	c.Set("user_id", claims.Id)
	return true, nil

}
