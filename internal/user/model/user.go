package model

import "gorm.io/gorm"

type User struct {
gorm.Model
Username string `gorm:"uniqueIndex;size:255;not null"`
Password string `gorm:"not null"`
Nickname string `gorm:"size:50"`
Avatar   string `gorm:"size:255"`
}
