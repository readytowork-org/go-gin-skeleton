package services

import (
	"boilerplate-api/api/repository"
	"boilerplate-api/infrastructure"
	"boilerplate-api/models"
	"boilerplate-api/utils"
)

type UserProfileService struct {
	logger                infrastructure.Logger
	userProfileRepository repository.UserProfileRepository
}

func NewUserProfileService(logger infrastructure.Logger, userProfileRepository repository.UserProfileRepository) UserProfileService {
	return UserProfileService{
		logger:                logger,
		userProfileRepository: userProfileRepository,
	}
}

func (c UserProfileService) GetAllUserProfile(pagination utils.Pagination) ([]models.UserProfile, int64, error) {
	return c.userProfileRepository.GetAllUserProfile(pagination)
}

func (c UserProfileService) CreateUserProfile(userProfile models.UserProfile) (*models.UserProfile, error) {
	return c.userProfileRepository.CreateUserProfile(userProfile)
}

func (c UserProfileService) GetUserProfile(Id string) (models.UserProfileDetail, error) {
	return c.userProfileRepository.GetUserProfile(Id)
}
