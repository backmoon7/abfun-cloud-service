package model

import "gorm.io/gorm"

type Danmaku struct {
        gorm.Model `json:",inline"`
        VideoID    uint    `gorm:"index;not null;constraint:OnDelete:CASCADE" json:"video_id"`
        UserID     uint    `gorm:"index;not null;constraint:OnDelete:CASCADE" json:"user_id"`
        Content    string  `gorm:"size:255;not null" json:"content"`
        Time       float64 `gorm:"not null" json:"time"`
        Color      string  `gorm:"size:20;default:'#FFFFFF'" json:"color"`
}
