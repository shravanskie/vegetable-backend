package models

import (
	"time"
)

type Cart struct {
	ID          uint `gorm:"primaryKey"`
	UserID      uint `gorm:"index;not null"`
	VegetableID uint `gorm:"index;not null"`
	Quantity    int  `gorm:"not null;default:1"`
	CreatedAt   time.Time
	UpdatedAt   time.Time
}
