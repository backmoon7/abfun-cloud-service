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

func (r *VideoRepository) UpdateStatus(videoID uint, status int, videoURL, coverURL string) error {
	updates := map[string]interface{}{"status": status}
	if videoURL != "" {
		updates["video_url"] = videoURL
	}
	if coverURL != "" {
		updates["cover_url"] = coverURL
	}
	return database.DB.Model(&model.Video{}).Where("id = ?", videoID).Updates(updates).Error
}

func (r *VideoRepository) FindByID(id uint) (*model.Video, []uint, error) {
	var video model.Video
	if err := database.DB.First(&video, id).Error; err != nil {
		return nil, nil, err
	}

	var tagIDs []uint
	err := database.DB.Model(&model.VideoTag{}).Where("video_id = ?", id).Pluck("tag_id", &tagIDs).Error
	return &video, tagIDs, err
}

func (r *VideoRepository) GetVideoTagIDs(videoID uint) ([]uint, error) {
        var videoTags []model.VideoTag
        err := database.DB.Where("video_id = ?", videoID).Find(&videoTags).Error
        if err != nil {
                return nil, err
        }
        
        var tagIDs []uint
        for _, vt := range videoTags {
                tagIDs = append(tagIDs, vt.TagID)
        }
        return tagIDs, nil
}
