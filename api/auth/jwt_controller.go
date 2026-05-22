package auth

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"boilerplate-api/api/admin/user"
	"boilerplate-api/lib/api_errors"
	"boilerplate-api/lib/auth"
	"boilerplate-api/lib/config"
	"boilerplate-api/lib/constants"
	"boilerplate-api/lib/json_response"
	"boilerplate-api/lib/request_validator"
	"boilerplate-api/lib/types"
	"boilerplate-api/lib/utils"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
)

// JwtAuthController struct
type JwtAuthController struct {
	logger       config.Logger
	userService  user.Service
	jwtService   auth.JWTAuthService
	env          config.Env
	validator    request_validator.Validator
	refreshTokens RefreshTokenRepository
}

// NewJwtAuthController constructor
func NewJwtAuthController(
	logger config.Logger,
	userService user.Service,
	jwtService auth.JWTAuthService,
	env config.Env,
	validator request_validator.Validator,
	refreshTokens RefreshTokenRepository,
) JwtAuthController {
	return JwtAuthController{
		logger:        logger,
		userService:   userService,
		jwtService:    jwtService,
		env:           env,
		validator:     validator,
		refreshTokens: refreshTokens,
	}
}

// issueAccessAndRefresh signs a fresh access+refresh pair for the given user
// and persists the refresh token. Callers MUST treat any error here as a 500.
func (cc JwtAuthController) issueAccessAndRefresh(userID uint32) (accessTok, refreshTok string, accessExpiresAt time.Time, err error) {
	idStr := fmt.Sprintf("%v", userID)
	accessExpiresAt = time.Now().Add(time.Minute * time.Duration(cc.env.JwtAccessTokenExpiresAt))
	refreshExpiresAt := time.Now().Add(time.Hour * time.Duration(cc.env.JwtRefreshTokenExpiresAt))

	accessClaims := auth.JWTClaims{
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(accessExpiresAt),
			ID:        idStr,
		},
	}
	if accessTok, err = cc.jwtService.GenerateToken(accessClaims, cc.env.JwtAccessSecret); err != nil {
		return
	}

	refreshClaims := auth.JWTClaims{
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(refreshExpiresAt),
			ID:        idStr,
		},
	}
	if refreshTok, err = cc.jwtService.GenerateToken(refreshClaims, cc.env.JwtRefreshSecret); err != nil {
		return
	}

	if err = cc.refreshTokens.Store(userID, refreshTok, refreshExpiresAt); err != nil {
		return
	}
	return
}

func (cc JwtAuthController) LoginUserWithJWT(c *gin.Context) {
	reqData := JWTLoginRequestData{}
	if err := c.ShouldBindJSON(&reqData); err != nil {
		api_errors.RespondError(c, api_errors.Wrap(err, http.StatusBadRequest, api_errors.CodeBadRequest, "Failed to bind request data"))
		return
	}

	if validationErr := cc.validator.Struct(reqData); validationErr != nil {
		api_errors.RespondError(c, api_errors.WithValidation(cc.validator.GenerateValidationResponse(validationErr), "Invalid input information"))
		return
	}

	userData, err := cc.userService.GetOneUserWithEmail(reqData.Email)
	if err != nil {
		api_errors.RespondError(c, api_errors.New(http.StatusUnauthorized, api_errors.CodeUnauthorized, "Invalid user credentials"))
		return
	}

	if !utils.CompareHashAndPlainPassword(userData.Password, reqData.Password) {
		api_errors.RespondError(c, api_errors.New(http.StatusUnauthorized, api_errors.CodeUnauthorized, "Invalid user credentials"))
		return
	}

	accessTok, refreshTok, _, issueErr := cc.issueAccessAndRefresh(userData.ID)
	if issueErr != nil {
		api_errors.RespondError(c, api_errors.Wrap(issueErr, http.StatusInternalServerError, api_errors.CodeInternal, "Failed to issue tokens"))
		return
	}

	c.JSON(http.StatusOK, json_response.Data[types.MapString]{Data: types.MapString{
		"user":          userData,
		"access_token":  accessTok,
		"refresh_token": refreshTok,
	}})
}

