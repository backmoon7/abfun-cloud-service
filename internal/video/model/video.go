package model

import "gorm.io/gorm"

type Video struct {
        gorm.Model   `json:",inline"`
        UserID       uint   `gorm:"index;constraint:OnDelete:CASCADE" json:"user_id"`
        Title        string `gorm:"size:255;not null" json:"title"`
        Description  string `gorm:"type:text" json:"description"`
        VideoURL     string `gorm:"size:255;not null" json:"video_url"`
        CoverURL     string `gorm:"size:255" json:"cover_url"`
        Status       int    `gorm:"index;default:0" json:"status"`
        LikeCount    int    `gorm:"default:0" json:"like_count"`
        CommentCount int    `gorm:"default:0" json:"comment_count"`
        TagIDs       []uint `gorm:"-" json:"tag_ids"`
}

type Tag struct {
        gorm.Model `json:",inline"`
        Name       string `gorm:"uniqueIndex;size:50" json:"name"`
}

type VideoTag struct {
        VideoID uint `gorm:"primaryKey;index;constraint:OnDelete:CASCADE" json:"video_id"`
        TagID   uint `gorm:"primaryKey;index;constraint:OnDelete:CASCADE" json:"tag_id"`
}
