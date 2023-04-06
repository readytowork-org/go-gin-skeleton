package controllers

import (
	"boilerplate-api/api/responses"
	"boilerplate-api/api/services"
	"boilerplate-api/api/validators"
	"boilerplate-api/errors"
	"boilerplate-api/infrastructure"
	"boilerplate-api/models"
	"boilerplate-api/utils"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

type UserProfileController struct {
	logger             infrastructure.Logger
	userProfileService services.UserProfileService
	validator          validators.UserProfileValidator
}

func NewUserProfileController(logger infrastructure.Logger, userProfileService services.UserProfileService, validator validators.UserProfileValidator) UserProfileController {
	return UserProfileController{
		logger:             logger,
		userProfileService: userProfileService,
		validator:          validator,
	}
}

func (cc UserProfileController) GetAllUserProfile(c *gin.Context) {
	pagination := utils.BuildPagination(c)
	userProfiles, count, err := cc.userProfileService.GetAllUserProfile(pagination)
	if err != nil {
		cc.logger.Zap.Error("Error Finding userProfiles", err)
		err := errors.InternalError.Wrap(err, "Failed to get userProfiles")
		responses.HandleError(c, err)
		return
	}
	responses.JSONCount(c, http.StatusOK, userProfiles, count)
}

func (cc UserProfileController) CreateUserProfile(c *gin.Context) {
	userProfile := models.UserProfile{}
	c.ShouldBindJSON(&userProfile)
	if validationErr := cc.validator.Validate.Struct(userProfile); validationErr != nil {
		err := errors.BadRequest.Wrap(validationErr, "Validation error")
		err = errors.SetCustomMessage(err, "Invalid input information")
		err = errors.AddErrorContextBlock(err, cc.validator.GenerateValidationResponse(validationErr))
		responses.HandleError(c, err)
		return
	}
	created, err := cc.userProfileService.CreateUserProfile(userProfile)
	cc.logger.Zap.Info("got err?????0", err)
	if err != nil {
		if strings.Contains(err.Error(), "1062") {
			cc.logger.Zap.Error("----1062------")
			err = errors.BadRequest.Wrap(err, "Error creating user")
			custom_msg := ""
			if strings.Contains(err.Error(), "UQ_user_profile_user_id") {
				cc.logger.Zap.Info("iside user id string----------")
				custom_msg = "User Profile already exists."
			} else if strings.Contains(err.Error(), "UQ_user_profile_contact") {
				cc.logger.Zap.Info("iside phone number string----------")
				custom_msg = "Phone number already taken"
			}
			err = errors.SetCustomMessage(err, custom_msg)
			responses.ErrorJSON(c, http.StatusBadRequest, custom_msg)
			return

		}
		if strings.Contains(err.Error(), "1452") {
			cc.logger.Zap.Error("----1062------")
			err = errors.BadRequest.Wrap(err, "Error creating user")
			custom_msg := ""
			if strings.Contains(err.Error(), "user_profile_ibfk_user_id") {
				cc.logger.Zap.Info("Invalid user id----------")
				custom_msg = "Invalid user id."
			} else {
				cc.logger.Zap.Info("Invalid user id----------")
				custom_msg = "Invalid user id."
			}
			err = errors.SetCustomMessage(err, custom_msg)
			// responses.ErrorJSON(c, http.StatusBadRequest, custom_msg)
			// return

		}
		cc.logger.Zap.Error("Error Creating UserProfile", err.Error(), created)
		err = errors.InternalError.Wrap(err, "Failed to create userrofile")
		responses.HandleError(c, err)
		return
	}
	responses.SuccessJSON(c, http.StatusCreated, created)
}

func (cc UserProfileController) GetUserProfile(c *gin.Context) {
	if c.Param("id") == "" {
		responses.JSON(c, http.StatusBadRequest, "Id required in url")
		return
	}
	userProfile, err := cc.userProfileService.GetUserProfile(c.Param("id"))
	if err != nil {
		cc.logger.Zap.Error("Error Creating UserProfile", err.Error())
		// err := errors.InternalError.Wrap(err, "Failed to create userrofile")
		responses.HandleError(c, err)
		return
	}
	responses.SuccessJSON(c, http.StatusCreated, userProfile)

}
