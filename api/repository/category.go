package repository

import (
	"boilerplate-api/infrastructure"
	"boilerplate-api/models"
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
