package model

import "gorm.io/gorm"

type User struct {
	gorm.Model
	Username string `gorm:"uniqueIndex;size:255;not null"`
	Password string `gorm:"not null"`
}