// RefreshJwtToken rotates the caller's refresh token: the supplied token is
// validated and revoked, and a fresh access+refresh pair is issued and stored.
// Re-use of a revoked token is rejected with 401 (defence against token theft).
func (cc JwtAuthController) RefreshJwtToken(c *gin.Context) {
	header := c.GetHeader(constants.Headers.Authorization.ToString())

	tokenString, err := cc.jwtService.GetTokenFromHeader(header)
	if err != nil {
		api_errors.RespondError(c, api_errors.New(http.StatusUnauthorized, api_errors.CodeUnauthorized, err.Message))
		return
	}

	parsedToken, parseErr := cc.jwtService.ParseAndVerifyToken(tokenString, cc.env.JwtRefreshSecret)
	if parseErr != nil {
		api_errors.RespondError(c, api_errors.New(http.StatusUnauthorized, api_errors.CodeUnauthorized, parseErr.Message))
		return
	}

	claims, verifyErr := cc.jwtService.RetrieveClaims(parsedToken)
	if verifyErr != nil {
		api_errors.RespondError(c, api_errors.New(http.StatusUnauthorized, api_errors.CodeUnauthorized, verifyErr.Message))
		return
	}

	stored, lookupErr := cc.refreshTokens.FindActive(tokenString)
	if lookupErr != nil {
		if errors.Is(lookupErr, ErrRefreshTokenNotFound) || errors.Is(lookupErr, ErrRefreshTokenRevoked) {
			// If the token parsed but isn't in the active set, treat it as
			// compromise: revoke every refresh token for the claimed user.
			if userID, convErr := strconv.ParseUint(claims.ID, 10, 32); convErr == nil {
				_ = cc.refreshTokens.RevokeAllForUser(uint32(userID))
			}
			api_errors.RespondError(c, api_errors.New(http.StatusUnauthorized, api_errors.CodeUnauthorized, "Invalid refresh token"))
			return
		}
		api_errors.RespondError(c, api_errors.Wrap(lookupErr, http.StatusInternalServerError, api_errors.CodeInternal, "Failed to validate refresh token"))
		return
	}

	if err := cc.refreshTokens.Revoke(stored.ID); err != nil {
		api_errors.RespondError(c, api_errors.Wrap(err, http.StatusInternalServerError, api_errors.CodeInternal, "Failed to rotate refresh token"))
		return
	}

	accessTok, refreshTok, accessExp, issueErr := cc.issueAccessAndRefresh(stored.UserID)
	if issueErr != nil {
		api_errors.RespondError(c, api_errors.Wrap(issueErr, http.StatusInternalServerError, api_errors.CodeInternal, "Failed to issue tokens"))
		return
	}

	c.JSON(http.StatusOK, json_response.Data[types.MapString]{Data: types.MapString{
		"access_token":  accessTok,
		"refresh_token": refreshTok,
		"expires_at":    accessExp,
	}})
}

// Logout revokes the refresh token supplied in the Authorization header.
// Always returns 200 to avoid leaking token validity to clients.
func (cc JwtAuthController) Logout(c *gin.Context) {
	header := c.GetHeader(constants.Headers.Authorization.ToString())
	tokenString, err := cc.jwtService.GetTokenFromHeader(header)
	if err != nil {
		c.JSON(http.StatusOK, json_response.Message{Msg: "Logged out"})
		return
	}

	if stored, lookupErr := cc.refreshTokens.FindActive(tokenString); lookupErr == nil {
		_ = cc.refreshTokens.Revoke(stored.ID)
	}
	c.JSON(http.StatusOK, json_response.Message{Msg: "Logged out"})
}
