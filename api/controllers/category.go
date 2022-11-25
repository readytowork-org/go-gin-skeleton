package controllers

import (
	"boilerplate-api/api/responses"
	"boilerplate-api/api/services"
	"boilerplate-api/api/validators"
	"boilerplate-api/constants"
	"boilerplate-api/errors"
	"boilerplate-api/infrastructure"
	"boilerplate-api/models"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

type CategoryController struct {
	logger          infrastructure.Logger
	categoryService services.CategoryService
	bucketService   services.StorageBucketService
	validator       validators.CategoryValidator
}

func NewCategoryController(logger infrastructure.Logger, categoryService services.CategoryService, bucketService services.StorageBucketService, validator validators.CategoryValidator) CategoryController {
	return CategoryController{
		logger:          logger,
		categoryService: categoryService,
		bucketService:   bucketService,
		validator:       validator,
	}
}

func (cc CategoryController) CreateCategory(c *gin.Context) {
	role := fmt.Sprintf("%v", c.MustGet(constants.Role))
	fmt.Println("---role----", role, "params::::::", c.Query("key"))
	if role != constants.RoleUser {
		err := errors.Unauthorized.New("Unauthorised user")
		err = errors.SetCustomMessage(err, "Unauthorised user")
		responses.HandleError(c, err)
		return
	}
	responses.SuccessJSON(c, http.StatusOK, "Authorized user")
	category := models.Category{}
	if err := c.ShouldBindJSON(&category); err != nil {
		cc.logger.Zap.Error("Error [CreateCategory] (ShouldBindJson) : ", err)
		err := errors.BadRequest.Wrap(err, "Failed to bind category data")
		responses.HandleError(c, err)
		return
	}
	if validationErr := cc.validator.Validate.Struct(category); validationErr != nil {
		err := errors.BadRequest.Wrap(validationErr, "Validation error")
		err = errors.SetCustomMessage(err, "Invalid input information")
		err = errors.AddErrorContextBlock(err, cc.validator.GenerateValidationResponse(validationErr))
		responses.HandleError(c, err)
		return
	}
	created_category, err := cc.categoryService.CreateCategory(category)
	if err != nil {
		cc.logger.Zap.Error("Error [CreateCategory] (ShouldBindJson) : ", err)
		err := errors.BadRequest.Wrap(err, "Failed to bind category data")
		responses.HandleError(c, err)
		return
	}
	responses.SuccessJSON(c, http.StatusCreated, created_category)
	return
}