package user

import (
	"boilerplate-api/app/global/infrastructure"
	"boilerplate-api/app/models"
	"boilerplate-api/resources/utils"

	"gorm.io/gorm"
)

// Repository -> database structure
type Repository struct {
	db     infrastructure.Database
	logger infrastructure.Logger
}

// UserRepository -> creates a new User repository
func UserRepository(db infrastructure.Database, logger infrastructure.Logger) Repository {
	return Repository{
		db:     db,
		logger: logger,
	}
}

// WithTrx enables repository with transaction
func (c Repository) WithTrx(trxHandle *gorm.DB) Repository {
	if trxHandle == nil {
		c.logger.Zap.Error("Transaction Database not found in gin context. ")
		return c
	}
	c.db.DB = trxHandle
	return c
}

// Save -> User
func (c Repository) Create(User models.User) error {
	return c.db.DB.Create(&User).Error
}

// GetAllUser -> Get All users
func (c Repository) GetAllUsers(pagination utils.Pagination) ([]models.User, int64, error) {
	var users []models.User
	var totalRows int64 = 0
	queryBuilder := c.db.DB.Limit(pagination.PageSize).Offset(pagination.Offset).Order("created_at desc")
	queryBuilder = queryBuilder.Model(&models.User{})

	if pagination.Keyword != "" {
		searchQuery := "%" + pagination.Keyword + "%"
		queryBuilder.Where(c.db.DB.Where("`users`.`name` LIKE ?", searchQuery))
	}

	err := queryBuilder.
		Find(&users).
		Offset(-1).
		Limit(-1).
		Count(&totalRows).Error
	return users, totalRows, err
}
