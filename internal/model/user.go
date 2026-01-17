package model

import "time"

type User struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	Name      string    `gorm:"type:varchar(255);not null" json:"name" validate:"required"`
	Email     string    `gorm:"type:varchar(255);unique;not null" json:"email" validate:"required,email"`
	Password  string    `gorm:"type:varchar(255);not null" json:"password" validate:"required,min=6"` // json:"-" agar tidak ditampilkan di response
	Role      string    `gorm:"type:varchar(50);default:'customer'" json:"role"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}