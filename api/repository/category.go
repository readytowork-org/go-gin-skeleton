package repository

import (
	"boilerplate-api/infrastructure"
	"boilerplate-api/models"
	"boilerplate-api/utils"
)

type CategoryRepository struct {
	db     infrastructure.Database
	env    infrastructure.Env
	logger infrastructure.Logger
}

func NewCategoryRepository(db infrastructure.Database, env infrastructure.Env, logger infrastructure.Logger) CategoryRepository {
	return CategoryRepository{
		db:     db,
		env:    env,
		logger: logger,
	}
}

func (c CategoryRepository) CreateCategory(category models.Category) (*models.Category, error) {
	if err := c.db.DB.Create(&category).Error; err != nil {
		return nil, err
	}
	return &category, nil
}

func (c CategoryRepository) GetAllCategory(pagination utils.Pagination) ([]models.Category, int64, error) {
	var Category []models.Category
	var totalRows int64 = 0
	return Category, totalRows, c.db.DB.Find(&Category).Error
}

func (c CategoryRepository) GetOneCategory(Id string) (*models.Category, error) {
	var Category models.Category
	return &Category, c.db.DB.Find(&Category, Id).Error
}
