package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"

	_ "github.com/mattn/go-sqlite3"

	_ "embed"

	"github.com/gin-gonic/gin"

	"github.com/johncmanuel/cpsc449-project2/db/sqlite"
	"github.com/johncmanuel/cpsc449-project2/pkgs/canvas"
	"github.com/johncmanuel/cpsc449-project2/pkgs/utils"
)

var uploadedFilesDir = "./uploads"

//go:embed db/sqlite/schema.sql
var ddl string

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

// Fetch assignments from Canvas and insert into sqlite db
func HandleAssignments(c *canvas.CanvasClient, q *sqlite.Queries) {
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
				// TODO: add course name to the DB
				params := sqlite.UpsertAssignmentParams{
					ID:       int64(assignment.ID),
					CourseID: int64(courseID),
					Name:     assignment.Name,
					DueDate:  utils.ConvertToNullTime(assignment.DueAt),
				}
				if _, err := q.UpsertAssignment(context.Background(), params); err != nil {
					fmt.Printf("Error inserting assignment: %v\n", err)
				}
			}
		}
	}
}

func SetupRouter(cli *canvas.CanvasClient, q *sqlite.Queries) *gin.Engine {
	r := gin.Default()

	// Setup routes and other endpoints here

	// Test route
	r.GET("/test", func(c *gin.Context) {
		ExampleCanvasAssignmentFetcher(cli)
	})

	r.GET("/assignments", func(c *gin.Context) {
		// ExampleCanvasAssignmentFetcher(cli)
		HandleAssignments(cli, q)
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

	db, err := sql.Open("sqlite3", "./db/canvas.db")
	if err != nil {
		panic(fmt.Sprintf("Error opening database: %v", err))
	}
	defer db.Close()

	// create the tables if they don't exist
	if _, err := db.ExecContext(context.Background(), ddl); err != nil {
		panic(fmt.Sprintf("Error creating tables: %v", err))
	}

	q := sqlite.New(db)

	client := canvas.NewCanvasClient(BASE_CANVAS_URL, CANVAS_TOKEN)

	router := SetupRouter(client, q)

	// test redis
	// r := redis.GetInstance()
	// r.Set("test", "value")
	// val, err := r.Get("test")
	// if err != nil {
	// 	fmt.Println("REDIS: Error getting key:", err)
	// }
	// fmt.Println(val)
	// r.Delete("test")
	// if _, err := r.Get("test"); err != nil {
	// 	fmt.Println("REDIS: Key not found, which is expected", err)
	// }

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	fmt.Printf("Server starting on port %s...\n", port)
	router.Run(":" + port)
}
