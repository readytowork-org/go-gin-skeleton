package repository

import (
	"boilerplate-api/infrastructure"
	"boilerplate-api/models"
	"boilerplate-api/utils"
)

// FoodRepository database structure
type FoodRepository struct {
	db     infrastructure.Database
	logger infrastructure.Logger
}

// NewFoodRepository creates a new Food repository
func NewFoodRepository(db infrastructure.Database, logger infrastructure.Logger) FoodRepository {
	return FoodRepository{
		db:     db,
		logger: logger,
	}
}


// Create Food
func (c FoodRepository) Create(Food models.Food) (models.Food, error) {
	return Food, c.db.DB.Create(&Food).Error
}

// GetAllFood -> Get All foods
func (c FoodRepository) GetAllFood(pagination utils.Pagination) ([]models.Food, int64, error) {
	var foods []models.Food
	var totalRows int64 = 0
	queryBuider := c.db.DB.Model(&models.Food{}).Offset(pagination.Offset).Order(pagination.Sort)
	
	if !pagination.All{
		queryBuider=queryBuider.Limit(pagination.PageSize)
	}
	
	if pagination.Keyword != "" {
		searchQuery := "%" + pagination.Keyword + "%"
		queryBuider.Where(c.db.DB.Where("`food`.`title` LIKE ?", searchQuery))
	}

	err := queryBuider.
		Find(&foods).
		Offset(-1).
		Limit(-1).
		Count(&totalRows).Error
	return foods, totalRows, err
}

// GetOneFood -> Get One Food By Id
func (c FoodRepository) GetOneFood(ID int64) (models.Food, error) {
	Food := models.Food{}
	return Food, c.db.DB.
		Where("id = ?", ID).First(&Food).Error
}

// UpdateOneFood -> Update One Food By Id
func (c FoodRepository) UpdateOneFood(Food models.Food) error {
	return c.db.DB.Model(&models.Food{}).
		Where("id = ?", Food.ID).
		Updates(map[string]interface{}{
			"title":           Food.Title,
		}).Error
}

// DeleteOneFood -> Delete One Food By Id
func (c FoodRepository) DeleteOneFood(ID int64) error {
	return c.db.DB.
		Where("id = ?", ID).
		Delete(&models.Food{}).
		Error
}
