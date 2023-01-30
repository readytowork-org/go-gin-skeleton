package controllers

import (
	"boilerplate-api/api/middlewares"
	"boilerplate-api/api/responses"
	"boilerplate-api/api/services"
	"boilerplate-api/api/validators"
	"boilerplate-api/errors"
	"boilerplate-api/infrastructure"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
	"golang.org/x/crypto/bcrypt"
)

// JwtAuthController -> struct
type JwtAuthController struct {
	logger      infrastructure.Logger
	userService services.UserService
	env         infrastructure.Env
	validator   validators.UserValidator
}

// NewJwtAuthController -> constructor
func NewJwtAuthController(
	logger infrastructure.Logger,
	userService services.UserService,
	env infrastructure.Env,
	validator validators.UserValidator,
) JwtAuthController {
	return JwtAuthController{
		logger:      logger,
		userService: userService,
		env:         env,
		validator:   validator,
	}
}

func (cc JwtAuthController) ObtainJwtToken(c *gin.Context) {
	var reqData struct {
		Email    string `json:"email" validate:"required"`
		Password string `json:"password" validate:"required"`
	}
	// Bind the request payload to a reqData struct
	if err := c.ShouldBindJSON(&reqData); err != nil {
		cc.logger.Zap.Error("Error [ShouldBindJSON] : ", err.Error())
		err := errors.InternalError.Wrap(err, "Failed to bind request data")
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
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(reqData.Password)); err != nil {
		cc.logger.Zap.Error("[CompareHashAndPassword] hash and plain password doesnot match")
		responses.ErrorJSON(c, http.StatusBadRequest, "Invalid user credentials")
		return
	}

	// Create a new JWT claims object
	claims := middlewares.JWTClaims{
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Minute * time.Duration(cc.env.JWT_TOKEN_EXPIRES_AT)).Unix(),
			Id:        fmt.Sprintf("%v", user.ID),
		},
	}

	// Create a new JWT token using the claims and the secret key
	tokenClaim := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	token, tokenErr := tokenClaim.SignedString([]byte(cc.env.JWT_SECRET))
	if tokenErr != nil {
		cc.logger.Zap.Error("[SignedString] Error getting token: ", tokenErr.Error())
		responses.ErrorJSON(c, http.StatusInternalServerError, tokenErr.Error())
		return
	}
	data := map[string]interface{}{
		"user":       user.ToMap(),
		"token":      token,
		"expires_at": claims.ExpiresAt,
	}
	responses.SuccessJSON(c, http.StatusOK, data)
	return
}
