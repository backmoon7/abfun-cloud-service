package service

import (
"bilibili-clone/internal/video/model"
"bilibili-clone/internal/video/repository"
"bilibili-clone/pkg/mq"
"encoding/json"
"fmt"
"os"

"github.com/meilisearch/meilisearch-go"
)

type VideoService struct {
repo        *repository.VideoRepository
meiliClient *meilisearch.Client
}

func NewVideoService() *VideoService {
meiliHost := os.Getenv("MEILISEARCH_HOST")
meiliKey := os.Getenv("MEILISEARCH_API_KEY")

var client *meilisearch.Client
if meiliHost != "" {
client = meilisearch.NewClient(meilisearch.ClientConfig{
Host:   meiliHost,
APIKey: meiliKey,
})
// Ensure index exists
_, _ = client.CreateIndex(&meilisearch.IndexConfig{
Uid:        "videos",
PrimaryKey: "id",
})
// Update filterable attributes
_, _ = client.Index("videos").UpdateFilterableAttributes(&[]string{"tag_ids", "user_id"})
// Update searchable attributes
_, _ = client.Index("videos").UpdateSearchableAttributes(&[]string{"title", "description"})
}

return &VideoService{
repo:        repository.NewVideoRepository(),
meiliClient: client,
}
}

func (s *VideoService) CreateVideo(video *model.Video) error {
// Set status to Processing (0)
video.Status = 0
if err := s.repo.Create(video, video.TagIDs); err != nil {
return err
}

// Publish event
event := map[string]interface{}{
"video_id":  video.ID,
"video_url": video.VideoURL,
}
bytes, err := json.Marshal(event)
if err != nil {
return err
}
return mq.Publish("video_uploaded", bytes)
}

func (s *VideoService) UpdateVideoStatus(videoID uint, status int, videoURL, coverURL string) error {
return s.repo.UpdateStatus(videoID, status, videoURL, coverURL)
}

func (s *VideoService) GetFeed(limit, offset int, tagID uint) ([]model.Video, error) {
return s.repo.GetVideos(limit, offset, tagID)
}

func (s *VideoService) ListTags() ([]model.Tag, error) {
return s.repo.GetAllTags()
}

func (s *VideoService) InitTags() {
s.repo.InitTags()
}

func (s *VideoService) SyncVideoToMeilisearch(videoID uint) error {
if s.meiliClient == nil {
return nil
}

video, tagIDs, err := s.repo.FindByID(videoID)
if err != nil {
return err
}

doc := map[string]interface{}{
"id":          video.ID,
"title":       video.Title,
"description": video.Description,
"cover_url":   video.CoverURL,
"user_id":     video.UserID,
"tag_ids":     tagIDs,
"created_at":  video.CreatedAt.Unix(),
}

_, err = s.meiliClient.Index("videos").AddDocuments(doc)
return err
}

func (s *VideoService) SearchVideos(query string) ([]model.Video, error) {
if s.meiliClient == nil {
return nil, fmt.Errorf("meilisearch client not initialized")
}

searchRes, err := s.meiliClient.Index("videos").Search(query, &meilisearch.SearchRequest{
Limit: 20,
})
if err != nil {
return nil, err
}

var videos []model.Video
hits, ok := searchRes.Hits.([]interface{})
if !ok {
return nil, fmt.Errorf("invalid search response")
}

for _, hit := range hits {
h := hit.(map[string]interface{})
v := model.Video{
ID:          uint(h["id"].(float64)),
Title:       h["title"].(string),
Description: h["description"].(string),
CoverURL:    h["cover_url"].(string),
UserID:      uint(h["user_id"].(float64)),
}
videos = append(videos, v)
}

return videos, nil
}
