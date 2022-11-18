package models

type Products struct {
	Base
	Title string `json:"title" binding:"required"`
	Stock int64  `json:"stock" binding:"required"`
}

// TableName gives table name of model
func (p Products) TableName() string {
	return "products"
}
