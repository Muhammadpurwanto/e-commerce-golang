package model

import "time"

type Order struct {
	ID          uint        `gorm:"primaryKey" json:"id"`
	UserID      uint        `gorm:"not null" json:"user_id"`
	User        User        `gorm:"foreignKey:UserID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;" json:"user"`
	OrderItems  []OrderItem `gorm:"foreignKey:OrderID" json:"order_items"`
	TotalAmount float64     `gorm:"type:decimal(10,2);not null" json:"total_amount"`
	Status      string      `gorm:"type:varchar(50);default:'pending'" json:"status"` // pending, paid, shipped, cancelled
	CreatedAt   time.Time   `json:"created_at"`
	UpdatedAt   time.Time   `json:"updated_at"`
}