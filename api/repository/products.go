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

func NewProductRepository(db infrastructure.Database, logger infrastructure.Logger) ProductRepository {
	return ProductRepository{
		db:     db,
		logger: logger,
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

func (pr ProductRepository) FilterUserProducts(id int64, pagination utils.Pagination) ([]models.Product, models.User, error) {
	var products []models.Product
	var user models.User

	queryBuilder := pr.db.DB.Limit(pagination.PageSize).Offset(pagination.Offset)
	queryBuilder = queryBuilder.Model(&models.Product{}).Preload("ReceivedUser").Where("received_by=?", id)

	pr.db.DB.Model(&models.User{}).Where("id=?", id).First(&user)

	// pr.logger.Zap.Info("recieved user", recievedUser)
	err := queryBuilder.Find(&products).
		Offset(-1).Limit(-1).Error
	return products, user, err

}

func (pr ProductRepository) SendProduct(id int64, product models.ProductSentInput) error {

	return pr.db.DB.Model(&models.Product{}).Where("item_id=?", id).Updates(map[string]interface{}{
		"item_id":  id,
		"sent_by":  product.SentBy,
		"sent_qty": product.SentQty,
	}).Error
}
