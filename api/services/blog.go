package services

import (
	"boilerplate-api/api/repository"
	"boilerplate-api/infrastructure"
	"boilerplate-api/models"
	"boilerplate-api/utils"
)

type BlogService struct {
	repository repository.BlogRepository
	logger     infrastructure.Logger
}

func NewBlogService(
	repository repository.BlogRepository,
	logger infrastructure.Logger,
) BlogService {
	return BlogService{
		repository: repository,
		logger:     logger,
	}
}

func (c BlogService) GetAllBlogs(pagination utils.Pagination) ([]models.Blog, int64, error) {
	return c.repository.GetAllBlogs(pagination)
}
