package services

import (
	"boilerplate-api/api/repository"
	"boilerplate-api/models"

	"gorm.io/gorm"
)

type ProductService struct {
	repository repository.ProductRepository
}

func NewProductService(pro repository.ProductRepository) ProductService {
	return ProductService{
		repository: pro,
	}
}

func (ps ProductService) AddProduct(product models.ProductCreateInput) error {
	return ps.repository.AddProduct(product)
}

func (ps ProductService) GetAllProduct() ([]models.Product, error) {
	return ps.repository.GetAllProducts()
}

func (ps ProductService) FilterUserProducts(id int64) *gorm.DB {
	return ps.repository.FilterUserProducts(id)
}
