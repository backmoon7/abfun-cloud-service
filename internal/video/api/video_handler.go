package api

import (
	"bilibili-clone/internal/video/model"
	"bilibili-clone/internal/video/service"
	"bilibili-clone/pkg/storage"
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
)

type VideoHandler struct {
	Service *service.VideoService
}

func NewVideoHandler() *VideoHandler {
	return &VideoHandler{
		Service: service.NewVideoService(),
	}
}

type FeedVideoResponse struct {
	ID           uint   `json:"id"`
	Title        string `json:"title"`
	CoverURL     string `json:"cover_url"`
	AuthorName   string `json:"author_name"`
	AuthorAvatar string `json:"author_avatar"`
}

type UserResponse struct {
	ID       uint   `json:"ID"`
	Nickname string `json:"Nickname"`
	Avatar   string `json:"Avatar"`
}

func (h *VideoHandler) Upload(c *gin.Context) {
	userIDVal, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}
	userID := userIDVal.(uint)

	file, err := c.FormFile("video")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}


        // File size limit: 500MB
        maxSize := int64(500 * 1024 * 1024)
        if file.Size > maxSize {
                c.JSON(http.StatusBadRequest, gin.H{"error": "file size exceeds 500MB limit"})
                return
        }

        // File type validation
        ext := strings.ToLower(filepath.Ext(file.Filename))
        allowedExts := []string{".mp4", ".avi", ".mov", ".mkv", ".flv", ".wmv"}
        isAllowed := false
        for _, allowed := range allowedExts {
                if ext == allowed {
                        isAllowed = true
                        break
                }
        }
        if !isAllowed {
                c.JSON(http.StatusBadRequest, gin.H{"error": "only video files are allowed (mp4, avi, mov, mkv, flv, wmv)"})
                return
        }
	filename := filepath.Base(file.Filename)
	savePath := filepath.Join(storage.LocalStoragePath, filename)
	if err := c.SaveUploadedFile(file, savePath); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	url, err := storage.LocalStorage.Upload(savePath)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Generate Cover using FFmpeg
	coverFilename := strings.TrimSuffix(filename, filepath.Ext(filename)) + ".jpg"
	coverPath := filepath.Join(storage.LocalStoragePath, coverFilename)

	// ffmpeg -i input.mp4 -ss 00:00:01 -vframes 1 output.jpg
	cmd := exec.Command("ffmpeg", "-i", savePath, "-ss", "00:00:01", "-vframes", "1", coverPath)
	if err := cmd.Run(); err != nil {
		fmt.Printf("Failed to generate cover: %v\n", err)
	}

	coverURL := ""
	// Simple file existence check by trying to upload
	if cURL, err := storage.LocalStorage.Upload(coverPath); err == nil {
		coverURL = cURL
	}

	title := c.PostForm("title")
	description := c.PostForm("description")

        if strings.TrimSpace(title) == "" {
                c.JSON(http.StatusBadRequest, gin.H{"error": "title is required"})
                return
        }

	// Parse tag_ids array from form
	tagIDStrs := c.PostFormArray("tag_ids[]")
	var tagIDs []uint
	for _, idStr := range tagIDStrs {
		if id, err := strconv.ParseUint(idStr, 10, 32); err == nil {
			tagIDs = append(tagIDs, uint(id))
		}
	}

	video := &model.Video{
		Title:       title,
		Description: description,
		VideoURL:    url, // Fixed field name from URL to VideoURL based on model definition
		CoverURL:    coverURL,
		UserID:      userID,
		TagIDs:      tagIDs,
	}

	if err := h.Service.CreateVideo(video); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": video})
}

