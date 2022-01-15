package models

// {{ucresource}} -> DB model
type {{ucresource}} struct {
	Base
	Title          string     `json:"title"`
}

// TableName  -> returns table name of model
func (c {{ucresource}}) TableName() string {
	return "{{resourcetable}}"
}

// ToMap ->  maps {{plcresource}}
func (c *{{ucresource}}) ToMap() map[string]interface{} {
	return map[string]interface{}{
		"id":              c.ID,
		"title":           c.Title,
	}
}
