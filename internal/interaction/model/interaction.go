package model

import "gorm.io/gorm"

type VideoLike struct {
	gorm.Model
	UserID  uint `gorm:"uniqueIndex:idx_user_video"`
	VideoID uint `gorm:"uniqueIndex:idx_user_video"`
}

type FavoriteFolder struct {
	gorm.Model
	UserID uint   `gorm:"index"`
	Name   string `gorm:"size:50"`
}

type FavoriteItem struct {
	gorm.Model
	FolderID uint `gorm:"uniqueIndex:idx_folder_video"`
	VideoID  uint `gorm:"uniqueIndex:idx_folder_video"`
}
