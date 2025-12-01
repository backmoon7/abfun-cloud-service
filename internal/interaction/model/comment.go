package model

import "gorm.io/gorm"

type Comment struct {
        gorm.Model `json:",inline"`
        UserID     uint   `gorm:"index;not null;constraint:OnDelete:CASCADE" json:"user_id"`
        VideoID    uint   `gorm:"index;not null;constraint:OnDelete:CASCADE" json:"video_id"`
        Content    string `gorm:"type:text;not null" json:"content"`
}
