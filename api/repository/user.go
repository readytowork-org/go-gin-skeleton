package repository

import (
	"boilerplate-api/api/user"
	"boilerplate-api/database/models"
	"boilerplate-api/dtos"
	"boilerplate-api/internal/config"
	"gorm.io/gorm"
)

// UserRepository database structure
type UserRepository struct {
	db     config.Database
	logger config.Logger
}

// NewUserRepository creates a new User repository
func NewUserRepository(db config.Database, logger config.Logger) UserRepository {
	return UserRepository{
		db:     db,
		logger: logger,
	}
}

// WithTrx enables repository with transaction
func (c UserRepository) WithTrx(trxHandle *gorm.DB) UserRepository {
	if trxHandle == nil {
		c.logger.Error("Transaction Database not found in gin context. ")
		return c
	}
	c.db.DB = trxHandle
	return c
}

// Create user
func (c UserRepository) Create(User models.User) error {
	return c.db.DB.Create(&User).Error
}

// GetAllUsers Get All users
func (c UserRepository) GetAllUsers(pagination user.Pagination) (users []dtos.GetUserResponse, count int64, err error) {
	queryBuilder := c.db.DB.Limit(pagination.PageSize).Offset(pagination.Offset).Order("created_at desc")
	queryBuilder = queryBuilder.Model(&models.User{})

	if pagination.Keyword != "" {
		searchQuery := "%" + pagination.Keyword + "%"
		queryBuilder.Where(c.db.DB.Where("`users`.`name` LIKE ?", searchQuery))
	}

	return users, count, queryBuilder.
		Find(&users).
		Offset(-1).
		Limit(-1).
		Count(&count).
		Error
}

func (c UserRepository) GetOneUser(Id string) (userModel dtos.GetUserResponse, err error) {
	return userModel, c.db.DB.
		Model(&userModel).
		Where("id = ?", Id).
		First(&userModel).
		Error
}

func (c UserRepository) GetOneUserWithEmail(Email string) (user models.User, err error) {
	return user, c.db.DB.Model(&user).
		Where("email = ?", Email).
		First(&user).
		Error
}

func (c UserRepository) GetOneUserWithPhone(Phone string) (user models.User, err error) {
	return user, c.db.DB.
		First(&user, "phone = ?", Phone).
		Error

}
