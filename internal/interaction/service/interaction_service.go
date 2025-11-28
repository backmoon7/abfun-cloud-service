package service

import (
	"bilibili-clone/internal/interaction/model"
	"bilibili-clone/internal/interaction/repository"
	"bilibili-clone/pkg/mq"
	"encoding/json"
)

type InteractionService struct {
	repo *repository.InteractionRepository
}

func NewInteractionService() *InteractionService {
	return &InteractionService{
		repo: repository.NewInteractionRepository(),
	}
}

func (s *InteractionService) LikeVideo(userID, videoID uint) error {
	if err := s.repo.CreateLike(userID, videoID); err != nil {
		return err
	}
	return s.publishEvent("like", userID, videoID)
}

func (s *InteractionService) UnlikeVideo(userID, videoID uint) error {
	if err := s.repo.DeleteLike(userID, videoID); err != nil {
		return err
	}
	return s.publishEvent("unlike", userID, videoID)
}

func (s *InteractionService) publishEvent(eventType string, userID, videoID uint) error {
	event := map[string]interface{}{
		"type":     eventType,
		"user_id":  userID,
		"video_id": videoID,
	}
	bytes, err := json.Marshal(event)
	if err != nil {
		return err
	}
	return mq.Publish("interaction_event", bytes)
}

func (s *InteractionService) CreateFolder(userID uint, name string) error {
	folder := &model.FavoriteFolder{
		UserID: userID,
		Name:   name,
	}
	return s.repo.CreateFolder(folder)
}

func (s *InteractionService) GetFolders(userID uint) ([]model.FavoriteFolder, error) {
	return s.repo.GetFolders(userID)
}

func (s *InteractionService) AddFavorite(folderID, videoID uint) error {
	item := &model.FavoriteItem{
		FolderID: folderID,
		VideoID:  videoID,
	}
	return s.repo.AddFavorite(item)
}
