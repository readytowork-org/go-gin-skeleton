package models

import "time"

type Product struct {
	ProductId    int64     `json:"id"`
	ProductName  string    `json:"product_name"`
	ReceivedQty  uint64    `json:"received_qty"`
	SentQty      uint64    `json:"sent_qty"`
	RemainingQty uint64    `json:"remaining_qty"`
	ReceivedBy   User      `json:"received_by" gorm:"references:Users;"`
	SentBy       User      `json:"sent_by" gorm:"references:Users;"`
	UpdatedAt    time.Time `json:"updated_at"`
	ReceivedAt   time.Time `json:"received_at"`
}

type Tabler interface {
	TableName() string
}

// TableName overrides the table name used by product to `items`
func (Product) TableName() string {
	return "items"
}
