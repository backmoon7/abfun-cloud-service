package main

import (
	"bilibili-clone/internal/danmaku/api"
	"bilibili-clone/internal/danmaku/model"
	"bilibili-clone/internal/danmaku/ws"
	"bilibili-clone/pkg/database"
	"bilibili-clone/pkg/middleware"
	"bilibili-clone/pkg/mq"
	"log"
	"os"

	"bilibili-clone/internal/danmaku/worker"
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
	database.DB.AutoMigrate(&model.Danmaku{})

	// Initialize RabbitMQ connection
	rabbitURL := os.Getenv("RABBITMQ_URL")
	if rabbitURL == "" {
		rabbitURL = "amqp://guest:guest@localhost:5672/"
	}
	if err := mq.InitRabbitMQ(rabbitURL); err != nil {
		log.Fatalf("Failed to connect RabbitMQ: %v", err)
	}

	// Initialize WebSocket Hub
	hub := ws.NewHub()
	go hub.Run()

	// Initialize RabbitMQ consumer (Background)
	consumer := worker.NewDanmakuConsumer(hub)
	if err := mq.Consume("danmaku_send", consumer.Handle); err != nil {
		log.Fatalf("Failed to start danmaku consumer: %v", err)
	}

	// Initialize Router
	r := gin.Default()
	r.Use(middleware.CORS())
	handler := api.NewDanmakuHandler(hub)

	v1 := r.Group("/api/v1/danmaku")
	{
		v1.GET("/list", handler.GetDanmakus)
                v1.POST("/send", middleware.AuthRequired(), handler.SendDanmaku)
		v1.GET("/ws", handler.ServeWS)
	}

	// Run Server
	if err := r.Run(":8083"); err != nil {
		log.Fatalf("Failed to run server: %v", err)
	}
}
