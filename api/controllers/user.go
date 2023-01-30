package controllers

import (
	"boilerplate-api/api/responses"
	"boilerplate-api/api/services"
	"boilerplate-api/api/validators"
	"boilerplate-api/constants"
	"boilerplate-api/errors"
	"boilerplate-api/infrastructure"
	"boilerplate-api/models"
	"boilerplate-api/utils"
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// UserController -> struct
type UserController struct {
	logger      infrastructure.Logger
	userService services.UserService
	env         infrastructure.Env
	validator   validators.UserValidator
}

// NewUserController -> constructor
func NewUserController(
	logger infrastructure.Logger,
	userService services.UserService,
	env infrastructure.Env,
	validator validators.UserValidator,
) UserController {
	return UserController{
		logger:      logger,
		userService: userService,
		env:         env,
		validator:   validator,
	}
}

// CreateUser -> Create User
func (cc UserController) CreateUser(c *gin.Context) {
	reqData := struct {
		models.User
		ConfirmPassword string `json:"confirm_password" validate:"required"`
	}{}
	trx := c.MustGet(constants.DBTransaction).(*gorm.DB)

	if err := c.ShouldBindJSON(&reqData); err != nil {
		cc.logger.Zap.Error("Error [CreateUser] (ShouldBindJson) : ", err)
		err := errors.BadRequest.Wrap(err, "Failed to bind user data")
		responses.HandleError(c, err)
		return
	}
	if validationErr := cc.validator.Validate.Struct(reqData); validationErr != nil {
		err := errors.BadRequest.Wrap(validationErr, "Validation error")
		err = errors.SetCustomMessage(err, "Invalid input information")
		err = errors.AddErrorContextBlock(err, cc.validator.GenerateValidationResponse(validationErr))
		responses.HandleError(c, err)
		return
	}

	if reqData.User.Password != reqData.ConfirmPassword {
		cc.logger.Zap.Error("Password and confirm password not matching : ")
		responses.ErrorJSON(c, http.StatusBadRequest, "Password and confirm password should be same.")
		return
	}

	if userEmail, _ := cc.userService.GetOneUserWithEmail(reqData.Email); userEmail != nil {
		cc.logger.Zap.Error("Error [CreateUser] [db CreateUser]: User with this email already exists")
		responses.ErrorJSON(c, http.StatusBadRequest, "User with this email already exists")
		return
	}

	if userContact, _ := cc.userService.GetOneUserWithPhone(reqData.Phone); userContact != nil {
		cc.logger.Zap.Error("Error [db GetOneUserWithPhone]: User with this phone already exists")
		responses.ErrorJSON(c, http.StatusBadRequest, "User with this phone already exists")
		return
	}

	if err := cc.userService.WithTrx(trx).CreateUser(reqData.User); err != nil {
		cc.logger.Zap.Error("Error [CreateUser] [db CreateUser]: ", err.Error())
		err := errors.InternalError.Wrap(err, "Failed to create user")
		responses.HandleError(c, err)
		return
	}

	responses.SuccessJSON(c, http.StatusOK, "User Created Sucessfully")
}

// GetAllUser -> Get All User
func (cc UserController) GetAllUsers(c *gin.Context) {
	pagination := utils.BuildPagination(c)
	users, count, err := cc.userService.GetAllUsers(pagination)

	if err != nil {
		cc.logger.Zap.Error("Error finding user records", err.Error())
		err := errors.InternalError.Wrap(err, "Failed to get users data")
		responses.HandleError(c, err)
		return
	}

	responses.JSONCount(c, http.StatusOK, users, count)
}
