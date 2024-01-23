package repository

import (
	"boilerplate-api/apps/user"
	"boilerplate-api/apps/user/models"
	"boilerplate-api/infrastructure"

	"gorm.io/gorm"
)

type UserRepository struct {
	logger infrastructure.Logger
	db     infrastructure.Database
}

func UserRepositoryConstuctor(
	logger infrastructure.Logger,
	db infrastructure.Database,
) UserRepository {
	return UserRepository{
		logger: logger,
		db:     db,
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

// Create user
func (c UserRepository) Create(obj models.User) error {
	return c.db.DB.Create(&obj).Error
}

// GetAllUsers Get All users
func (c UserRepository) GetAllUsers(pagination user.UserPagination) (users []user.GetUserResponse, count int64, err error) {
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

func (c UserRepository) GetOneUser(Id string) (userModel user.GetUserResponse, err error) {
	return userModel, c.db.DB.
		Model(&models.User{}).
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
