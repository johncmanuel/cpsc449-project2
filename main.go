package main

import (
	"fmt"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"

	"github.com/johncmanuel/cpsc449-project2/pkgs/canvas"
	"github.com/johncmanuel/cpsc449-project2/pkgs/utils"
)

// Just for printing and testing the API
func ExampleCanvasAssignmentFetcher(c *canvas.CanvasClient) {
	allAssignments, err := c.GetAllAssignmentsForCurrentTerm()
	if err != nil {
		fmt.Printf("Error fetching assignments: %v\n", err)
		return
	}

	for courseID, courseAssignments := range allAssignments {
		for courseName, assignments := range courseAssignments {
			fmt.Printf("Course: %s (ID: %d)\n", courseName, courseID)
			for _, assignment := range assignments {
				fmt.Printf("- Assignment: %s (ID: %d), Due: %s\n",
					assignment.Name, assignment.ID, assignment.DueAt)
			}
		}
	}
}

func SetupRouter(cli *canvas.CanvasClient) *gin.Engine {
	r := gin.Default()

	// Setup routes and other endpoints here
	r.GET("/ping", func(c *gin.Context) {
		c.String(http.StatusOK, "pong")
	})

	r.GET("/assignments", func(c *gin.Context) {
		ExampleCanvasAssignmentFetcher(cli)
	})

	return r
}

func main() {
	utils.LoadEnv()
	var (
		BASE_CANVAS_URL = utils.GetEnv("CANVAS_URL")
		CANVAS_TOKEN    = utils.GetEnv("CANVAS_TOKEN")
	)

	client := canvas.NewCanvasClient(BASE_CANVAS_URL, CANVAS_TOKEN)

	router := SetupRouter(client)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	fmt.Printf("Server starting on port %s...\n", port)
	router.Run(":" + port)
}
