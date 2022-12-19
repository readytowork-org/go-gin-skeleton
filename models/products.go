package models

import (
	"errors"
	"time"

	"gorm.io/gorm"
)

type Product struct {
	ItemId       int64     `json:"id"`
	ItemName     string    `json:"product_name"`
	ReceivedQty  uint64    `json:"received_qty"`
	SentQty      uint64    `json:"sent_qty"`
	RemainingQty uint64    `json:"remaining_qty"`
	ReceivedBy   int64     `json:"received_by"`
	ReceivedUser User      `gorm:"foreignKey:ReceivedBy" `
	SentBy       int64     `json:"sent_by" gorm:"-" `
	SentUser     User      `gorm:"foreignKey:SentBy" `
	UpdatedAt    time.Time `json:"updated_at"`
	CreatedAt    time.Time `json:"created_at"`
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

type ProductSentInput struct {
	SentQty uint64 `json:"sent_qty,string" binding:"required"`
	SentBy  int64  `json:"sent_by,string" binding:"required"`
}

func (products *Product) AfterUpdate(tx *gorm.DB) error {
	item := Product{}
	updateProduct := tx.Model(&Product{}).Where("item_id=?", products.ItemId).First(&item)
	if products.SentQty != 0 {
		if item.ReceivedQty < products.SentQty {
			err := errors.New("Sent qty is  greater than received qty")
			return err
		}
		remainingQty := item.ReceivedQty - products.SentQty
		return updateProduct.Update("remaining_qty", remainingQty).Error
	}
	return nil
}
