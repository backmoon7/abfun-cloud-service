package service

import (
	"bilibili-clone/internal/video/model"
	"bilibili-clone/internal/video/repository"
	"bilibili-clone/pkg/storage"
	"mime/multipart"
	"path/filepath"
	"strconv"
	"time"
)

type VideoService struct {
	repo *repository.VideoRepository
}

func NewVideoService() *VideoService {
	return &VideoService{
		repo: repository.NewVideoRepository(),
	}
}

func (s *VideoService) UploadVideo(userID uint, title, desc string, tagIDs []uint, file *multipart.FileHeader) error {
	// Upload to Local Storage
	filename := strconv.FormatInt(time.Now().UnixNano(), 10) + filepath.Ext(file.Filename)
	url, err := storage.UploadFile(file, filename)
	if err != nil {
		return err
	}

	video := &model.Video{
		UserID:      userID,
		Title:       title,
		Description: desc,
		VideoURL:    url,
		Status:      1, // Published
	}

	return s.repo.Create(video, tagIDs)
}

func (s *VideoService) GetFeed(limit, offset int, tagID uint) ([]model.Video, error) {
	return s.repo.GetVideos(limit, offset, tagID)
}

func (s *VideoService) GetTags() ([]model.Tag, error) {
	return s.repo.GetAllTags()
}

func (s *VideoService) InitTags() {
	s.repo.InitTags()
}
