package user

import (
	"boilerplate-api/infrastructure"

	"gorm.io/gorm"
)

type Repository struct {
	logger infrastructure.Logger
	db     infrastructure.Database
}

func RepositoryConstuctor(
	logger infrastructure.Logger,
	db infrastructure.Database,
) Repository {
	return Repository{
		logger: logger,
		db:     db,
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

// Create user
func (c Repository) Create(obj User) error {
	return c.db.DB.Create(&obj).Error
}

// GetAllUsers Get All users
func (c Repository) GetAllUsers(pagination UserPagination) (users []GetUserResponse, count int64, err error) {
	queryBuilder := c.db.DB.Limit(pagination.PageSize).Offset(pagination.Offset).Order("created_at desc")
	queryBuilder = queryBuilder.Model(&User{})

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

func (c Repository) GetOneUser(Id string) (userModel GetUserResponse, err error) {
	return userModel, c.db.DB.
		Model(&User{}).
		Where("id = ?", Id).
		First(&userModel).
		Error
}

func (c Repository) GetOneUserWithEmail(Email string) (user User, err error) {
	return user, c.db.DB.Model(&user).
		Where("email = ?", Email).
		First(&user).
		Error
}

func (c Repository) GetOneUserWithPhone(Phone string) (user User, err error) {
	return user, c.db.DB.
		First(&user, "phone = ?", Phone).
		Error

}
