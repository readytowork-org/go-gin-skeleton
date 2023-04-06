package repository

import (
	"boilerplate-api/errors"
	"boilerplate-api/infrastructure"
	"boilerplate-api/models"
	"boilerplate-api/utils"
	"fmt"
	"strings"
)

type UserProfileRepository struct {
	logger infrastructure.Logger
	db     infrastructure.Database
	env    infrastructure.Env
}

func NewUserProfileRepository(logger infrastructure.Logger, db infrastructure.Database, env infrastructure.Env) UserProfileRepository {
	return UserProfileRepository{
		logger: logger,
		db:     db,
		env:    env,
	}
}

func (c UserProfileRepository) GetAllUserProfile(pagination utils.Pagination) ([]models.UserProfile, int64, error) {

	var userProfiles []models.UserProfile
	var totalRows int64 = 0
	queryBuilder := c.db.DB.Limit(pagination.PageSize).Offset(pagination.Offset).Order("created_at desc")
	queryBuilder = queryBuilder.Model(&models.UserProfile{})
	qr := fmt.Sprintf("SELECT u.id, u.email, u.full_name, up.address, up.contact, up.id FROM user_profile AS up INNER JOIN users AS u ON up.user_id = u.id;")
	data := c.db.DB.Raw(qr).Scan(
		&models.UserProfile{},
	)

	// if pagination.Keyword != "" {
	// 	rawQuery := fmt.Sprintf("selct u.id, u.email, u.full_name, up.address, up.contact, up.id FROM user_profile AS up INNER JOIN user AS u ON up.user_id = u.id;")
	// 	searchQuery := "%" + pagination.Keyword + "%"
	// 	queryBuilder.Where(c.db.DB.Raw(rawQuery))
	// }
	err := data.Find(&userProfiles).Offset(-1).Limit(10).Count(&totalRows).Error
	return userProfiles, totalRows, err
}

func (c UserProfileRepository) CreateUserProfile(UserProfile models.UserProfile) (*models.UserProfile, error) {
	if err := c.db.DB.Create(&UserProfile).Error; err != nil {
		if strings.Contains(err.Error(), "1062") {
			c.logger.Zap.Info("iside error string----------")
			err = errors.BadRequest.Wrap(err, "Error creating user")
			custom_msg := ""
			if strings.Contains(err.Error(), "UQ_user_profile_user_id") {
				c.logger.Zap.Info("iside user id string----------")
				custom_msg = "User Profile already exists."
			} else if strings.Contains(err.Error(), "UQ_user_profile_contact") {
				c.logger.Zap.Info("iside phone number string----------")
				custom_msg = "Phone number already taken"
			}
			err = errors.SetCustomMessage(err, custom_msg)
		} else {
			c.logger.Zap.Info("iside else string----------")
			err = errors.InternalError.Wrap(err, "Error updating user")
		}

		return nil, err

	}
	return &UserProfile, nil
}

func (c UserProfileRepository) GetUserProfile(Id string) (models.UserProfileDetail, error) {
	userProfile := models.UserProfile{}
	user := models.User{}
	err := c.db.DB.First(&userProfile, Id).Error
	err2 := c.db.DB.First(&user, Id).Error
	fmt.Println(err2)
	var userProfileDetail models.UserProfileDetail
	userProfileDetail.UserDetail = user
	userProfileDetail.UserProfile = userProfile

	return userProfileDetail, err

}
