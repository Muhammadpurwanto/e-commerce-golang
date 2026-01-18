package model

import "time"

type CartItem struct {
	ID        uint    `gorm:"primaryKey" json:"id"`
	CartID    uint    `gorm:"not null;uniqueIndex:idx_cart_product" json:"cart_id"`
	ProductID uint    `gorm:"not null;uniqueIndex:idx_cart_product" json:"product_id"`
	Product   Product `gorm:"foreignKey:ProductID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"product"`
	Quantity  uint    `gorm:"not null" json:"quantity"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}