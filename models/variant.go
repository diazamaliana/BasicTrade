package models

import "gorm.io/gorm"

type Variant struct {
	gorm.Model
	UUID        string `gorm:"unique;not null"`
	VariantName string `gorm:"not null"`
	Quantity    int    `gorm:"not null"`
	ProductID   uint   `gorm:"not null"`
}
