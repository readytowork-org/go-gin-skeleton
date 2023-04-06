package controllers

import (
	"boilerplate-api/api/responses"
	"boilerplate-api/api/services"
	"boilerplate-api/errors"
	"boilerplate-api/infrastructure"
	"boilerplate-api/models"
	"boilerplate-api/utils"
	"net/http"

	"github.com/gin-gonic/gin"
)

type CategoryController struct {
	logger          infrastructure.Logger
	categoryService services.CategoryService
	env             infrastructure.Env
}

func NewCategoryController(services services.CategoryService, logger infrastructure.Logger, env infrastructure.Env) CategoryController {
	return CategoryController{
		categoryService: services,
		logger:          logger,
		env:             env,
	}
}

func (cc CategoryController) GetAllCategories(c *gin.Context) {
	pagination := utils.BuildPagination(c)
	categories, count, err := cc.categoryService.GetAllCategories(pagination)
	if err != nil {
		cc.logger.Zap.Error("Error Finding categories", err)
		err := errors.InternalError.Wrap(err, "Failed get categories")
		responses.HandleError(c, err)
		return
	}
	responses.JSONCount(c, http.StatusOK, categories, count)
}

func (cc CategoryController) GetCategory(c *gin.Context) {
	if c.Param("id") == "" {
		responses.JSON(c, http.StatusBadRequest, "Id required in url")
		return
	}
	category, err := cc.categoryService.GetCategory(c.Param("id"))
	if err != nil {
		cc.logger.Zap.Error("Error Finding category", err)
		err := errors.InternalError.Wrap(err, "Failed get category")
		responses.HandleError(c, err)
		return
	}
	responses.JSON(c, http.StatusOK, category)
}

func (cc CategoryController) CreateCategory(c *gin.Context) {
	category := models.Category{}
	c.ShouldBindJSON(&category)
	created, err := cc.categoryService.CreateCategory(category)
	if err != nil {
		cc.logger.Zap.Error("Error Creating categories", err)
		err := errors.InternalError.Wrap(err, "Failed to create categories")
		responses.HandleError(c, err)
		return
	}
	responses.SuccessJSON(c, http.StatusCreated, created)
}

func (cc CategoryController) DeleteCategory(c *gin.Context) {
	if c.Param("id") == "" {
		responses.JSON(c, http.StatusBadRequest, "Id required in url")
		return
	}
	if err := cc.categoryService.DeleteCategory(c.Param("id")); err != nil {
		cc.logger.Zap.Error("Failed to delete the category", err)
		err := errors.InternalError.Wrap(err, "Failed to delete the category")
		responses.HandleError(c, err)
		return
	}
	responses.JSON(c, http.StatusOK, "Category deleted successfully.")

}

func (cc CategoryController) UpdateCategory(c *gin.Context) {
	cc.logger.Zap.Info("----------- UPDATE CATEGORY ---------")
	if c.Param("id") == "" {
		responses.JSON(c, http.StatusBadRequest, "Id required in url")
		return
	}
	categoryObj := models.Category{}
	c.ShouldBindJSON(&categoryObj)
	category, err := cc.categoryService.UpdateCategory(categoryObj, c.Param("id"))

	if err != nil {
		cc.logger.Zap.Error("Failed to delete the category", err)
		err := errors.InternalError.Wrap(err, "Failed to delete the category")
		responses.HandleError(c, err)
		return
	}
	cc.logger.Zap.Info("-----category id---", category)
	if category.ID == 0 {
		cc.logger.Zap.Error("category not found")
		responses.ErrorJSON(c, http.StatusBadRequest, "Failed to update the category.")
		return
	}
	responses.JSON(c, http.StatusOK, category)

}
