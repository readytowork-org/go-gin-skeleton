package models

type Category struct {
	Base
	Title string `json:"title" validate:"required"`
}

func (m *Category) TableName() string {
	return "category"
}

// ToMap convert User to map
func (m Category) ToMap() map[string]interface{} {
	return map[string]interface{}{
		"title": m.Title,
	}
}

type BlogCategories struct {
	Base
	BlogId     int64 `json:"blog_id" validate:"required"`
	CategoryId int64 `json:"category_id" validate:"required"`
}

func (m *BlogCategories) TableName() string {
	return "blog_categories"
}

// ToMap convert User to map
func (m BlogCategories) ToMap() map[string]interface{} {
	return map[string]interface{}{
		"blog_id":     m.BlogId,
		"category_id": m.CategoryId,
	}
}

type Blog struct {
	Base
	Title          int64  `json:"title" validate:"required"`
	Content        string `json:"content" validate:"required"`
	ThumbnailImage string `json:"thumbnail_image" validate:"required"`
	IsPublished    bool   `json:"is_published" validate:"required"`
	CreatedBy      int64  `json:"created_by" validate:"required"`
	UpdatedBy      int64  `json:"updated_by" validate:"required"`
}

func (m *Blog) TableName() string {
	return "blog_categories"
}

// ToMap convert User to map
func (m Blog) ToMap() map[string]interface{} {
	return map[string]interface{}{
		"title":           m.Title,
		"content":         m.Content,
		"thumbnail_image": m.ThumbnailImage,
		"is_published":    m.IsPublished,
		"created_by":      m.CreatedBy,
		"updated_by":      m.UpdatedBy,
	}
}
