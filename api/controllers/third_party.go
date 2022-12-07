package controllers

import (
	"boilerplate-api/api/responses"
	"boilerplate-api/api/services"
	"boilerplate-api/infrastructure"
	"boilerplate-api/models"
	"net/http"

	"github.com/gin-gonic/gin"
)

// TwitterController -> struct
type ThirdPartyController struct {
	logger            infrastructure.Logger
	thirdPartyService services.ThirdPartyService
	env               infrastructure.Env
}

// NewTwitterController -> constructor
func NewThirdPartyController(
	logger infrastructure.Logger,
	thirdPartyService services.ThirdPartyService,
	env infrastructure.Env,
) ThirdPartyController {
	return ThirdPartyController{
		logger:            logger,
		thirdPartyService: thirdPartyService,
		env:               env,
	}
}

// CreateUser -> Create User
func (cc ThirdPartyController) MerchanntRegister(c *gin.Context) {
	var merchant models.MerchanntRegisterInput

	if err := c.Bind(&merchant); err != nil {
		cc.logger.Zap.Error("Bind json error", err)
		return
	}
	resp, err := cc.thirdPartyService.MerchantRegister(merchant)
	if err != nil {
		cc.logger.Zap.Error("response eerror", err)
	}
	responses.JSON(c, http.StatusOK, resp)

	responses.SuccessJSON(c, http.StatusOK, "Third pArty heath check")
}
