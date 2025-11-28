package repository

import (
	"bilibili-clone/internal/video/model"
	"bilibili-clone/pkg/database"
	"gorm.io/gorm"
)

type VideoRepository struct{}

func NewVideoRepository() *VideoRepository {
	return &VideoRepository{}
}

func (r *VideoRepository) Create(video *model.Video, tagIDs []uint) error {
	return database.DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(video).Error; err != nil {
			return err
		}
		for _, tagID := range tagIDs {
			if err := tx.Create(&model.VideoTag{VideoID: video.ID, TagID: tagID}).Error; err != nil {
				return err
			}
		}
		return nil
	})
}

func (r *VideoRepository) GetVideos(limit, offset int, tagID uint) ([]model.Video, error) {
	var videos []model.Video
	db := database.DB.Limit(limit).Offset(offset).Order("created_at desc")
	
	if tagID > 0 {
		db = db.Joins("JOIN video_tags ON video_tags.video_id = videos.id").Where("video_tags.tag_id = ?", tagID)
	}
	
	err := db.Find(&videos).Error
	return videos, err
}

func (r *VideoRepository) GetAllTags() ([]model.Tag, error) {
	var tags []model.Tag
	err := database.DB.Find(&tags).Error
	return tags, err
}

func (r *VideoRepository) InitTags() {
	tags := []string{"Anime", "Gaming", "Music", "Tech", "Life"}
	for _, name := range tags {
		database.DB.FirstOrCreate(&model.Tag{Name: name}, model.Tag{Name: name})
	}
}
