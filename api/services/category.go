package services

import (
	"boilerplate-api/api/repository"
	"boilerplate-api/infrastructure"
	"boilerplate-api/models"
	"boilerplate-api/utils"
)

type CategoryService struct {
	repository repository.CategoryRepository
	logger     infrastructure.Logger
}

func NewCategoryRepository(repository repository.CategoryRepository, logger infrastructure.Logger) CategoryService {
	return CategoryService{
		repository: repository,
		logger:     logger,
	}
}

func (c CategoryService) CreateCategory(catergory models.Category) (*models.Category, error) {
	return c.repository.CreateCategory(catergory)

}

func (c CategoryService) GetAllCategory(pagination utils.Pagination) ([]models.Category, error) {
	return c.repository.GetAllCategory(pagination)

}

func (c CategoryService) GetOneCategory(Id string) (*models.Category, error) {
	return c.repository.GetOneCategory(Id)

}

func (c CategoryService) DeleteOneCategory(Id string) error {
	return c.repository.DeleteOneCategory(Id)

}

func (c CategoryService) UpdateOneCategory(Category models.Category) (*models.Category, error) {
	return c.repository.UpdateOneCategory(Category)

}
