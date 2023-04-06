package models

type Category struct {
	Base
	Title string `json:"title"`
}

func (c Category) TableName() string {
	return "categories"
}

func (c Category) ToMap() map[string]interface{} {
	return map[string]interface{}{
		"title": c.Title,
	}
}
