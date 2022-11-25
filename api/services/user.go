package services

import (
	"boilerplate-api/api/repository"
	"boilerplate-api/models"
	"boilerplate-api/utils"

	"gorm.io/gorm"
)

// UserService -> struct
type UserService struct {
	repository repository.UserRepository
}

// NewUserService -> creates a new Userservice
func NewUserService(repository repository.UserRepository) UserService {
	return UserService{
		repository: repository,
	}
}

// WithTrx -> enables repository with transaction
func (c UserService) WithTrx(trxHandle *gorm.DB) UserService {
	c.repository = c.repository.WithTrx(trxHandle)
	return c
}

// CreateUser -> call to create the User
func (c UserService) CreateUser(user *models.User) (*models.User, error) {
	return c.repository.Create(user)
}

// GetAllUser -> call to get all the User
func (c UserService) GetAllUsers(pagination utils.Pagination) ([]models.User, int64, error) {
	return c.repository.GetAllUsers(pagination)
}

// upate user partially
func (c UserService) UpdatePartial(ID int64, map_update map[string]interface{}) (*models.User, error) {
	return c.repository.UpdatePartial(ID, map_update)
}

func (c UserService) UpdateUser(ID string, map_update map[string]interface{}) (*models.User, error) {
	return c.repository.UpdateUser(ID, map_update)
}

func (c UserService) GetOneUser(Id string) (*models.User, error) {
	return c.repository.GetOneUser(Id)
}
func (c UserService) DeleteOneUser(Id string) (*string, error) {
	return c.repository.DeleteOneUser(Id)
}
