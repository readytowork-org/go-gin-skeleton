package user

import (
	"errors"
	"net/http"

	"boilerplate-api/lib/api_errors"
	"boilerplate-api/lib/config"
	"boilerplate-api/lib/constants"
	"boilerplate-api/lib/json_response"
	"boilerplate-api/lib/request_validator"
	"boilerplate-api/lib/utils"
	"github.com/gin-gonic/gin"

	"gorm.io/gorm"
)

type Controller struct {
	logger      config.Logger
	userService Service
	env         config.Env
	validator   request_validator.Validator
}

// NewController Creates New user controller
func NewController(
	logger config.Logger,
	userService Service,
	env config.Env,
	validator request_validator.Validator,
) Controller {
	return Controller{
		logger:      logger,
		userService: userService,
		env:         env,
		validator:   validator,
	}
}

// @Tags			UserManagementApi
// @Summary		Create CreateUser
// @Description	Create one user
// @Security		Bearer
// @Produce		application/json
// @Param			data	body		CreateUserRequestData	true	"Enter JSON"
// @Success		200		{object}	json_response.Message	"CUser Created Successfully"
// @Failure		400		{object}	api_errors.Envelope
// @Failure		409		{object}	api_errors.Envelope
// @Failure		422		{object}	api_errors.Envelope
// @Failure		500		{object}	api_errors.Envelope
// @Router			/api/v1/users [post]
// @Id				CreateUser
func (cc Controller) CreateUser(c *gin.Context) {
	reqData := CreateUserRequestData{}
	trx := c.MustGet(constants.DBTransaction).(*gorm.DB)

	if err := c.ShouldBindJSON(&reqData); err != nil {
		api_errors.RespondError(c, api_errors.Wrap(err, http.StatusBadRequest, api_errors.CodeBadRequest, "Failed to bind user data"))
		return
	}
	if validationErr := cc.validator.Struct(reqData); validationErr != nil {
		api_errors.RespondError(c, api_errors.WithValidation(cc.validator.GenerateValidationResponse(validationErr), "Invalid input information"))
		return
	}

	if err := cc.userService.WithTrx(trx).CreateUser(reqData); err != nil {
		switch {
		case errors.Is(err, ErrPasswordMismatch):
			api_errors.RespondError(c, api_errors.New(http.StatusBadRequest, api_errors.CodeBadRequest, "Password and confirm password should be same."))
		case errors.Is(err, ErrEmailAlreadyExists):
			api_errors.RespondError(c, api_errors.New(http.StatusConflict, api_errors.CodeConflict, "User with this email already exists"))
		case errors.Is(err, ErrPhoneAlreadyExists):
			api_errors.RespondError(c, api_errors.New(http.StatusConflict, api_errors.CodeConflict, "User with this phone already exists"))
		default:
			cc.logger.Error("Error [CreateUser]: ", err.Error())
			api_errors.RespondError(c, api_errors.Wrap(err, http.StatusInternalServerError, api_errors.CodeInternal, "Failed to create user"))
		}
		return
	}

	c.JSON(
		http.StatusOK, json_response.Message{
			Msg: "CUser Created Successfully",
		},
	)
}

// @Tags			UserManagementApi
// @Summary		All users
// @Description	get all users
// @Security		Bearer
// @Produce		application/json
// @Param			pagination	query		Pagination	false	"query param"
// @Success		200			{object}	json_response.DataCount[GetUserResponse]
// @Failure		500			{object}	api_errors.Envelope
// @Router			/api/v1/users [get]
// @Id				GetAllUsers
func (cc Controller) GetAllUsers(c *gin.Context) {
	pagination := utils.BuildPagination[*Pagination](c)

	users, count, err := cc.userService.GetAllUsers(*pagination)
	if err != nil {
		api_errors.RespondError(c, api_errors.Wrap(err, http.StatusInternalServerError, api_errors.CodeInternal, "Failed to get users data"))
		return
	}

	c.JSON(
		http.StatusOK, json_response.DataCount[GetUserResponse]{
			Count: count,
			Data:  users,
		},
	)
}

// @Tags			UserManagementApi
// @Summary		CreateUser Profile
// @Description	get user profile
// @Security		Bearer
// @Produce		application/json
// @Success		200	{object}	json_response.Data[GetUserResponse]
// @Failure		500	{object}	api_errors.Envelope
// @Router			/api/v1/{id} [get]
// @Id				GetOneUser
func (cc Controller) GetOneUser(c *gin.Context) {
	userID, errResponse := utils.StringToInt64(c.Param("id"))
	if errResponse != nil {
		api_errors.RespondError(c, api_errors.New(http.StatusInternalServerError, api_errors.CodeInternal, errResponse.Message))
		return
	}

	user, err := cc.userService.GetOneUser(userID)
	if err != nil {
		api_errors.RespondError(c, api_errors.Wrap(err, http.StatusInternalServerError, api_errors.CodeInternal, "Failed to get user"))
		return
	}

	c.JSON(
		http.StatusOK, json_response.Data[GetUserResponse]{
			Data: user,
		},
	)
}
