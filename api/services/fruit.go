package services

import (
	"boilerplate-api/api/repository"
	"boilerplate-api/models"
	"boilerplate-api/utils"

)

// FruitService -> struct
type FruitService struct {
	repository repository.FruitRepository
}

// NewFruitService  -> creates a new Fruitservice
func NewFruitService(repository repository.FruitRepository) FruitService {
	return FruitService{
		repository: repository,
	}
}

// CreateFruit -> call to create the Fruit
func (c FruitService) CreateFruit(fruit models.Fruit) (models.Fruit, error) {
	return c.repository.Create(fruit)
}

// GetAllFruit -> call to create the Fruit
func (c FruitService) GetAllFruit(pagination utils.Pagination) ([]models.Fruit, int64, error) {
	return c.repository.GetAllFruit(pagination)
}

// GetOneFruit -> Get One Fruit By Id
func (c FruitService) GetOneFruit(ID int64) (models.Fruit, error) {
	return c.repository.GetOneFruit(ID)
}

// UpdateOneFruit -> Update One Fruit By Id
func (c FruitService) UpdateOneFruit(fruit models.Fruit) error {
	return c.repository.UpdateOneFruit(fruit)
}

// DeleteOneFruit -> Delete One Fruit By Id
func (c FruitService) DeleteOneFruit(ID int64) error {
	return c.repository.DeleteOneFruit(ID)

}