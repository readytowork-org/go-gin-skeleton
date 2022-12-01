package services

import "boilerplate-api/api/repository"

type ProductService struct {
	repository repository.ProductRepository
}

func NewProductService(pro repository.ProductRepository) ProductService {
	return ProductService{
		repository: pro,
	}
}

func (ps ProductService) AddProduct() {
	return
}

func (ps ProductService) GetAllProduct() {
	return
}
