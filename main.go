package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"

	"github.com/gin-gonic/gin"

	"github.com/johncmanuel/cpsc449-project2/pkgs/canvas"
	"github.com/johncmanuel/cpsc449-project2/pkgs/utils"
)

var uploadedFilesDir = "./uploads"

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

	// It's possible to get the syllabus through Canvas API; however, some
	// teachers only upload their syllabus file in the the syllabus page. The API returns
	// the HTML content of the syllabus page, so not sure if the API would be able to return the file
	r.POST("/syllabus", func(c *gin.Context) {
		f, err := c.FormFile("file")
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "No file uploaded",
			})
			return
		}

		fpath := filepath.Join(uploadedFilesDir, f.Filename)

		if err := os.MkdirAll(uploadedFilesDir, os.ModePerm); err != nil {
			log.Println("Error creating directory:", err)
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "Internal server error",
			})
			return
		}

		if err := c.SaveUploadedFile(f, fpath); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "Failed to save file",
			})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"filename": f.Filename,
			"message":  "File uploaded successfully",
		})
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
