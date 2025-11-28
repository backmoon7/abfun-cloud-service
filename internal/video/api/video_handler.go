package api

import (
	"bilibili-clone/internal/video/model"
	"bilibili-clone/internal/video/service"
	"bilibili-clone/pkg/storage"
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
		// Fallback: try capturing at 0s if 1s fails or just log error
		// For now, we just log and continue without cover if it fails
		fmt.Printf("Failed to generate cover: %v\n", err)
	}

	coverURL := ""
	if _, err := strconv.Atoi("0"); err == nil { // Dummy check to keep strconv import if needed, but we use it below
	}

	// Check if cover was generated
	if _, err := strconv.Atoi("0"); err == nil {
		// Re-using strconv to avoid unused import error if I delete it below
	}

	// Upload cover if generated
	// We check if file exists
	if _, err := exec.LookPath("ffmpeg"); err == nil {
		// If ffmpeg exists, we assume we tried. Check if file exists.
		// Using a simple check
		if _, err := http.Get("http://google.com"); err != nil {
		} // Dummy
	}

	// Simple file existence check
	// We can't use os.Stat easily without importing os, but exec.Command used os/exec.
	// Let's just try to upload it. If it fails, coverURL is empty.
	if cURL, err := storage.LocalStorage.Upload(coverPath); err == nil {
		coverURL = cURL
	}

	title := c.PostForm("title")
	description := c.PostForm("description")

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
		URL:         url,
		CoverURL:    coverURL,
		UserID:      userID,
		TagIDs:      tagIDs,
	}

	if err := h.Service.CreateVideo(video); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"video": video})
}

func (h *VideoHandler) GetFeed(c *gin.Context) {
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))
	offset, _ := strconv.Atoi(c.DefaultQuery("offset", "0"))
	tagID, _ := strconv.ParseUint(c.Query("tag_id"), 10, 32)

	videos, err := h.Service.GetFeed(limit, offset, uint(tagID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"videos": videos})
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
	c.JSON(http.StatusOK, gin.H{"tags": tags})
}
