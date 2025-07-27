package models

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
	"time"
)

// User represents a user in the system.
type User struct {
	ID           uuid.UUID      `gorm:"primaryKey;type:uuid;default:uuid_generate_v4()" json:"id"`
	FullName     string         `gorm:"type:varchar(100);not null" json:"full_name" validate:"required,min=2,max=100"`
	Email        string         `gorm:"type:varchar(255);uniqueIndex;not null" json:"email" validate:"required,email"`
	PasswordHash string         `gorm:"type:varchar(255);not null" json:"password" validate:"required,min=8,max=100"`
	Verified     bool           `gorm:"default:false" json:"verified"`
	CreatedAt    time.Time      `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt    time.Time      `gorm:"autoUpdateTime" json:"updated_at"`
	DeletedAt    gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty"`
}
