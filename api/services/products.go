package services

import (
	"boilerplate-api/api/repository"
	"boilerplate-api/models"
	"boilerplate-api/utils"
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

func (ps ProductService) GetAllProduct(pagination utils.Pagination) ([]models.Product, int64, error) {
	return ps.repository.GetAllProducts(pagination)
}

func (ps ProductService) FilterUserProducts(id int64, pagination utils.Pagination) ([]models.Product, models.User, error) {
	return ps.repository.FilterUserProducts(id, pagination)
}

// func (ps ProductService) SendProduct(id int64, product models.ProductSentInput) error {
// 	return ps.repository.SendProduct(id, product)
// }
