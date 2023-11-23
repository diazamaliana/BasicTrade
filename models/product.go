package models

import (
	"time"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Product struct {
	ID        uint   `gorm:"primaryKey" json:"id"`
	UUID      string `gorm:"type:varchar(36);unique;not null" json:"uuid"`
	ProductName     string `gorm:"not null" json:"product_name"`
	ImageURL string `json:"image_url"`
	AdminUUID  uuid.UUID `gorm:"type:varchar(36);not null" json:"admin_uuid"`
	CreatedAt time.Time `json:"created_at,omitempty"`
	UpdatedAt time.Time `json:"updated_at,omitempty"`
	Variants  []Variant `gorm:"foreignKey:ProductUUID;references:UUID"`
}

// BeforeCreate generates a UUID for the admin before creating a record.
func (product *Product) BeforeCreate(tx *gorm.DB) error {
	product.UUID = uuid.New().String()
	return nil
}
