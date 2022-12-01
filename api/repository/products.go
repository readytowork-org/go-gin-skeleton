package repository

import "boilerplate-api/infrastructure"

type ProductRepository struct {
	db infrastructure.Database
}

func NewProductRepository(db infrastructure.Database) ProductRepository {
	return ProductRepository{
		db: db,
	}
}

func (pr ProductRepository) AddProduct() {

}

func (pr ProductRepository) GetAllProducts() {

}
