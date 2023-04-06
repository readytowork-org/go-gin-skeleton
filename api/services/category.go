package services

import (
	"boilerplate-api/api/repository"
	"boilerplate-api/models"
	"boilerplate-api/utils"
)

type CategoryService struct {
	repository repository.CategoryRepository
}

func NewCategoryService(repository repository.CategoryRepository) CategoryService {
	return CategoryService{
		repository: repository,
	}
}

func (c CategoryService) GetAllCategories(pagination utils.Pagination) ([]models.Category, int64, error) {
	return c.repository.GetAllCategories(pagination)
}

func (c CategoryService) CreateCategory(category models.Category) (models.Category, error) {
	return c.repository.CreateCategory(category)
}

func (c CategoryService) DeleteCategory(categoryId string) error {

	return c.repository.DeleteCategory(categoryId)
}

func (c CategoryService) GetCategory(categoryId string) (models.Category, error) {
	return c.repository.GetCategory(categoryId)
}
func (c CategoryService) UpdateCategory(category models.Category, categoryId string) (models.Category, error) {
	return c.repository.UpdateCategory(category, categoryId)
}
