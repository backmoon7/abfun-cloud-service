package worker

import (
        "bilibili-clone/internal/video/model"
        "bilibili-clone/pkg/cache"
        "bilibili-clone/pkg/database"
        "context"
        "encoding/json"
        "fmt"
        "log"
)

type LikeConsumer struct{}

func NewLikeConsumer() *LikeConsumer {
        return &LikeConsumer{}
}

type likeEvent struct {
        Type    string `json:"type"`
        UserID  uint   `json:"user_id"`
        VideoID uint   `json:"video_id"`
}

func (c *LikeConsumer) Handle(msg []byte) {
        var event likeEvent
        if err := json.Unmarshal(msg, &event); err != nil {
                log.Printf("invalid like event: %v", err)
                return
        }

        key := fmt.Sprintf("video:%d:likes", event.VideoID)
        ctx := context.Background()

        switch event.Type {
        case "like":
                if err := cache.RDB.Incr(ctx, key).Err(); err != nil {
                        log.Printf("failed to increment likes in cache: %v", err)
                }
                // Update database
                if err := database.DB.Model(&model.Video{}).Where("id = ?", event.VideoID).
                        UpdateColumn("like_count", database.DB.Raw("like_count + 1")).Error; err != nil {
                        log.Printf("failed to increment like_count in db: %v", err)
                }
        case "unlike":
                if err := cache.RDB.Decr(ctx, key).Err(); err != nil {
                        log.Printf("failed to decrement likes in cache: %v", err)
                }
                // Update database
                if err := database.DB.Model(&model.Video{}).Where("id = ? AND like_count > 0", event.VideoID).
                        UpdateColumn("like_count", database.DB.Raw("like_count - 1")).Error; err != nil {
                        log.Printf("failed to decrement like_count in db: %v", err)
                }
        default:
                log.Printf("unknown event type: %s", event.Type)
        }
}
