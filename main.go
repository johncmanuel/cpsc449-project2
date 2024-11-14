package main

import (
	"fmt"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func loadEnv() {
	err := godotenv.Load()
	if err != nil {
		panic("Error loading .env file")
	}
}

func setupRouter() *gin.Engine {
	r := gin.Default()

	// Setup routes and other endpoints here
	r.GET("/ping", func(c *gin.Context) {
		c.String(http.StatusOK, "pong")
	})

	// Setup middleware
	r.Use(gin.Recovery())

	return r
}

func main() {
	loadEnv()

	router := setupRouter()

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	fmt.Printf("Server starting on port %s...\n", port)
	router.Run(":" + port)
}
