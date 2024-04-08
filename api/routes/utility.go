package routes

import (
	"boilerplate-api/api/controllers"
	"boilerplate-api/internal/config"
	"boilerplate-api/internal/router"
)

// UtilityRoutes utility routes struct
type UtilityRoutes struct {
	router            router.Router
	Logger            config.Logger
	UtilityController controllers.UtilityController
}

// NewUtilityRoutes returns new utility route
func NewUtilityRoutes(
	logger config.Logger,
	router router.Router,
	UtilityController controllers.UtilityController,
) UtilityRoutes {
	return UtilityRoutes{
		Logger:            logger,
		router:            router,
		UtilityController: UtilityController,
	}
}

// Setup sets up route for util entities
func (u UtilityRoutes) Setup() {
	utils := u.router.Group("/utils")
	{
		utils.POST("/file-upload", u.UtilityController.FileUploadHandler)
		utils.POST("/s3-file-upload", u.UtilityController.FileUploadS3Handler)
	}
}
