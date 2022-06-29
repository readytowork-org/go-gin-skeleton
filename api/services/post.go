package services

import (
	"boilerplate-api/api/repository"
	"boilerplate-api/models"
	"boilerplate-api/utils"

)

// PostService -> struct
type PostService struct {
	repository repository.PostRepository
}

// NewPostService  -> creates a new Postservice
func NewPostService(repository repository.PostRepository) PostService {
	return PostService{
		repository: repository,
	}
}

// CreatePost -> call to create the Post
func (c PostService) CreatePost(post models.Post) (models.Post, error) {
	return c.repository.Create(post)
}

// GetAllPost -> call to create the Post
func (c PostService) GetAllPost(pagination utils.Pagination) ([]models.Post, int64, error) {
	return c.repository.GetAllPost(pagination)
}

// GetOnePost -> Get One Post By Id
func (c PostService) GetOnePost(ID int64) (models.Post, error) {
	return c.repository.GetOnePost(ID)
}

// UpdateOnePost -> Update One Post By Id
func (c PostService) UpdateOnePost(post models.Post) error {
	return c.repository.UpdateOnePost(post)
}

// DeleteOnePost -> Delete One Post By Id
func (c PostService) DeleteOnePost(ID int64) error {
	return c.repository.DeleteOnePost(ID)

}