package models

// StudentInfo -> DB model
type StudentInfo struct {
	Base
	Name    string `json:"name"`
	Address string `json:"address"`
	Class   int    `json:"class"`
	Detail  string `json:"detail"`
}

// TableName  -> returns table name of model
func (c StudentInfo) TableName() string {
	return "student_infos"
}

// ToMap ->  maps studentInfos
func (c *StudentInfo) ToMap() map[string]interface{} {
	return map[string]interface{}{
		"id":      c.ID,
		"name":    c.Name,
		"address": c.Address,
		"class":   c.Class,
		"detail":  c.Detail,
	}
}
