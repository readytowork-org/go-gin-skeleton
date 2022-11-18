package services

import (
	"boilerplate-api/api/repository"
	"boilerplate-api/models"
	"boilerplate-api/utils"
)

type ProductService struct {
	repository repository.ProductRepository
}

func NewProductService(repository repository.ProductRepository) ProductService {
	return ProductService{
		repository: repository,
	}
}

func (c ProductService) CreateProduct(product models.Products) error {
	err := c.repository.Create(product)
	return err
}

func (c ProductService) GetAllProducts(pagination utils.Pagination) ([]models.Products, int64, error) {
	return c.repository.GetAllProducts(pagination)
}

func (c ProductService) GetProduct(product models.Products) (models.Products, error) {
	return c.repository.GetProduct(product)
}
