package models

// Post -> DB model
type Post struct {
	Base
	Title          string     `json:"title"`
}

// TableName  -> returns table name of model
func (c Post) TableName() string {
	return "posts"
}

// ToMap ->  maps posts
func (c *Post) ToMap() map[string]interface{} {
	return map[string]interface{}{
		"id":              c.ID,
		"title":           c.Title,
	}
}
