package repository

import (
	"boilerplate-api/errors"
	"boilerplate-api/infrastructure"
	"boilerplate-api/models"
	"boilerplate-api/utils"
	"encoding/json"
	"fmt"
	"strings"

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
func (c UserRepository) Create(User *models.User) (*models.User, error) {
	c.logger.Zap.Info(User, "---User")
	b, err := json.MarshalIndent(User, "", "")
	if err == nil {
		fmt.Println(string(b))
	}
	if err := c.db.DB.Create(&User).Error; err != nil {

		return nil, err
	}
	return User, nil
}

// GetAllUser -> Get All users
func (c UserRepository) GetAllUsers(pagination utils.Pagination) ([]models.User, int64, error) {
	var users []models.User
	var totalRows int64 = 0
	queryBuilder := c.db.DB.Limit(pagination.PageSize).Offset(pagination.Offset).Order("created_at desc")
	queryBuilder = queryBuilder.Model(&models.User{})

	if pagination.Keyword != "" {
		searchQuery := "%" + pagination.Keyword + "%"
		queryBuilder.Where(c.db.DB.Where("`user`.`username` LIKE ?", searchQuery).
			Or("`user`.`username` LIKE ?", searchQuery).
			Or("`user`.`email` LIKE ?", searchQuery).
			Or("`user`.`full_name` LIKE ?", searchQuery))
	}

	err := queryBuilder.
		Find(&users).
		Offset(-1).
		Limit(-1).
		Count(&totalRows).Error
	return users, totalRows, err
}

// Partial update of user
func (c UserRepository) UpdatePartial(ID int64, map_update map[string]interface{}) (*models.User, error) {
	user := models.User{}
	if err := c.db.DB.Model(&models.User{}).
		Where("id = ?", ID).
		Updates(map_update).Find(&user).Error; err != nil {
		if strings.Contains(err.Error(), "1062") {
			err = errors.BadRequest.Wrap(err, "Error updating user")
			custom_msg := ""
			if strings.Contains(err.Error(), "UQ_user_email") {
				custom_msg = "Email address already taken"
			} else if strings.Contains(err.Error(), "users.UQ_user_phone") {
				custom_msg = "Phone number already taken"
			}
			err = errors.SetCustomMessage(err, custom_msg)
		} else {
			err = errors.InternalError.Wrap(err, "Error updating user")
		}
		return nil, err
	}
	return &user, nil
}

func (c UserRepository) GetOneUser(Id string) (*models.User, error) {
	user := models.User{}
	if err := c.db.DB.First(&user, Id).Error; err != nil {
		return nil, err
	}
	return &user, nil

}

func (c UserRepository) DeleteOneUser(Id string) (*string, error) {
	user := models.User{}
	err := c.db.DB.First(&user, Id).Delete(&user, Id).Error
	return &user.FirebaseUID, err
}
