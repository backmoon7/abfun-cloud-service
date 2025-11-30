package main

import (
"bilibili-clone/internal/video/api"
"bilibili-clone/internal/video/model"
"bilibili-clone/internal/video/service"
"bilibili-clone/internal/video/worker"
"bilibili-clone/pkg/database"
"bilibili-clone/pkg/middleware"
"bilibili-clone/pkg/mq"
"bilibili-clone/pkg/storage"
"log"
"os"

"github.com/gin-gonic/gin"
)

func main() {
// Initialize Database
dsn := os.Getenv("MYSQL_DSN")
if dsn == "" {
dsn = "root:root@tcp(localhost:3306)/bilibili?charset=utf8mb4&parseTime=True&loc=Local"
}
database.InitMySQL(dsn)

// Initialize Local Storage
storagePath := os.Getenv("STORAGE_PATH")
if storagePath == "" {
storagePath = "./uploads"
}
domainURL := os.Getenv("DOMAIN_URL")
if domainURL == "" {
domainURL = "http://localhost:80/uploads"
}
storage.InitLocalStorage(storagePath, domainURL)

// Auto Migrate
err := database.DB.AutoMigrate(&model.Video{}, &model.Tag{}, &model.VideoTag{})
if err != nil {
log.Fatalf("Failed to migrate database: %v", err)
}

// Initialize RabbitMQ connection
rabbitURL := os.Getenv("RABBITMQ_URL")
if rabbitURL == "" {
rabbitURL = "amqp://guest:guest@localhost:5672/"
}
if err := mq.InitRabbitMQ(rabbitURL); err != nil {
log.Fatalf("Failed to connect RabbitMQ: %v", err)
}

// Initialize RabbitMQ consumer (background)
consumer := worker.NewTranscodeConsumer()
if err := mq.Consume("video_uploaded", consumer.Handle); err != nil {
log.Fatalf("Failed to start transcode consumer: %v", err)
}

// Init Tags
svc := service.NewVideoService()
svc.InitTags()

// Initialize Router
r := gin.Default()
r.Use(middleware.CORS())
videoHandler := api.NewVideoHandler()

v1 := r.Group("/api/v1")
{
v1.GET("/tags", videoHandler.ListTags)
v1.GET("/videos/feed", videoHandler.Feed)
        v1.GET("/search", videoHandler.Search)

// Protected routes
protected := v1.Group("")
protected.Use(middleware.AuthRequired())
{
protected.POST("/videos", videoHandler.Upload)
}
}

// Run Server
if err := r.Run(":8082"); err != nil {
log.Fatalf("Failed to run server: %v", err)
}
}
