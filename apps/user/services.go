package user

import (
	"boilerplate-api/infrastructure"

	"gorm.io/gorm"
)

type Service struct {
	logger     infrastructure.Logger
	repository Repository
}

func ServiceConstuctor(
	logger infrastructure.Logger,
	repository Repository,
) Service {
	return Service{
		logger:     logger,
		repository: repository,
	}
}

// WithTrx repository with transaction
func (c Service) WithTrx(trxHandle *gorm.DB) Service {
	c.repository = c.repository.WithTrx(trxHandle)
	return c
}

// CreateUser to create the User
func (c Service) CreateUser(user User) error {
	err := c.repository.Create(user)
	return err
}

// GetAllUsers to get all the User
func (c Service) GetAllUsers(pagination UserPagination) ([]GetUserResponse, int64, error) {
	return c.repository.GetAllUsers(pagination)
}

// GetOneUser one user
func (c Service) GetOneUser(Id string) (GetUserResponse, error) {
	return c.repository.GetOneUser(Id)
}

// GetOneUserWithEmail Get one user with email
func (c Service) GetOneUserWithEmail(Email string) (User, error) {
	return c.repository.GetOneUserWithEmail(Email)
}

// GetOneUserWithPhone Get one user with phone
func (c Service) GetOneUserWithPhone(Phone string) (User, error) {
	return c.repository.GetOneUserWithPhone(Phone)
}
