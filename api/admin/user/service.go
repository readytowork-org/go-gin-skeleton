package user

import (
	"boilerplate-api/api/user/user"
	"errors"

	"gorm.io/gorm"
)

var (
	ErrPasswordMismatch   = errors.New("password and confirm password do not match")
	ErrEmailAlreadyExists = errors.New("user with this email already exists")
	ErrPhoneAlreadyExists = errors.New("user with this phone already exists")
)

type Service struct {
	repository UserRepository
}

// NewService Creates New user service
func NewService(repository UserRepository) Service {
	return Service{
		repository: repository,
	}
}

// WithTrx repository with transaction
func (c Service) WithTrx(trxHandle *gorm.DB) Service {
	c.repository = c.repository.WithTrx(trxHandle)
	return c
}

// CreateUser validates the request and creates a new user. Duplicate email or phone and password mismatches are reported via sentinel errors so callers can distinguish them with errors.Is.
func (c Service) CreateUser(req CreateUserRequestData) error {
	if req.Password != req.ConfirmPassword {
		return ErrPasswordMismatch
	}

	if _, err := c.repository.GetOneUserWithEmail(req.Email); err == nil {
		return ErrEmailAlreadyExists
	} else if !errors.Is(err, gorm.ErrRecordNotFound) {
		return err
	}

	if _, err := c.repository.GetOneUserWithPhone(req.Phone); err == nil {
		return ErrPhoneAlreadyExists
	} else if !errors.Is(err, gorm.ErrRecordNotFound) {
		return err
	}

	return c.repository.Create(req.CUser)
}

// GetAllUsers to get all the CreateUser
func (c Service) GetAllUsers(pagination Pagination) ([]GetUserResponse, int64, error) {
	return c.repository.GetAllUsers(pagination)
}

// GetOneUser one user
func (c Service) GetOneUser(Id int64) (GetUserResponse, error) {
	return c.repository.GetOneUser(Id)
}

// GetOneUserWithEmail Get one user with email
func (c Service) GetOneUserWithEmail(Email string) (user.CUser, error) {
	return c.repository.GetOneUserWithEmail(Email)
}

// GetOneUserWithPhone Get one user with phone
func (c Service) GetOneUserWithPhone(Phone string) (user.CUser, error) {
	return c.repository.GetOneUserWithPhone(Phone)
}
