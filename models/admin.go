// models/admin.go

package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// Admin represents the admin model.
type Admin struct {
	ID        uint   `gorm:"primaryKey" json:"id"`
	UUID      string `gorm:"type:varchar(36);unique;not null" json:"uuid"`
	Name      string `gorm:"not null" json:"name"`
	Email     string `gorm:"unique;not null" json:"email"`
	Password  string `gorm:"not null" json:"password"`
	CreatedAt time.Time `json:"created_at,omitempty"`
	UpdatedAt time.Time `json:"updated_at,omitempty"`
}

// BeforeCreate generates a UUID for the admin before creating a record.
func (admin *Admin) BeforeCreate(tx *gorm.DB) error {
	admin.UUID = uuid.New().String()
	return nil
}
