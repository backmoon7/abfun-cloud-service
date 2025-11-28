package repository

import (
	"bilibili-clone/internal/interaction/model"
	"bilibili-clone/pkg/database"

	"gorm.io/gorm/clause"
)

type InteractionRepository struct{}

func NewInteractionRepository() *InteractionRepository {
	return &InteractionRepository{}
}

func (r *InteractionRepository) CreateFolder(folder *model.FavoriteFolder) error {
	return database.DB.Create(folder).Error
}

func (r *InteractionRepository) GetFolders(userID uint) ([]model.FavoriteFolder, error) {
	var folders []model.FavoriteFolder
	err := database.DB.Where("user_id = ?", userID).Find(&folders).Error
	return folders, err
}

func (r *InteractionRepository) AddFavorite(item *model.FavoriteItem) error {
	return database.DB.Create(item).Error
}

func (r *InteractionRepository) CreateLike(userID, videoID uint) error {
	like := &model.VideoLike{UserID: userID, VideoID: videoID}
	return database.DB.Clauses(clause.OnConflict{DoNothing: true}).Create(like).Error
}

func (r *InteractionRepository) DeleteLike(userID, videoID uint) error {
	return database.DB.Where("user_id = ? AND video_id = ?", userID, videoID).Delete(&model.VideoLike{}).Error
}
