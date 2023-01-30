package repository

import (
	"boilerplate-api/infrastructure"
	"boilerplate-api/models"
	"boilerplate-api/utils"

	"gorm.io/gorm"
)

// UserRepository -> database structure
type UserRepository struct {
	db     infrastructure.Database
	logger infrastructure.Logger
}

// NewUserRepository -> creates a new User repository
func NewUserRepository(db infrastructure.Database, logger infrastructure.Logger) UserRepository {
	return UserRepository{
		db:     db,
		logger: logger,
	}
}

// WithTrx enables repository with transaction
func (c UserRepository) WithTrx(trxHandle *gorm.DB) UserRepository {
	if trxHandle == nil {
		c.logger.Zap.Error("Transaction Database not found in gin context. ")
		return c
	}
	c.db.DB = trxHandle
	return c
}

// Save -> User
func (c UserRepository) Create(User models.User) error {
	return c.db.DB.Create(&User).Error
}

// GetAllUser -> Get All users
func (c UserRepository) GetAllUsers(pagination utils.Pagination) ([]models.User, int64, error) {
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

func (c UserRepository) GetOneUser(Id string) (*models.User, error) {
	user := models.User{}
	err := c.db.DB.Model(&user).Where("id = ?", Id).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil

}

func (c UserRepository) GetOneUserWithEmail(Email string) (*models.User, error) {
	user := models.User{}
	err := c.db.DB.Model(&user).Where("email = ?", Email).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil

}

func (c UserRepository) GetOneUserWithPhone(Phone string) (*models.User, error) {
	user := models.User{}
	if err := c.db.DB.First(&user, "phone = ?", Phone).Error; err != nil {
		return nil, err
	}
	return &user, nil

}
