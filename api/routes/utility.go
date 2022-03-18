package routes

import (
	"boilerplate-api/api/controllers"
	"boilerplate-api/infrastructure"
)

// UtilityRoutes -> utility routes struct
type UtilityRoutes struct {
	router            infrastructure.Router
	Logger            infrastructure.Logger
	UtilityController controllers.UtilityController
}

//NewUtilityRoute -> returns new utility route
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

//Setup -> sets up route for util entities
func (u UtilityRoutes) Setup() {
	utils := u.router.Gin.Group("/utils")
	{
		utils.POST("/file-upload", u.UtilityController.FileUploadHandler)
		utils.POST("/s3-file-upload", u.UtilityController.FileUploadS3Handler)
	}
}
