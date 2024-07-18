package models

import (
	"time"

	"gorm.io/gorm"
)

type General struct {
	ID        string `gorm:"primaryKey;size:255;default:uuid_generate_v4()"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`
}
