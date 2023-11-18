package models

import "gorm.io/gorm"

type Product struct {
	gorm.Model
	UUID     string `gorm:"unique;not null"`
	Name     string `gorm:"not null"`
	ImageURL string `gorm:"not null"`
	AdminID  uint   `gorm:"not null"`
}
