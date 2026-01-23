package model

import "time"

// Product merepresentasikan model untuk tabel produk
type Product struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	Name      string    `gorm:"type:varchar(255);not null" json:"name"`
	SKU       string    `gorm:"type:varchar(100);unique;not null" json:"sku"`
	Price     float64   `gorm:"type:decimal(10,2);not null" json:"price"`
	ImageURL  string    `gorm:"type:varchar(255)" json:"image_url"`
	Stock     uint      `gorm:"not null" json:"stock"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}