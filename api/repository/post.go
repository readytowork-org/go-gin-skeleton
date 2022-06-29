package repository

import (
	"boilerplate-api/infrastructure"
	"boilerplate-api/models"
	"boilerplate-api/utils"
)

// PostRepository database structure
type PostRepository struct {
	db     infrastructure.Database
	logger infrastructure.Logger
}

// NewPostRepository creates a new Post repository
func NewPostRepository(db infrastructure.Database, logger infrastructure.Logger) PostRepository {
	return PostRepository{
		db:     db,
		logger: logger,
	}
}


// Create Post
func (c PostRepository) Create(Post models.Post) (models.Post, error) {
	return Post, c.db.DB.Create(&Post).Error
}

// GetAllPost -> Get All posts
func (c PostRepository) GetAllPost(pagination utils.Pagination) ([]models.Post, int64, error) {
	var posts []models.Post
	var totalRows int64 = 0
	queryBuider := c.db.DB.Model(&models.Post{}).Offset(pagination.Offset).Order(pagination.Sort)
	
	if !pagination.All{
		queryBuider=queryBuider.Limit(pagination.PageSize)
	}
	
	if pagination.Keyword != "" {
		searchQuery := "%" + pagination.Keyword + "%"
		queryBuider.Where(c.db.DB.Where("`post`.`title` LIKE ?", searchQuery))
	}

	err := queryBuider.
		Find(&posts).
		Offset(-1).
		Limit(-1).
		Count(&totalRows).Error
	return posts, totalRows, err
}

// GetOnePost -> Get One Post By Id
func (c PostRepository) GetOnePost(ID int64) (models.Post, error) {
	Post := models.Post{}
	return Post, c.db.DB.
		Where("id = ?", ID).First(&Post).Error
}

// UpdateOnePost -> Update One Post By Id
func (c PostRepository) UpdateOnePost(Post models.Post) error {
	return c.db.DB.Model(&models.Post{}).
		Where("id = ?", Post.ID).
		Updates(map[string]interface{}{
			"title":           Post.Title,
		}).Error
}

// DeleteOnePost -> Delete One Post By Id
func (c PostRepository) DeleteOnePost(ID int64) error {
	return c.db.DB.
		Where("id = ?", ID).
		Delete(&models.Post{}).
		Error
}
