package repository

import (
	"bilibili-clone/internal/danmaku/model"
	"bilibili-clone/pkg/database"
)

type DanmakuRepository struct{}

func NewDanmakuRepository() *DanmakuRepository {
	return &DanmakuRepository{}
}

func (r *DanmakuRepository) Create(danmaku *model.Danmaku) error {
	return database.DB.Create(danmaku).Error
}

func (r *DanmakuRepository) GetByVideoID(videoID uint) ([]model.Danmaku, error) {
	var danmakus []model.Danmaku
	err := database.DB.Where("video_id = ?", videoID).Order("time asc").Find(&danmakus).Error
	return danmakus, err
}
