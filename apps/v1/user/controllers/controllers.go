package controllers

import (
	"boilerplate-api/apps/v1/user"
	"boilerplate-api/apps/v1/user/services"
	"boilerplate-api/common/constants"
	"boilerplate-api/common/errors"
	"boilerplate-api/common/helpers"
	"boilerplate-api/common/responses"
	"boilerplate-api/infrastructure"
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type UserController struct {
	logger    infrastructure.Logger
	service   services.UserService
	validator user.UserValidator
}

func UserControllerConstuctor(
	logger infrastructure.Logger,
	service services.UserService,
	validator user.UserValidator,

) UserController {
	return UserController{
		logger:    logger,
		service:   service,
		validator: validator,
	}
}

// CreateUser Create User
// @Summary				Create User
// @Description			Create User
// @Param				data body CreateUserRequestData true "Enter JSON"
// @Param 				Authorization header string true "Insert your access token" default(Bearer <Add access token here>)
// @Produce				application/json
// @Tags				User
// @Success				200 {object} responses.Success "OK"
// @Failure      		400 {object} responses.Error
// @Failure      		500 {object} responses.Error
// @Router				/users [post]
func (cc UserController) CreateUser(c *gin.Context) {
	reqData := user.CreateUserRequestData{}
	trx := c.MustGet(constants.DBTransaction).(*gorm.DB)

	if err := c.ShouldBindJSON(&reqData); err != nil {
		cc.logger.Zap.Error("Error [CreateUser] (ShouldBindJson) : ", err)
		err := errors.BadRequest.Wrap(err, "Failed to bind user data")
		responses.HandleError(c, err)
		return
	}
	if validationErr := cc.validator.Validate.Struct(reqData); validationErr != nil {
		err := errors.BadRequest.Wrap(validationErr, "Invalid input information")
		err = errors.SetCustomMessage(err, "Invalid input information")
		err = errors.AddErrorContextBlock(err, cc.validator.GenerateValidationResponse(validationErr))
		responses.HandleError(c, err)
		return
	}

	if reqData.Password != reqData.ConfirmPassword {
		cc.logger.Zap.Error("Password and confirm password not matching : ")
		err := errors.BadRequest.New("Password and confirm password should be same.")
		responses.HandleError(c, err)
		return
	}

	if _, err := cc.service.GetOneUserWithEmail(reqData.Email); err != nil {
		if err != gorm.ErrRecordNotFound {
			cc.logger.Zap.Error("Error [CreateUser] [db CreateUser]: Failed to create user")
			responses.HandleError(c, err)
			return
		}
	}

	if _, err := cc.service.GetOneUserWithPhone(reqData.Phone); err != nil {
		if err != gorm.ErrRecordNotFound {
			cc.logger.Zap.Error("Error [CreateUser] [db CreateUser]: Failed to create user")
			responses.HandleError(c, err)
			return
		}
	}

	if err := cc.service.WithTrx(trx).CreateUser(reqData.User); err != nil {
		cc.logger.Zap.Error("Error [CreateUser] [db CreateUser]: ", err.Error())
		err := errors.InternalError.Wrap(err, "Failed to create user")
		responses.HandleError(c, err)
		return
	}

	responses.SuccessJSON(c, "User Created Successfully")
}

// GetAllUsers Get All User
// @Summary				Get all User.
// @Param				page_size query string false "10"
// @Param				page query string false "Page no" "1"
// @Param				keyword query string false "search by name"
// @Param				Keyword2 query string false "search by type"
// @Description			Return all the User
// @Produce				application/json
// @Param Authorization header string true "Insert your access token" default(Bearer <Add access token here>)
// @Tags				User
// @Success 			200 {array} responses.DataCount{data=[]dtos.GetUserResponse}
// @Failure      		500 {object} responses.Error
// @Router				/users [get]
func (cc UserController) GetAllUsers(c *gin.Context) {
	pagination := helpers.BuildPagination[*user.UserPagination](c)

	users, count, err := cc.service.GetAllUsers(*pagination)
	if err != nil {
		cc.logger.Zap.Error("Error finding user records", err.Error())
		err := errors.InternalError.Wrap(err, "Failed to get users data")
		responses.HandleError(c, err)
		return
	}

	responses.JSONCount(c, http.StatusOK, users, count)
}

// GetUserProfile Returns logged-in user profile
// @Summary				Get one user by id
// @Description			Get one user by id
// @Param Authorization header string true "Insert your access token" default(Bearer <Add access token here>)
// @Produce				application/json
// @Tags				User
// @Success 			200 {array} responses.Data{data=dtos.GetUserResponse}
// @Failure      		500 {object} responses.Error
// @Router				/profile [get]
func (cc UserController) GetUserProfile(c *gin.Context) {
	userID := c.MustGet(constants.UserID).(string)
	if userID == "" {
		err := errors.BadRequest.New("Unable to get User Id")
		responses.HandleError(c, err)
		return
	}
	user, err := cc.service.GetOneUser(userID)
	if err != nil {
		cc.logger.Zap.Error("Error finding user profile", err.Error())
		err := errors.InternalError.Wrap(err, "Failed to get users profile data")
		responses.HandleError(c, err)
		return
	}

	responses.JSON(c, http.StatusOK, user)
}
