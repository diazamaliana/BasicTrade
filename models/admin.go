package models

import "gorm.io/gorm"

type Admin struct {
	gorm.Model
	UUID     string `gorm:"unique;not null"`
	Name     string `gorm:"not null"`
	Email    string `gorm:"unique;not null"`
	Password string `gorm:"not null"`
}
