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

func (c ProductRepository) Create(Product models.Products) error {
	return c.db.DB.Create(&Product).Error
}

func (c ProductRepository) GetAllProducts(pagination utils.Pagination) ([]models.Products, int64, error) {
	var products []models.Products
	var totalRows int64 = 0
	queryBuilder := c.db.DB.Limit(pagination.PageSize).Offset(pagination.Offset).Order("created_at")
	queryBuilder = queryBuilder.Model(&models.Products{})

	if pagination.Keyword != "" {
		searchQuery := "%" + pagination.Keyword + "%"
		queryBuilder.Where(c.db.DB.Where("`products`.`title` LIKE ?", searchQuery))
	}

	err := queryBuilder.
		Find(&products).
		Offset(-1).
		Limit(-1).Count(&totalRows).Error
	return products, totalRows, err
}

func (c ProductRepository) GetProduct(product models.Products) (models.Products, error) {
	err := c.db.DB.Model(&product).First(&product).Error
	return product, err
}
