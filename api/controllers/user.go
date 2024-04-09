package controllers

import (
	"fmt"
	"net/http"

	"boilerplate-api/api/services"
	"boilerplate-api/api/user"
	"boilerplate-api/dtos"
	"boilerplate-api/internal/api_errors"
	"boilerplate-api/internal/config"
	"boilerplate-api/internal/constants"
	"boilerplate-api/internal/json_response"
	"boilerplate-api/internal/request_validator"
	"boilerplate-api/internal/utils"
	"github.com/gin-gonic/gin"

	"gorm.io/gorm"
)

type UserController struct {
	logger      config.Logger
	userService services.UserService
	env         config.Env
	validator   request_validator.Validator
}

// NewUserController Creates New user controller
func NewUserController(
	logger config.Logger,
	userService services.UserService,
	env config.Env,
	validator request_validator.Validator,
) UserController {
	return UserController{
		logger:      logger,
		userService: userService,
		env:         env,
		validator:   validator,
	}
}

//	@Summary		Create User
//	@Description	Create User
//	@Param			data			body	dtos.CreateUserRequestData	true	"Enter JSON"
//	@Param			Authorization	header	string						true	"Insert your access token"	default(Bearer <Add access token here>)
//	@Produce		application/json
//	@Tags			User
//	@Success		200	{object}	api_response.Success	"OK"
//	@Failure		400	{object}	json_response.Error
//	@Failure		500	{object}	json_response.Error
//	@Router			/users [post]
//	@Id				CreateUser
func (cc UserController) CreateUser(c *gin.Context) {
	reqData := dtos.CreateUserRequestData{}
	trx := c.MustGet(constants.DBTransaction).(*gorm.DB)

	if err := c.ShouldBindJSON(&reqData); err != nil {
		cc.logger.Error("Error [CreateUser] (ShouldBindJson) : ", err)
		err := api_errors.BadRequest.Wrap(err, "Failed to bind user data")
		status, errM := api_errors.HandleError(err)
		c.JSON(status, json_response.Error{Error: errM})
		return
	}
	if validationErr := cc.validator.Struct(reqData); validationErr != nil {
		err := api_errors.BadRequest.Wrap(validationErr, "Validation error")
		err = api_errors.SetCustomMessage(err, "Invalid input information")
		err = api_errors.AddErrorContextBlock(err, cc.validator.GenerateValidationResponse(validationErr))
		status, errM := api_errors.HandleError(err)
		c.JSON(status, json_response.Error{Error: errM})
		return
	}

	if reqData.Password != reqData.ConfirmPassword {
		cc.logger.Error("Password and confirm password not matching : ")
		err := api_errors.BadRequest.New("Password and confirm password should be same.")
		status, errM := api_errors.HandleError(err)
		c.JSON(status, json_response.Error{Error: errM})
		return
	}

	if _, err := cc.userService.GetOneUserWithEmail(reqData.Email); err != nil {
		cc.logger.Error("Error [CreateUser] [db CreateUser]: User with this email already exists")
		err := api_errors.BadRequest.New("User with this email already exists")
		status, errM := api_errors.HandleError(err)
		c.JSON(status, json_response.Error{Error: errM})
		return
	}

	if _, err := cc.userService.GetOneUserWithPhone(reqData.Phone); err != nil {
		cc.logger.Error("Error [db GetOneUserWithPhone]: User with this phone already exists")
		err := api_errors.BadRequest.New("User with this phone already exists")
		status, errM := api_errors.HandleError(err)
		c.JSON(status, json_response.Error{Error: errM})
		return
	}

	if err := cc.userService.WithTrx(trx).CreateUser(reqData.User); err != nil {
		cc.logger.Error("Error [CreateUser] [db CreateUser]: ", err.Error())
		err := api_errors.InternalError.Wrap(err, "Failed to create user")
		status, errM := api_errors.HandleError(err)
		c.JSON(status, json_response.Error{Error: errM})
		return
	}

	c.JSON(http.StatusOK, json_response.Message{
		Msg: "User Created Successfully",
	})
}

//	@Summary		Get all User.
//	@Param			page_size	query	string	false	"10"
//	@Param			page		query	string	false	"Page no"	"1"
//	@Param			keyword		query	string	false	"search by name"
//	@Param			Keyword2	query	string	false	"search by type"
//	@Description	Return all the User
//	@Produce		application/json
//	@Param			Authorization	header	string	true	"Insert your access token"	default(Bearer <Add access token here>)
//	@Tags			User
//	@Success		200	{object}	json_response.DataCount[dtos.GetUserResponse]
//	@Failure		500	{object}	json_response.Error
//	@Router			/users [get]
//	@Id				GetAllUsers
func (cc UserController) GetAllUsers(c *gin.Context) {
	pagination := utils.BuildPagination[*user.Pagination](c)

	users, count, err := cc.userService.GetAllUsers(*pagination)
	if err != nil {
		cc.logger.Error("Error finding user records", err.Error())
		err := api_errors.InternalError.Wrap(err, "Failed to get users data")
		status, errM := api_errors.HandleError(err)
		c.JSON(status, json_response.Error{Error: errM})
		return
	}

	c.JSON(http.StatusOK, json_response.DataCount[dtos.GetUserResponse]{
		Count: count,
		Data:  users,
	})
}

//	@Summary		Get one user by id
//	@Description	Get one user by id
//	@Param			Authorization	header	string	true	"Insert your access token"	default(Bearer <Add access token here>)
//	@Produce		application/json
//	@Tags			User
//	@Success		200	{object}	json_response.Data[dtos.GetUserResponse]
//	@Failure		500	{object}	json_response.Error
//	@Router			/profile [get]
//	@Id				GetUserProfile
func (cc UserController) GetUserProfile(c *gin.Context) {
	userID := fmt.Sprintf("%v", c.MustGet(constants.UserID))

	user, err := cc.userService.GetOneUser(userID)
	if err != nil {
		cc.logger.Error("Error finding user profile", err.Error())
		err := api_errors.InternalError.Wrap(err, "Failed to get users profile data")
		status, errM := api_errors.HandleError(err)
		c.JSON(status, json_response.Error{Error: errM})
		return
	}

	c.JSON(http.StatusOK, json_response.Data[dtos.GetUserResponse]{
		Data: user,
	})
}
