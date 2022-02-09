package services

import (
	"boilerplate-api/api/repository"
	"boilerplate-api/models"
	"boilerplate-api/utils"

)

// FoodService -> struct
type FoodService struct {
	repository repository.FoodRepository
}

// NewFoodService  -> creates a new Foodservice
func NewFoodService(repository repository.FoodRepository) FoodService {
	return FoodService{
		repository: repository,
	}
}

// CreateFood -> call to create the Food
func (c FoodService) CreateFood(food models.Food) (models.Food, error) {
	return c.repository.Create(food)
}

// GetAllFood -> call to create the Food
func (c FoodService) GetAllFood(pagination utils.Pagination) ([]models.Food, int64, error) {
	return c.repository.GetAllFood(pagination)
}

// GetOneFood -> Get One Food By Id
func (c FoodService) GetOneFood(ID int64) (models.Food, error) {
	return c.repository.GetOneFood(ID)
}

// UpdateOneFood -> Update One Food By Id
func (c FoodService) UpdateOneFood(food models.Food) error {
	return c.repository.UpdateOneFood(food)
}

// DeleteOneFood -> Delete One Food By Id
func (c FoodService) DeleteOneFood(ID int64) error {
	return c.repository.DeleteOneFood(ID)

}