package routes

import (
	"boilerplate-api/app/global/infrastructure"
	"boilerplate-api/app/packages/utility"

	"github.com/gin-gonic/gin"
)

// UtilityRoutes -> utility routes struct
type UtilityRoutes struct {
	router            infrastructure.Router
	Logger            infrastructure.Logger
	UtilityController utility.Controller
}

// NewUtilityRoute -> returns new utility route
func NewUtilityRoutes(
	logger infrastructure.Logger,
	router infrastructure.Router,
	UtilityController utility.Controller,
) UtilityRoutes {
	return UtilityRoutes{
		Logger:            logger,
		router:            router,
		UtilityController: UtilityController,
	}
}

// Setup -> sets up route for util entities
func (u UtilityRoutes) Setup(v1 *gin.RouterGroup) {
	utils := v1.Group("/utils")
	{
		utils.POST("/file-upload", u.UtilityController.FileUploadHandler)
		utils.POST("/s3-file-upload", u.UtilityController.FileUploadS3Handler)
	}
}
