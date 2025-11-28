package model

import "gorm.io/gorm"

type Video struct {
	gorm.Model
	UserID      uint   `gorm:"index"`
	Title       string `gorm:"size:255;not null"`
	Description string `gorm:"type:text"`
	VideoURL    string `gorm:"size:255;not null"`
	CoverURL    string `gorm:"size:255"`
	Status      int    `gorm:"default:0"` // 0: processing, 1: published
}

type Tag struct {
	gorm.Model
	Name string `gorm:"uniqueIndex;size:50"`
}

type VideoTag struct {
	VideoID uint `gorm:"primaryKey"`
	TagID   uint `gorm:"primaryKey"`
}
