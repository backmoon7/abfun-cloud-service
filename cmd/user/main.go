package main

import (
	"bilibili-clone/internal/user/api"
	"bilibili-clone/internal/user/model"
	"bilibili-clone/pkg/database"
	"bilibili-clone/pkg/middleware"
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
	err := database.DB.AutoMigrate(&model.User{})
	if err != nil {
		log.Fatalf("Failed to migrate database: %v", err)
	}

	// Initialize Router
	r := gin.Default()
	r.Use(middleware.CORS())
	userHandler := api.NewUserHandler()

	v1 := r.Group("/api/v1/user")
	{
		v1.POST("/register", userHandler.Register)
		v1.POST("/login", userHandler.Login)
	}

	// Run Server
	if err := r.Run(":8081"); err != nil {
		log.Fatalf("Failed to run server: %v", err)
	}
}
