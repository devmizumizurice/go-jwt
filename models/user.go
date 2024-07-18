package models

type User struct {
	General
	Name     string `gorm:"size:255;not null"`
	Email    string `gorm:"unique;size:255;not null"`
	Password string `gorm:"size:255;not null"`
}
