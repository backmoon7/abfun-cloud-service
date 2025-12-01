package worker

import (
	"bilibili-clone/internal/video/service"
	"bilibili-clone/pkg/storage"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

type TranscodeConsumer struct {
	svc *service.VideoService
}

func NewTranscodeConsumer() *TranscodeConsumer {
	return &TranscodeConsumer{
		svc: service.NewVideoService(),
	}
}

func (c *TranscodeConsumer) Handle(body []byte) error {
	var event struct {
		VideoID  uint   `json:"video_id"`
		VideoURL string `json:"video_url"`
	}
	if err := json.Unmarshal(body, &event); err != nil {
		return err
	}

	log.Printf("Processing video %d: %s", event.VideoID, event.VideoURL)

	// Construct paths
	filename := filepath.Base(event.VideoURL)
	inputPath := filepath.Join(storage.LocalStoragePath, filename)

	// Output directory: /app/uploads/videos/{video_id}/
	outputDir := filepath.Join(storage.LocalStoragePath, "videos", fmt.Sprintf("%d", event.VideoID))
	if err := os.MkdirAll(outputDir, 0755); err != nil {
		return err
	}

	// 1. Generate Cover
	coverFilename := "cover.jpg"
	coverPath := filepath.Join(outputDir, coverFilename)
	cmdCover := exec.Command("ffmpeg", "-i", inputPath, "-ss", "00:00:01", "-vframes", "1", coverPath)
	if err := cmdCover.Run(); err != nil {
		log.Printf("Failed to generate cover: %v", err)
	}

	// 2. Transcode to HLS
	hlsFilename := "index.m3u8"
	hlsPath := filepath.Join(outputDir, hlsFilename)
	// ffmpeg -i input.mp4 -c:v libx264 -c:a aac -hls_time 10 -hls_list_size 0 -f hls index.m3u8
	cmdHLS := exec.Command("ffmpeg", "-i", inputPath, "-c:v", "libx264", "-c:a", "aac", "-hls_time", "10", "-hls_list_size", "0", "-f", "hls", hlsPath)

	if output, err := cmdHLS.CombinedOutput(); err != nil {
		log.Printf("Transcoding failed: %v, output: %s", err, string(output))
		c.svc.UpdateVideoStatus(event.VideoID, -1, "", "")
		return err
	}

	// 3. Update DB
	domainURL := os.Getenv("DOMAIN_URL")
	if domainURL == "" {
		domainURL = "https://api.abfun.me/uploads"
	}
	domainURL = strings.TrimSuffix(domainURL, "/")

	newVideoURL := fmt.Sprintf("%s/videos/%d/%s", domainURL, event.VideoID, hlsFilename)
	newCoverURL := fmt.Sprintf("%s/videos/%d/%s", domainURL, event.VideoID, coverFilename)

	if err := c.svc.UpdateVideoStatus(event.VideoID, 1, newVideoURL, newCoverURL); err != nil {
		return err
	}

	// Sync to Meilisearch
	if err := c.svc.SyncVideoToMeilisearch(event.VideoID); err != nil {
		log.Printf("Failed to sync video %d to Meilisearch: %v", event.VideoID, err)
	}

	// 4. Delete original file
	if err := os.Remove(inputPath); err != nil {
		log.Printf("Failed to remove original file: %v", err)
	}

	log.Printf("Video %d processed successfully", event.VideoID)
	return nil
}
