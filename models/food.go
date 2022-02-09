package models

// Food -> DB model
type Food struct {
	Base
	Title          string     `json:"title"`
}

// TableName  -> returns table name of model
func (c Food) TableName() string {
	return "foods"
}

// ToMap ->  maps foods
func (c *Food) ToMap() map[string]interface{} {
	return map[string]interface{}{
		"id":              c.ID,
		"title":           c.Title,
	}
}
