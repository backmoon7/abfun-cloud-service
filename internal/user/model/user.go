package model

import "gorm.io/gorm"

type User struct {
        gorm.Model `json:",inline"`
        Email      string `gorm:"uniqueIndex;size:255;not null" json:"email"`
        Password   string `gorm:"not null" json:"-"`
        Nickname   string `gorm:"size:50" json:"nickname"`
        Avatar     string `gorm:"size:255" json:"avatar"`
}
