package repository

import (
	"boilerplate-api/infrastructure"
	"boilerplate-api/models"
	"boilerplate-api/utils"
	"encoding/json"
	"fmt"
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

func (c CategoryRepository) GetAllCategory(pagination utils.Pagination) ([]models.Category, error) {
	var Category []models.Category
	return Category, c.db.DB.Find(&Category).Error
}

func (c CategoryRepository) GetOneCategory(Id string) (*models.Category, error) {
	var Category models.Category
	return &Category, c.db.DB.Find(&Category, Id).Error
}

func (c CategoryRepository) DeleteOneCategory(Id string) error {
	var Category models.Category
	err := c.db.DB.First(&Category, Id).Delete(&Category, Id).Error
	if err != nil {
		c.logger.Zap.Info(err, "____category err_______")
		return err

	}
	return nil
}

func (c CategoryRepository) UpdateOneCategory(Category models.Category) (*models.Category, error) {

	b, errr := json.MarshalIndent(Category, "", " ")
	if errr == nil {
		fmt.Println(string(b))
	}
	err := c.db.DB.Model(&Category).Where("id=?", Category.ID).Updates(models.Category{Title: Category.Title}).Error

	if err != nil {
		c.logger.Zap.Info(err, "____category err_______")
		return nil, err

	}
	return &Category, nil
}
