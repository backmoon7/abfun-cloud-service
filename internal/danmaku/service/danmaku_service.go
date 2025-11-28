package service

import (
	"bilibili-clone/internal/danmaku/model"
	"bilibili-clone/internal/danmaku/repository"
)

type DanmakuService struct {
	repo *repository.DanmakuRepository
}

func NewDanmakuService() *DanmakuService {
	return &DanmakuService{
		repo: repository.NewDanmakuRepository(),
	}
}

func (s *DanmakuService) AddDanmaku(d *model.Danmaku) error {
	return s.repo.Create(d)
}

func (s *DanmakuService) GetDanmakus(videoID uint) ([]model.Danmaku, error) {
	return s.repo.GetByVideoID(videoID)
}
