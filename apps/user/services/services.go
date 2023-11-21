package services

import (
	"boilerplate-api/apps/user"
	"boilerplate-api/apps/user/models"
	"boilerplate-api/apps/user/repository"
	"boilerplate-api/infrastructure"

	"gorm.io/gorm"
)

type UserService struct {
	logger     infrastructure.Logger
	repository repository.UserRepository
}

func UserServiceConstuctor(
	logger infrastructure.Logger,
	repository repository.UserRepository,
) UserService {
	return UserService{
		logger:     logger,
		repository: repository,
	}
}

// WithTrx repository with transaction
func (c UserService) WithTrx(trxHandle *gorm.DB) UserService {
	c.repository = c.repository.WithTrx(trxHandle)
	return c
}

// CreateUser to create the User
func (c UserService) CreateUser(user models.User) error {
	err := c.repository.Create(user)
	return err
}

// GetAllUsers to get all the User
func (c UserService) GetAllUsers(pagination user.UserPagination) ([]user.GetUserResponse, int64, error) {
	return c.repository.GetAllUsers(pagination)
}

// GetOneUser one user
func (c UserService) GetOneUser(Id string) (user.GetUserResponse, error) {
	return c.repository.GetOneUser(Id)
}

// GetOneUserWithEmail Get one user with email
func (c UserService) GetOneUserWithEmail(Email string) (models.User, error) {
	return c.repository.GetOneUserWithEmail(Email)
}

// GetOneUserWithPhone Get one user with phone
func (c UserService) GetOneUserWithPhone(Phone string) (models.User, error) {
	return c.repository.GetOneUserWithPhone(Phone)
}