func (h *VideoHandler) GetFeed(c *gin.Context) {
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))
	offset, _ := strconv.Atoi(c.DefaultQuery("offset", "0"))
	tagID, _ := strconv.ParseUint(c.Query("tag_id"), 10, 32)

        // Validate pagination parameters
        if limit < 1 {
                limit = 20
        }
        if limit > 100 {
                limit = 100
        }
        if offset < 0 {
                offset = 0
        }

	videos, err := h.Service.GetFeed(limit, offset, uint(tagID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Collect UserIDs
	userIDs := make([]uint, 0)
	for _, v := range videos {
		userIDs = append(userIDs, v.UserID)
	}

	// Fetch Users from User Service
	userMap := make(map[uint]UserResponse)
	if len(userIDs) > 0 {
		reqBody, _ := json.Marshal(map[string][]uint{"ids": userIDs})
		resp, err := http.Post("http://user-service:8081/api/v1/user/batch", "application/json", bytes.NewBuffer(reqBody))
		if err == nil && resp.StatusCode == http.StatusOK {
			defer resp.Body.Close()
			var userResp struct {
				Users []UserResponse `json:"users"`
			}
			if err := json.NewDecoder(resp.Body).Decode(&userResp); err == nil {
				for _, u := range userResp.Users {
					userMap[u.ID] = u
				}
			}
		} else {
			fmt.Printf("Failed to fetch users: %v\n", err)
		}
	}

	// Build Response
	var response []FeedVideoResponse
	for _, v := range videos {
		u := userMap[v.UserID]
		response = append(response, FeedVideoResponse{
			ID:           v.ID,
			Title:        v.Title,
			CoverURL:     v.CoverURL,
			AuthorName:   u.Nickname,
			AuthorAvatar: u.Avatar,
		})
	}

	c.JSON(http.StatusOK, gin.H{"data": response})
}

func (h *VideoHandler) Feed(c *gin.Context) {
	h.GetFeed(c)
}

func (h *VideoHandler) ListTags(c *gin.Context) {
	tags, err := h.Service.ListTags()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": tags})
}

func (h *VideoHandler) Search(c *gin.Context) {
	query := c.Query("q")
	if query == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "query is required"})
		return
	}

        // Query validation
        query = strings.TrimSpace(query)
        if len(query) > 100 {
                c.JSON(http.StatusBadRequest, gin.H{"error": "query too long (max 100 characters)"})
                return
        }

	videos, err := h.Service.SearchVideos(query)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Collect UserIDs
	userIDs := make([]uint, 0)
	for _, v := range videos {
		userIDs = append(userIDs, v.UserID)
	}

	// Fetch Users from User Service
	userMap := make(map[uint]UserResponse)
	if len(userIDs) > 0 {
		reqBody, _ := json.Marshal(map[string][]uint{"ids": userIDs})
		resp, err := http.Post("http://user-service:8081/api/v1/user/batch", "application/json", bytes.NewBuffer(reqBody))
		if err == nil && resp.StatusCode == http.StatusOK {
			defer resp.Body.Close()
			var userResp struct {
				Users []UserResponse `json:"users"`
			}
			if err := json.NewDecoder(resp.Body).Decode(&userResp); err == nil {
				for _, u := range userResp.Users {
					userMap[u.ID] = u
				}
			}
		} else {
			fmt.Printf("Failed to fetch users: %v\n", err)
		}
	}

	// Build Response
	var response []FeedVideoResponse
	for _, v := range videos {
		u := userMap[v.UserID]
		response = append(response, FeedVideoResponse{
			ID:           v.ID,
			Title:        v.Title,
			CoverURL:     v.CoverURL,
			AuthorName:   u.Nickname,
			AuthorAvatar: u.Avatar,
		})
	}

	c.JSON(http.StatusOK, gin.H{"data": response})
}

func (h *VideoHandler) GetVideoDetail(c *gin.Context) {
videoID, err := strconv.ParseUint(c.Param("id"), 10, 32)
if err != nil {
c.JSON(http.StatusBadRequest, gin.H{"error": "invalid video id"})
return
}

video, err := h.Service.GetVideoByID(uint(videoID))
if err != nil {
c.JSON(http.StatusNotFound, gin.H{"error": "video not found"})
return
}

c.JSON(http.StatusOK, gin.H{"data": video})
}
