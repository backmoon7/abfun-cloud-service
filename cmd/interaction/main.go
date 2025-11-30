package main

import (
"bilibili-clone/internal/interaction/api"
"bilibili-clone/internal/interaction/model"
"bilibili-clone/internal/interaction/worker"
"bilibili-clone/pkg/cache"
"bilibili-clone/pkg/database"
"bilibili-clone/pkg/middleware"
"bilibili-clone/pkg/mq"
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

// Auto Migrate
database.DB.AutoMigrate(&model.VideoLike{}, &model.FavoriteFolder{}, &model.FavoriteItem{}, &model.Comment{})

// Initialize Redis
redisAddr := os.Getenv("REDIS_ADDR")
if redisAddr == "" {
redisAddr = "localhost:6379"
}
cache.InitRedis(redisAddr)

// Initialize RabbitMQ connection
rabbitURL := os.Getenv("RABBITMQ_URL")
if rabbitURL == "" {
rabbitURL = "amqp://guest:guest@localhost:5672/"
}
if err := mq.InitRabbitMQ(rabbitURL); err != nil {
log.Fatalf("Failed to connect RabbitMQ: %v", err)
}

// Initialize RabbitMQ consumer (background)
consumer := worker.NewLikeConsumer()
if err := mq.Consume("interaction_event", consumer.Handle); err != nil {
log.Fatalf("Failed to start interaction consumer: %v", err)
}

// Initialize Router
r := gin.Default()
r.Use(middleware.CORS())
handler := api.NewInteractionHandler()

v1 := r.Group("/api/v1/interaction")

// Public routes
v1.GET("/comments", handler.GetComments)

// Protected routes
protected := v1.Group("/")
protected.Use(middleware.AuthRequired())
{
protected.POST("/like", handler.Like)
protected.DELETE("/like", handler.Unlike)
protected.POST("/favorites/folders", handler.CreateFolder)
protected.GET("/favorites/folders", handler.GetFolders)
protected.POST("/favorites/add", handler.AddFavorite)

// Comments
protected.POST("/comments", handler.PostComment)
protected.DELETE("/comments/:id", handler.DeleteComment)
}

// Run Server
if err := r.Run(":8084"); err != nil {
log.Fatalf("Failed to run server: %v", err)
}
}
