package repository

import (
	"boilerplate-api/infrastructure"
	"boilerplate-api/models"

	"gorm.io/gorm"
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

func (pr ProductRepository) GetAllProducts() ([]models.Product, error) {
	// var products
	// return pr.db.DB.Find(&models.Product{})
	var products []models.Product
	// pr.logger.Zap.Info(pr.db.DB.Find(products))
	return products, pr.db.DB.Model(&models.Product{}).Preload("ReceivedUser").Find(&products).Error
}

func (pr ProductRepository) FilterUserProducts(id int64) *gorm.DB {
	var products []models.Product
	return pr.db.DB.Where("received_by=?", id).Find(&products)
}
