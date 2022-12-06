package controllers

import (
	"boilerplate-api/api/responses"
	"boilerplate-api/api/services"
	"boilerplate-api/infrastructure"
	"net/http"

	"github.com/gin-gonic/gin"
)

// TwitterController -> struct
type TwitterController struct {
	logger         infrastructure.Logger
	twitterService services.TwitterService
	env            infrastructure.Env
}

// NewTwitterController -> constructor
func NewTwitterController(
	logger infrastructure.Logger,
	twitterService services.TwitterService,
	env infrastructure.Env,
) TwitterController {
	return TwitterController{
		logger:         logger,
		twitterService: twitterService,
		env:            env,
	}
}

// CreateUser -> Create User
func (cc TwitterController) CreateTweet(c *gin.Context) {

	responses.SuccessJSON(c, http.StatusOK, "Twitter heath check")
}
