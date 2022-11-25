package repository

import (
	"boilerplate-api/api/validators"
	"boilerplate-api/infrastructure"
	"boilerplate-api/models"
	"boilerplate-api/utils"
)

type BlogRepository struct {
	logger    infrastructure.Logger
	db        infrastructure.Database
	env       infrastructure.Env
	validator validators.BlogValidator
}

func NewBlogRepository(
	logger infrastructure.Logger,
	db infrastructure.Database,
	env infrastructure.Env,
	validator validators.BlogValidator,
) BlogRepository {
	return BlogRepository{
		logger:    logger,
		db:        db,
		env:       env,
		validator: validator,
	}
}

func (c BlogRepository) GetAllBlogs(pagination utils.Pagination) ([]models.Blog, int64, error) {
	var blogs []models.Blog
	var totalRows int64 = 0
	err := c.db.DB.Model(&models.Blog{}).
		Count(&totalRows).Find(&blogs).Error

	if err != nil {
		c.logger.Zap.Info("---error---", err)
	}

	return blogs, totalRows, err

}
