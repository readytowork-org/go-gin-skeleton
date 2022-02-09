package controllers

import (
	"net/http"
	"strconv"
	"boilerplate-api/api/responses"
	"boilerplate-api/api/services"
	"boilerplate-api/infrastructure"
	"boilerplate-api/models"
	"boilerplate-api/utils"

	"github.com/gin-gonic/gin"
)

// FoodController -> struct
type FoodController struct {
	logger                 infrastructure.Logger
	FoodService  services.FoodService
}

// NewFoodController -> constructor
func NewFoodController(
	logger infrastructure.Logger,
	FoodService services.FoodService,
) FoodController {
	return FoodController{
		logger:                  logger,
		FoodService:  FoodService,
	}
}

// CreateFood -> Create Food
func (cc FoodController) CreateFood(c *gin.Context) {
	food := models.Food{}

	if err := c.ShouldBindJSON(&food); err != nil {
		cc.logger.Zap.Error("Error [CreateFood] (ShouldBindJson) : ", err)
		responses.ErrorJSON(c, http.StatusBadRequest, "Failed To Create Food")
		return
	}

	if _, err := cc.FoodService.CreateFood(food); err != nil {
		cc.logger.Zap.Error("Error [CreateFood] [db CreateFood]: ", err.Error())
		responses.ErrorJSON(c, http.StatusInternalServerError, "Failed To Create Food")
		return
	}

	responses.SuccessJSON(c, http.StatusOK, "Food Created Sucessfully")
}

// GetAllFood -> Get All Food
func (cc FoodController) GetAllFood(c *gin.Context) {

	pagination := utils.BuildPagination(c)
	pagination.Sort = "created_at desc"
	foods, count, err := cc.FoodService.GetAllFood(pagination)

	if err != nil {
		cc.logger.Zap.Error("Error finding Food records", err.Error())
		responses.ErrorJSON(c, http.StatusBadRequest, "Failed to Find Food")
		return
	}
	responses.JSONCount(c, http.StatusOK, foods, count)

}

// GetOneFood -> Get One Food
func (cc FoodController) GetOneFood(c *gin.Context) {
	ID, _ := strconv.ParseInt(c.Param("id"), 10, 64)
	food, err := cc.FoodService.GetOneFood(ID)

	if err != nil {
		cc.logger.Zap.Error("Error [GetOneFood] [db GetOneFood]: ", err.Error())
		responses.ErrorJSON(c, http.StatusBadRequest, "Failed to Find Food")
		return
	}
	responses.JSON(c, http.StatusOK, food)

}

// UpdateOneFood -> Update One Food By Id
func (cc FoodController) UpdateOneFood(c *gin.Context) {
	ID, _ := strconv.ParseInt(c.Param("id"), 10, 64)
	food := models.Food{}

	if err := c.ShouldBindJSON(&food); err != nil {
		cc.logger.Zap.Error("Error [UpdateFood] (ShouldBindJson) : ", err)
		responses.ErrorJSON(c, http.StatusBadRequest, "failed to update food")
		return
	}
	food.ID = ID

	if err := cc.FoodService.UpdateOneFood(food); err != nil {
		cc.logger.Zap.Error("Error [UpdateFood] [db UpdateFood]: ", err.Error())
		responses.ErrorJSON(c, http.StatusInternalServerError, "failed to update food")
		return
	}

	responses.SuccessJSON(c, http.StatusOK, "Food Updated Sucessfully")
}

// DeleteOneFood -> Delete One Food By Id
func (cc FoodController) DeleteOneFood(c *gin.Context) {
	ID, _ := strconv.ParseInt(c.Param("id"), 10, 64)
	err := cc.FoodService.DeleteOneFood(ID)

	if err != nil {
		cc.logger.Zap.Error("Error [DeleteOneFood] [db DeleteOneFood]: ", err.Error())
		responses.ErrorJSON(c, http.StatusBadRequest, "Failed to Delete Food")
		return
	}

	responses.SuccessJSON(c, http.StatusOK, "Food Deleted Sucessfully")
}
