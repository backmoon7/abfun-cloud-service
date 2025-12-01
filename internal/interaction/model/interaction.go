package model

import "gorm.io/gorm"

type VideoLike struct {
        gorm.Model `json:",inline"`
        UserID     uint `gorm:"uniqueIndex:idx_user_video;constraint:OnDelete:CASCADE" json:"user_id"`
        VideoID    uint `gorm:"uniqueIndex:idx_user_video;constraint:OnDelete:CASCADE" json:"video_id"`
}

type FavoriteFolder struct {
        gorm.Model `json:",inline"`
        UserID     uint   `gorm:"index;constraint:OnDelete:CASCADE" json:"user_id"`
        Name       string `gorm:"size:50" json:"name"`
}

type FavoriteItem struct {
        gorm.Model `json:",inline"`
        FolderID   uint `gorm:"uniqueIndex:idx_folder_video;constraint:OnDelete:CASCADE" json:"folder_id"`
        VideoID    uint `gorm:"uniqueIndex:idx_folder_video;constraint:OnDelete:CASCADE" json:"video_id"`
}
