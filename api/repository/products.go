package repository

import (
	"boilerplate-api/infrastructure"
	"boilerplate-api/models"
	"boilerplate-api/utils"
)

type ProductRepository struct {
	db     infrastructure.Database
	logger infrastructure.Logger
}

func NewProductRepository(db infrastructure.Database) ProductRepository {
	return ProductRepository{
		db: db,
	}
}

func (pr ProductRepository) AddProduct(product models.ProductCreateInput) error {
	productCreate := models.Product{
		ItemName:    product.ProductName,
		ReceivedQty: product.ReceivedQty,
		ReceivedBy:  product.ReceivedBy,
	}
	return pr.db.DB.Create(&productCreate).Error
}

func (pr ProductRepository) GetAllProducts(pagination utils.Pagination) ([]models.Product, int64, error) {
	// var products
	// return pr.db.DB.Find(&models.Product{})
	var products []models.Product
	var totalRows int64 = 0

	queryBuilder := pr.db.DB.Limit(pagination.PageSize).Offset(pagination.Offset)
	queryBuilder = queryBuilder.Model(&models.Product{}).Preload("ReceivedUser")
	// pr.logger.Zap.Info(pr.db.DB.Find(products))
	err := queryBuilder.Find(&products).
		Offset(-1).
		Limit(-1).
		Count(&totalRows).Error
	return products, totalRows, err
}

func (pr ProductRepository) FilterUserProducts(id int64, pagination utils.Pagination) ([]models.Product, error) {
	var products []models.Product
	queryBuilder := pr.db.DB.Limit(pagination.PageSize).Offset(pagination.Offset)
	queryBuilder = queryBuilder.Model(&models.Product{}).Preload("ReceivedUser").Where("received_by=?", id)
	err := queryBuilder.Find(&products).
		Offset(-1).Limit(-1).Error
	return products, err
}
