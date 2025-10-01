package models

import "gorm.io/gorm"

type Vegetable struct {
	gorm.Model
	Name         string  `json:"name"`
	Unit         string  `json:"unit"`
	ImageURL     string  `json:"image_url"`
	UnitQuantity float64 `json:"unit_quantity"` // e.g., 1 for "kg", 0.5 for "half kg"
}

// Datacenter Master
type Datacenter struct {
	gorm.Model
	Name     string `json:"name" gorm:"unique;not null"`
	Location string `json:"location"`
}

// VegetablePrice maps Vegetable <-> Datacenter <-> Unit -> Price
type VegetablePrice struct {
	gorm.Model
	VegetableID  uint       `json:"vegetable_id"`
	Vegetable    Vegetable  `gorm:"foreignKey:VegetableID"`
	DatacenterID uint       `json:"datacenter_id"`
	Datacenter   Datacenter `gorm:"foreignKey:DatacenterID"`
	Unit         string     `json:"unit" gorm:"not null"` // "kg", "gram", "piece"
	Price        float64    `json:"price"`
}

type VegetableInput struct {
	Name         string  `json:"name" binding:"required"`
	Unit         string  `json:"unit" binding:"required"`
	UnitQuantity float64 `json:"unit_quantity" binding:"required"`
	ImagePath    string  `json:"image_path" binding:"required"`
}
