package models

import "time"

type Product struct {
	ItemId       int64  `json:"id"`
	ItemName     string `json:"product_name"`
	ReceivedQty  uint64 `json:"received_qty"`
	SentQty      uint64 `json:"sent_qty"`
	RemainingQty uint64 `json:"remaining_qty"`
	ReceivedBy   int64  `json:"received_by"`
	ReceivedUser User   `gorm:"foreignKey:ReceivedBy" json:"received_user"`
	// SentBy       int64     `json:"sent_by"`
	UpdatedAt time.Time `json:"updated_at"`
	CreatedAt time.Time `json:"created_at"`
	// SentUser     User      `gorm:"foreignKey:SentBy;"`
}

type Tabler interface {
	TableName() string
}

// TableName overrides the table name used by product to `items`
func (Product) TableName() string {
	return "items"
}

type ProductCreateInput struct {
	ProductName string `json:"product_name" binding:"required"`
	ReceivedQty uint64 `json:"received_qty,string" binding:"required"`
	ReceivedBy  int64  `json:"received_by,string" binding:"required"`
}
