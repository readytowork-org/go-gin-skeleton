package routes

import (
	"boilerplate-api/app/http/controllers/v1"
	"boilerplate-api/app/infrastructure"

	"github.com/gin-gonic/gin"
)

// UtilityRoutes -> utility routes struct
type UtilityRoutes struct {
	router            infrastructure.Router
	Logger            infrastructure.Logger
	UtilityController controllers.UtilityController
}

// NewUtilityRoute -> returns new utility route
func NewUtilityRoutes(
	logger infrastructure.Logger,
	router infrastructure.Router,
	UtilityController controllers.UtilityController,
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
