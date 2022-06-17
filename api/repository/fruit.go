package repository

import (
	"boilerplate-api/infrastructure"
	"boilerplate-api/models"
	"boilerplate-api/utils"
)

// FruitRepository database structure
type FruitRepository struct {
	db     infrastructure.Database
	logger infrastructure.Logger
}

// NewFruitRepository creates a new Fruit repository
func NewFruitRepository(db infrastructure.Database, logger infrastructure.Logger) FruitRepository {
	return FruitRepository{
		db:     db,
		logger: logger,
	}
}


// Create Fruit
func (c FruitRepository) Create(Fruit models.Fruit) (models.Fruit, error) {
	return Fruit, c.db.DB.Create(&Fruit).Error
}

// GetAllFruit -> Get All fruits
func (c FruitRepository) GetAllFruit(pagination utils.Pagination) ([]models.Fruit, int64, error) {
	var fruits []models.Fruit
	var totalRows int64 = 0
	queryBuider := c.db.DB.Model(&models.Fruit{}).Offset(pagination.Offset).Order(pagination.Sort)
	
	if !pagination.All{
		queryBuider=queryBuider.Limit(pagination.PageSize)
	}
	
	if pagination.Keyword != "" {
		searchQuery := "%" + pagination.Keyword + "%"
		queryBuider.Where(c.db.DB.Where("`fruit`.`title` LIKE ?", searchQuery))
	}

	err := queryBuider.
		Find(&fruits).
		Offset(-1).
		Limit(-1).
		Count(&totalRows).Error
	return fruits, totalRows, err
}

// GetOneFruit -> Get One Fruit By Id
func (c FruitRepository) GetOneFruit(ID int64) (models.Fruit, error) {
	Fruit := models.Fruit{}
	return Fruit, c.db.DB.
		Where("id = ?", ID).First(&Fruit).Error
}

// UpdateOneFruit -> Update One Fruit By Id
func (c FruitRepository) UpdateOneFruit(Fruit models.Fruit) error {
	return c.db.DB.Model(&models.Fruit{}).
		Where("id = ?", Fruit.ID).
		Updates(map[string]interface{}{
			"name":           Fruit.Name,
			"season":			Fruit.Season,
			"metric":			Fruit.Metric,
			"price":			Fruit.Price,
		}).Error
}

// DeleteOneFruit -> Delete One Fruit By Id
func (c FruitRepository) DeleteOneFruit(ID int64) error {
	return c.db.DB.
		Where("id = ?", ID).
		Delete(&models.Fruit{}).
		Error
}
