package repository

import (
	"boilerplate-api/infrastructure"
	"boilerplate-api/models"
	"boilerplate-api/utils"
)

type CategoryRepository struct {
	db     infrastructure.Database
	logger infrastructure.Logger
}

func NewCategoryRepository(db infrastructure.Database, logger infrastructure.Logger) CategoryRepository {
	return CategoryRepository{
		db:     db,
		logger: logger,
	}
}

func (c CategoryRepository) GetAllCategories(pagination utils.Pagination) ([]models.Category, int64, error) {
	var categories []models.Category
	var totalRows int64 = 0
	queryBuilder := c.db.DB.Limit(pagination.PageSize).Offset(pagination.Offset).Order("created_at desc")
	queryBuilder = queryBuilder.Model(&models.Category{})

	if pagination.Keyword != "" {
		searchQuery := "%" + pagination.Keyword + "%"
		// queryBuilder.Where(c.db.DB.Where("`categories`.`title`LIKE ?", searchQuery).Or("`users`.`email` LIKE ?", searchQuery))
		queryBuilder.Where(c.db.DB.Where("`categories`.`title`LIKE ?", searchQuery))
	}
	err := queryBuilder.Find(&categories).Offset(-1).Limit(-1).Count(&totalRows).Error
	return categories, totalRows, err
}

func (c CategoryRepository) CreateCategory(Category models.Category) (models.Category, error) {
	return Category, c.db.DB.Create(&Category).Find(&Category).Error
}

func (c CategoryRepository) DeleteCategory(Id string) error {
	return c.db.DB.Delete(&models.Category{}, Id).Error
}

func (c CategoryRepository) GetCategory(Id string) (models.Category, error) {
	category := models.Category{}
	err := c.db.DB.First(&category, Id)
	return category, err.Error
}
func (c CategoryRepository) UpdateCategory(category models.Category, Id string) (models.Category, error) {

	err := c.db.DB.Model(&category).Where("Id = ?", Id).Updates(
		map[string]interface{}{
			"title": category.Title,
		}).Find(&category).Error
	return category, err
}
