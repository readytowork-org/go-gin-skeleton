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

// FruitController -> struct
type FruitController struct {
	logger                 infrastructure.Logger
	FruitService  services.FruitService
}

// NewFruitController -> constructor
func NewFruitController(
	logger infrastructure.Logger,
	FruitService services.FruitService,
) FruitController {
	return FruitController{
		logger:                  logger,
		FruitService:  FruitService,
	}
}

// CreateFruit -> Create Fruit
func (cc FruitController) CreateFruit(c *gin.Context) {
	fruit := models.Fruit{}

	if err := c.ShouldBindJSON(&fruit); err != nil {
		cc.logger.Zap.Error("Error [CreateFruit] (ShouldBindJson) : ", err)
		responses.ErrorJSON(c, http.StatusBadRequest, "Failed To Create Fruit")
		return
	}

	if _, err := cc.FruitService.CreateFruit(fruit); err != nil {
		cc.logger.Zap.Error("Error [CreateFruit] [db CreateFruit]: ", err.Error())
		responses.ErrorJSON(c, http.StatusInternalServerError, "Failed To Create Fruit")
		return
	}

	responses.SuccessJSON(c, http.StatusOK, "Fruit Created Sucessfully")
}

// GetAllFruit -> Get All Fruit
func (cc FruitController) GetAllFruit(c *gin.Context) {

	pagination := utils.BuildPagination(c)
	pagination.Sort = "created_at desc"
	fruits, count, err := cc.FruitService.GetAllFruit(pagination)

	if err != nil {
		cc.logger.Zap.Error("Error finding Fruit records", err.Error())
		responses.ErrorJSON(c, http.StatusBadRequest, "Failed to Find Fruit")
		return
	}
	responses.JSONCount(c, http.StatusOK, fruits, count)

}

// GetOneFruit -> Get One Fruit
func (cc FruitController) GetOneFruit(c *gin.Context) {
	ID, _ := strconv.ParseInt(c.Param("id"), 10, 64)
	fruit, err := cc.FruitService.GetOneFruit(ID)

	if err != nil {
		cc.logger.Zap.Error("Error [GetOneFruit] [db GetOneFruit]: ", err.Error())
		responses.ErrorJSON(c, http.StatusBadRequest, "Failed to Find Fruit")
		return
	}
	responses.JSON(c, http.StatusOK, fruit)

}

// UpdateOneFruit -> Update One Fruit By Id
func (cc FruitController) UpdateOneFruit(c *gin.Context) {
	ID, _ := strconv.ParseInt(c.Param("id"), 10, 64)
	fruit := models.Fruit{}

	if err := c.ShouldBindJSON(&fruit); err != nil {
		cc.logger.Zap.Error("Error [UpdateFruit] (ShouldBindJson) : ", err)
		responses.ErrorJSON(c, http.StatusBadRequest, "failed to update fruit")
		return
	}
	fruit.ID = ID

	if err := cc.FruitService.UpdateOneFruit(fruit); err != nil {
		cc.logger.Zap.Error("Error [UpdateFruit] [db UpdateFruit]: ", err.Error())
		responses.ErrorJSON(c, http.StatusInternalServerError, "failed to update fruit")
		return
	}

	responses.SuccessJSON(c, http.StatusOK, "Fruit Updated Sucessfully")
}

// DeleteOneFruit -> Delete One Fruit By Id
func (cc FruitController) DeleteOneFruit(c *gin.Context) {
	ID, _ := strconv.ParseInt(c.Param("id"), 10, 64)
	err := cc.FruitService.DeleteOneFruit(ID)

	if err != nil {
		cc.logger.Zap.Error("Error [DeleteOneFruit] [db DeleteOneFruit]: ", err.Error())
		responses.ErrorJSON(c, http.StatusBadRequest, "Failed to Delete Fruit")
		return
	}

	responses.SuccessJSON(c, http.StatusOK, "Fruit Deleted Sucessfully")
}
