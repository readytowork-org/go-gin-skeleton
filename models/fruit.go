package models

// Fruit -> DB model
type Fruit struct {
	Base
	Name        string  `json:"name"`
	Season		string	`json:"season"`
	Metric		string	`json:"metric"`
	Price		string	`json:"price"`
}

// TableName  -> returns table name of model
func (c Fruit) TableName() string {
	return "fruits"
}

// ToMap ->  maps fruits
func (c *Fruit) ToMap() map[string]interface{} {
	return map[string]interface{}{
		"id":              	c.ID,
		"name":           	c.Name,
		"season":			c.Season,
		"metric":			c.Metric,
		"Price":			c.Price,
	}
}
