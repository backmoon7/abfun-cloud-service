package model

import "gorm.io/gorm"

type Comment struct {
gorm.Model
UserID  uint   `gorm:"index;not null"`
VideoID uint   `gorm:"index;not null"`
Content string `gorm:"type:text;not null"`
}
