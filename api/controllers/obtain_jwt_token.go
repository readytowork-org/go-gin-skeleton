package controllers

import (
	"boilerplate-api/api/responses"
	"boilerplate-api/api/services"
	"boilerplate-api/api/validators"
	"boilerplate-api/dtos"
	"boilerplate-api/errors"
	"boilerplate-api/infrastructure"
	"boilerplate-api/utils"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
)

// JwtAuthController -> struct
type JwtAuthController struct {
	logger      infrastructure.Logger
	userService services.UserService
	jwtService  services.JWTAuthService
	env         infrastructure.Env
	validator   validators.UserValidator
}

// NewJwtAuthController -> constructor
func NewJwtAuthController(
	logger infrastructure.Logger,
	userService services.UserService,
	jwtService services.JWTAuthService,
	env infrastructure.Env,
	validator validators.UserValidator,
) JwtAuthController {
	return JwtAuthController{
		logger:      logger,
		userService: userService,
		jwtService:  jwtService,
		env:         env,
		validator:   validator,
	}
}

func (cc JwtAuthController) ObtainJwtToken(c *gin.Context) {
	reqData := dtos.JWTLoginRequestData{}
	// Bind the request payload to a reqData struct
	if err := c.ShouldBindJSON(&reqData); err != nil {
		cc.logger.Zap.Error("Error [ShouldBindJSON] : ", err.Error())
		err := errors.BadRequest.Wrap(err, "Failed to bind request data")
		responses.HandleError(c, err)
		return
	}
	// validating using custom validator
	if validationErr := cc.validator.Validate.Struct(reqData); validationErr != nil {
		cc.logger.Zap.Error("[Validate Struct] Validation error: ", validationErr.Error())
		err := errors.BadRequest.Wrap(validationErr, "Validation error")
		err = errors.SetCustomMessage(err, "Invalid input information")
		err = errors.AddErrorContextBlock(err, cc.validator.GenerateValidationResponse(validationErr))
		responses.HandleError(c, err)
		return
	}
	// Check if the user exists with provided email address
	user, err := cc.userService.GetOneUserWithEmail(reqData.Email)
	if err != nil {
		responses.ErrorJSON(c, http.StatusBadRequest, "Invalid user credentials1")
		return
	}

	// Check if the password is correct
	// Thus password is encrypted and saved in DB, comparing plain text with it's hash
	isValidPassword := utils.CompareHashAndPlainPassword(user.Password, reqData.Password)
	if !isValidPassword {
		cc.logger.Zap.Error("[CompareHashAndPassword] hash and plain password doesnot match")
		responses.ErrorJSON(c, http.StatusBadRequest, "Invalid user credentials")
		return
	}

	// Create a new JWT access claims object
	accessClaims := services.JWTClaims{
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Minute * time.Duration(cc.env.JWT_ACCESS_TOKEN_EXPIRES_AT)).Unix(),
			Id:        fmt.Sprintf("%v", user.ID),
		},
		//Add other claims
	}
	// Create a new JWT Access token using the claims and the secret key
	accessToken, tokenErr := cc.jwtService.GenerateToken(accessClaims, cc.env.JWT_ACCESS_SECRET)
	if tokenErr != nil {
		cc.logger.Zap.Error("[SignedString] Error getting token: ", tokenErr.Error())
		responses.ErrorJSON(c, http.StatusInternalServerError, tokenErr.Error())
		return
	}
	// Create a new JWT refresh claims object
	refreshClaims := services.JWTClaims{
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Hour * time.Duration(cc.env.JWT_REFRESH_TOKEN_EXPIRES_AT)).Unix(),
			Id:        fmt.Sprintf("%v", user.ID),
		},
	}
	// Create a new JWT Refresh token using the claims and the secret key
	refreshToken, refreshTokenErr := cc.jwtService.GenerateToken(refreshClaims, cc.env.JWT_REFRESH_SECRET)
	if refreshTokenErr != nil {
		cc.logger.Zap.Error("[SignedString] Error getting token: ", refreshTokenErr.Error())
		responses.ErrorJSON(c, http.StatusInternalServerError, refreshTokenErr.Error())
		return
	}
	data := map[string]interface{}{
		"user":               user.ToMap(),
		"access_token":       accessToken,
		"refresh_token":      refreshToken,
	}
	responses.SuccessJSON(c, http.StatusOK, data)
}

func (cc JwtAuthController) RefreshJwtToken(c *gin.Context) {
	tokenString, err := cc.jwtService.GetTokenFromHeader(c)
	if err != nil {
		cc.logger.Zap.Error("Error getting token from header: ", err.Error())
		err = errors.Unauthorized.Wrap(err, "Something went wrong")
		responses.HandleError(c, err)
		return
	}
	parsedToken, parseErr := cc.jwtService.ParseToken(tokenString, cc.env.JWT_REFRESH_SECRET)
	if parseErr != nil {
		cc.logger.Zap.Error("Error parsing token: ", parseErr.Error())
		err = errors.Unauthorized.Wrap(parseErr, "Something went wrong")
		responses.HandleError(c, err)
		return
	}
	claims, verifyErr := cc.jwtService.VerifyToken(parsedToken)
	if verifyErr != nil {
		cc.logger.Zap.Error("Error veriefying token: ", verifyErr.Error())
		err = errors.Unauthorized.Wrap(verifyErr, "Something went wrong")
		responses.HandleError(c, err)
		return
	}
	// Create a new JWT Access claims
	accessClaims := services.JWTClaims{
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Minute * time.Duration(cc.env.JWT_ACCESS_TOKEN_EXPIRES_AT)).Unix(),
			Id:        fmt.Sprintf("%v", claims.Id),
		},
		// Add other claims
	}
	// Create a new JWT token using the claims and the secret key
	accessToken, tokenErr := cc.jwtService.GenerateToken(accessClaims, cc.env.JWT_ACCESS_SECRET)
	if tokenErr != nil {
		cc.logger.Zap.Error("[SignedString] Error getting token: ", tokenErr.Error())
		responses.ErrorJSON(c, http.StatusInternalServerError, tokenErr.Error())
		return
	}
	data := map[string]interface{}{
		"access_token": accessToken,
		"expires_at":   accessClaims.ExpiresAt,
	}
	responses.SuccessJSON(c, http.StatusOK, data)
	return

}
