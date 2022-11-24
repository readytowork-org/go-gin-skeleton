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
	logger          infrastructure.Logger
	userService     services.UserService
	env             infrastructure.Env
	validator       validators.UserValidator
	firebaseService services.FirebaseService
}

// NewUserController -> constructor
func NewUserController(
	logger infrastructure.Logger,
	userService services.UserService,
	env infrastructure.Env,
	validator validators.UserValidator,
	firebaseService services.FirebaseService,
) UserController {
	return UserController{
		logger:          logger,
		userService:     userService,
		env:             env,
		validator:       validator,
		firebaseService: firebaseService,
	}
}

// CreateUser -> Create User
func (cc UserController) CreateUser(c *gin.Context) {
	reqData := struct {
		models.User
		Password        string `json:"password" validate:"required"`
		ConfirmPassword string `json:"confirm_password" validate:"required"`
	}{}
	trx := c.MustGet(constants.DBTransaction).(*gorm.DB)
	if err := c.ShouldBindJSON(&reqData); err != nil {
		cc.logger.Zap.Error("Error [CreateUser] (ShouldBindJson) : ", err)
		err := errors.BadRequest.Wrap(err, "Failed to bind user data")
		responses.HandleError(c, err)
		return
	}
	if reqData.Password != reqData.ConfirmPassword {
		cc.logger.Zap.Error("Password and confirm password not matching : ")
		responses.ErrorJSON(c, http.StatusBadRequest, "Password and confirm password should be same.")
		return
	}
	if !utils.IsValidEmail(reqData.User.Email) {
		cc.logger.Zap.Error("Invalid email")
		responses.ErrorJSON(c, http.StatusBadRequest, "Invalid Email")
		return
	}
	if validationErr := cc.validator.Validate.Struct(reqData); validationErr != nil {
		err := errors.BadRequest.Wrap(validationErr, "Validation error")
		err = errors.SetCustomMessage(err, "Invalid input information")
		err = errors.AddErrorContextBlock(err, cc.validator.GenerateValidationResponse(validationErr))
		responses.HandleError(c, err)
		return
	}

	fb_user := cc.firebaseService.GetUserByEmail(reqData.Email)
	if fb_user != "" {
		err := errors.BadRequest.New("Firebase user already exists")
		err = errors.SetCustomMessage(err, "Email address already taken")
		responses.HandleError(c, err)
		return
	}
	created_user, err := cc.userService.WithTrx(trx).CreateUser(&reqData.User)
	if err != nil {
		cc.logger.Zap.Error("Error [CreateUser] [db CreateUser]: ", err.Error())
		err := errors.InternalError.Wrap(err, "Failed to create user")
		responses.HandleError(c, err)
		return
	}
	fb_auth_user := models.FirebaseAuthUser{
		Enabled:     1,
		Email:       created_user.Email,
		DisplayName: created_user.FullName,
		Password:    reqData.Password,
		Role:        constants.RoleUser,
		UserId:      utils.Int64ToString(created_user.ID),
	}
	fb_uid, err := cc.firebaseService.CreateUser(fb_auth_user)
	if err != nil {
		cc.logger.Zap.Error("Error creating client in firebase: ", err.Error())
		err = errors.InternalError.Wrap(err, "Error creating client in firebase")
		responses.HandleError(c, err)
		return
	}
	updated_user, err := cc.userService.WithTrx(trx).UpdatePartial(
		created_user.ID, map[string]interface{}{
			"firebase_uid": fb_uid,
		})
	if err != nil {
		cc.logger.Zap.Error("Error Updating user: ", err.Error())
		err := errors.InternalError.Wrap(err, "Error deleting firebase user")
		responses.HandleError(c, err)
		return
	}
	user_data := map[string]interface{}{
		"id":           updated_user.ID,
		"firebase_uid": updated_user.FirebaseUID,
		"full_name":    updated_user.FullName,
		"email":        updated_user.Email,
		"username":     updated_user.Username,
		"phone":        updated_user.Phone,
		"address":      updated_user.Address,
	}

	responses.SuccessJSON(c, http.StatusOK, user_data)
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

func (cc UserController) GetOneUser(c *gin.Context) {
	id := c.Param("id")
	user, err := cc.userService.GetOneUser(id)
	if err != nil {
		cc.logger.Zap.Error("Error finding user records", err.Error())
		err := errors.InternalError.Wrap(err, "Failed to get users data")
		responses.HandleError(c, err)
		return
	}
	responses.SuccessJSON(c, http.StatusOK, user)
	return
}

func (cc UserController) DeleteOneUser(c *gin.Context) {
	trx := c.MustGet(constants.DBTransaction).(*gorm.DB)
	f_uid, err := cc.userService.WithTrx(trx).DeleteOneUser(c.Param("id"))
	if err != nil {
		cc.logger.Zap.Error("Error Deleting user record", err.Error())
		err := errors.InternalError.Wrap(err, "Failed to delete user data")
		responses.HandleError(c, err)
		return
	}
	if err := cc.firebaseService.DeleteUser(*f_uid); err != nil {
		cc.logger.Zap.Error("Error Deleting user record from firebase", err.Error())
		err := errors.InternalError.Wrap(err, "Failed to delete user data from firebase")
		responses.HandleError(c, err)
		return
	}
	responses.SuccessJSON(c, http.StatusOK, "User deleted successfully")
	return
}
