package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Variant struct {
	ID        uint   `gorm:"primaryKey" json:"id"`
	UUID      string `gorm:"type:varchar(36);unique;not null" json:"uuid"`
	VariantName string `gorm:"not null" json:"variant_name"`
	Quantity    uint    `gorm:"not null" json:"quantity"`
	ProductUUID  uuid.UUID `gorm:"type:varchar(36);not null" json:"admin_uuid"`
	CreatedAt time.Time `json:"created_at,omitempty"`
	UpdatedAt time.Time `json:"updated_at,omitempty"`
}

func (variant *Variant) BeforeCreate(tx *gorm.DB) error {
	variant.UUID = uuid.New().String()
	return nil
}
