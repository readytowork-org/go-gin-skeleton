package routes

import (
	"boilerplate-api/api/controllers"
	"boilerplate-api/infrastructure"
)

type UserProfileRoutes struct {
	logger                infrastructure.Logger
	router                infrastructure.Router
	userProfileController controllers.UserProfileController
}

func NewUserProfileRoutes(logger infrastructure.Logger, router infrastructure.Router, userProfileController controllers.UserProfileController) UserProfileRoutes {
	return UserProfileRoutes{
		logger:                logger,
		router:                router,
		userProfileController: userProfileController,
	}
}

func (i UserProfileRoutes) Setup() {
	i.logger.Zap.Info("settting up category routes")
	userProfiles := i.router.Gin.Group("/user-profile")
	{
		userProfiles.GET("", i.userProfileController.GetAllUserProfile)
		userProfiles.GET("/:id", i.userProfileController.GetUserProfile)
		userProfiles.POST("", i.userProfileController.CreateUserProfile)
	}
}
