package services

import (
	"boilerplate-api/api/repository"
	"boilerplate-api/dtos"
	"boilerplate-api/models"
	"boilerplate-api/utils"

	"gorm.io/gorm"
)

type UserService struct {
	repository repository.UserRepository
}

// NewUserService Creates New user service
func NewUserService(repository repository.UserRepository) UserService {
	return UserService{
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
func (c UserService) GetAllUsers(pagination utils.UserPagination) ([]dtos.GetUserResponse, int64, error) {
	return c.repository.GetAllUsers(pagination)
}

// GetOneUser one user
func (c UserService) GetOneUser(Id string) (dtos.GetUserResponse, error) {
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
