package user

import (
	"boilerplate-api/app/helpers"
	"boilerplate-api/app/models"

	"gorm.io/gorm"
)

// Service -> struct
type Service struct {
	repository Repository
}

// UserService -> creates a new Service
func UserService(repository Repository) Service {
	return Service{
		repository: repository,
	}
}

// WithTrx -> enables repository with transaction
func (c Service) WithTrx(trxHandle *gorm.DB) Service {
	c.repository = c.repository.WithTrx(trxHandle)
	return c
}

// CreateUser -> call to create the User
func (c Service) CreateUser(user models.User) error {
	err := c.repository.Create(user)
	return err
}

// GetAllUser -> call to get all the User
func (c Service) GetAllUsers(pagination helpers.Pagination) ([]models.User, int64, error) {
	return c.repository.GetAllUsers(pagination)
}
