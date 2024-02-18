package controllers

import (
	"boilerplate-api/api/services"
	"boilerplate-api/api/validators"
	"boilerplate-api/constants"
	"boilerplate-api/dtos"
	"boilerplate-api/errors"
	"boilerplate-api/infrastructure"
	"boilerplate-api/models"
	"boilerplate-api/responses"
	"boilerplate-api/url_query"
	"boilerplate-api/utils"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/gin-gonic/gin"

	"gorm.io/gorm"
)

type UserController struct {
	logger       infrastructure.Logger
	userService  services.UserService
	env          infrastructure.Env
	validator    validators.UserValidator
	oAuthService services.OAuthService
}

// NewUserController Creates New user controller
func NewUserController(
	logger infrastructure.Logger,
	userService services.UserService,
	env infrastructure.Env,
	validator validators.UserValidator,
	oAuthService services.OAuthService,
) UserController {
	return UserController{
		logger:       logger,
		userService:  userService,
		env:          env,
		validator:    validator,
		oAuthService: oAuthService,
	}
}

// CreateUser Create User
// @Summary				Create User
// @Description			Create User
// @Param				data body dtos.CreateUserRequestData true "Enter JSON"
// @Param 				Authorization header string true "Insert your access token" default(Bearer <Add access token here>)
// @Produce				application/json
// @Tags				User
// @Success				200 {object} responses.Success "OK"
// @Failure      		400 {object} responses.Error
// @Failure      		500 {object} responses.Error
// @Router				/users [post]
func (cc UserController) CreateUser(c *gin.Context) {
	reqData := dtos.CreateUserRequestData{}
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

	if reqData.Password != reqData.ConfirmPassword {
		cc.logger.Zap.Error("Password and confirm password not matching : ")
		err := errors.BadRequest.New("Password and confirm password should be same.")
		responses.HandleError(c, err)
		return
	}

	if _, err := cc.userService.GetOneUserWithEmail(reqData.Email); err != nil {
		cc.logger.Zap.Error("Error [CreateUser] [db CreateUser]: User with this email already exists")
		err := errors.BadRequest.New("User with this email already exists")
		responses.HandleError(c, err)
		return
	}

	if _, err := cc.userService.GetOneUserWithPhone(reqData.Phone); err != nil {
		cc.logger.Zap.Error("Error [db GetOneUserWithPhone]: User with this phone already exists")
		err := errors.BadRequest.New("User with this phone already exists")
		responses.HandleError(c, err)
		return
	}

	if err := cc.userService.WithTrx(trx).CreateUser(reqData.User); err != nil {
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
	pagination := url_query.BuildPagination[*url_query.UserPagination](c)

	users, count, err := cc.userService.GetAllUsers(*pagination)
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
	userID := fmt.Sprintf("%v", c.MustGet(constants.UserID))

	user, err := cc.userService.GetOneUser(userID)
	if err != nil {
		cc.logger.Zap.Error("Error finding user profile", err.Error())
		err := errors.InternalError.Wrap(err, "Failed to get users profile data")
		responses.HandleError(c, err)
		return
	}

	responses.JSON(c, http.StatusOK, user)
}

// OAuthSignIn Redirects to sign in page
// @Summary				redirect to sign in page
// @Description			redirect to sign in page
// @Produce				text/html
// @Tags				User
// @Success 			200 {array} text/html
// @Failure      		500 {object} responses.Error
// @Router				/oauth/sign-in [post]
func (cc UserController) OAuthSignIn(c *gin.Context) {
	randomString := utils.GenerateRandomCode(8)
	url := cc.oAuthService.GetURL(randomString)
	fmt.Println(url)
	c.Redirect(http.StatusTemporaryRedirect, url)

}

/*
This call back route should be added to CallbackURL in G-cloud OAuthClient Services
This will enable SSO, call this API from any Application, Or you can setup similar functions in another application using same G-Client-key and certificates.
*/
func (cc UserController) OAuthCallback(c *gin.Context) {
	trx := c.MustGet(constants.DBTransaction).(*gorm.DB)
	resData := models.OAuthUser{}

	code := c.Request.FormValue("code")

	token, err := cc.oAuthService.GetToken(code)

	if err != nil {
		cc.logger.Zap.Error("error getting oauth token")
	}

	response, err := http.Get("https://www.googleapis.com/oauth2/v2/userinfo?access_token=" + token.AccessToken)

	if err != nil {
		cc.logger.Zap.Error("error getting oauth response")
	}
	defer response.Body.Close()
	data, err := ioutil.ReadAll(response.Body)
	if err != nil {
		cc.logger.Zap.Error("error reading oauth response body")
	}

	if err := json.Unmarshal(data, &resData); err != nil {
		cc.logger.Zap.Error("Error [OAuthSignUp] (ShouldBindJson) : ", err)
		err := errors.BadRequest.Wrap(err, "Failed to bind oauth user data")
		responses.HandleError(c, err)
		return
	}

	// TODO:
	// Perform the logic as per requirement. For now, we are adding the new user to database if the user doesn't exist
	userInfo, checkUserErr := cc.userService.GetOneUserWithEmail(resData.Email)

	// Add User if User not found
	if checkUserErr == gorm.ErrRecordNotFound {
		// Use id from oauth as Password and ask user to update it later or make changes in DB to set empty password and add a password later
		user := models.User{}
		user.Email = resData.Email
		user.Password = resData.OAuthId
		user.FullName = resData.Name
		// Use this expiry time in middleware and redirect to sign in page again if expired.
		user.Token = &token.AccessToken
		user.TokenExpiryTime = &token.Expiry

		if err := cc.userService.WithTrx(trx).CreateUser(user); err != nil {
			cc.logger.Zap.Error("Error [CreateUser] [db CreateUser]: ", err.Error())
			err := errors.InternalError.Wrap(err, "Failed to create user")
			responses.HandleError(c, err)
			return
		}

		responses.JSON(c, http.StatusOK, user.Token)
		return
	}

	if checkUserErr != nil {
		cc.logger.Zap.Error("Error finding user profile", checkUserErr.Error())
		err := errors.InternalError.Wrap(checkUserErr, "Failed to get users profile data")
		responses.HandleError(c, err)
		return
	}

	// Update the token exipry time on re signin
	userInfo.TokenExpiryTime = &token.Expiry

	if err := cc.userService.WithTrx(trx).Update(userInfo); err != nil {
		cc.logger.Zap.Error("Error [UpdateUser] [db UpdateUser]: ", err.Error())
		err := errors.InternalError.Wrap(err, "Failed to update user")
		responses.HandleError(c, err)
		return
	}

	responses.JSON(c, http.StatusOK, userInfo.Token)

}
