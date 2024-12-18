package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"

	_ "github.com/mattn/go-sqlite3"
	"github.com/openai/openai-go"
	"github.com/openai/openai-go/option"

	_ "embed"

	"github.com/gin-gonic/gin"

	"github.com/johncmanuel/cpsc449-project2/db/sqlite"
	"github.com/johncmanuel/cpsc449-project2/pkgs/canvas"
	"github.com/johncmanuel/cpsc449-project2/pkgs/redis"
	"github.com/johncmanuel/cpsc449-project2/pkgs/utils"
)

var uploadedFilesDir = "./uploads"

//go:embed db/sqlite/schema.sql
var ddl string

// for OpenAI request
type Assignment struct {
	ID         string `json:"id"`
	Name       string `json:"name"`
	CourseID   string `json:"course_id"`
	DueDate    string `json:"due_date"`
	Difficulty string `json:"difficulty,omitempty"`
	Length     string `json:"length,omitempty"`
}

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
					ID:         int64(assignment.ID),
					CourseID:   int64(courseID),
					Name:       assignment.Name,
					DueDate:    utils.ConvertToNullTime(assignment.DueAt),
					Difficulty: sql.NullInt64{},
					Length:     sql.NullInt64{},
				}
				if _, err := q.UpsertAssignment(context.Background(), params); err != nil {
					fmt.Printf("Error inserting assignment: %v\n", err)
				}
			}
		}
	}
}

// Gets an individual assignment from the DB
func GetAssignment(c *gin.Context, q *sqlite.Queries) {
	courseID := c.Param("courseID")
	assignmentID := c.Param("assignmentID")
	r := redis.GetInstance()
	keys := redis.GenerateTupleKey(courseID, assignmentID)

	// Check if it exists in cache first
	// if it does, return it
	// if not, get it from the DB and cache it
	// if it doesn't exist in the DB, return 404
	e, err := r.Exists(keys)
	if err != nil {
		fmt.Printf("Error checking for key: %v\n", err)
	}
	if !e {
		fmt.Println("Value not found in cache, retrieving from DB.")
		// get from the DB
		params := sqlite.GetAssignmentParams{
			CourseID: utils.ConvertStringToInt64(courseID),
			ID:       utils.ConvertStringToInt64(assignmentID),
		}
		assignment, err := q.GetAssignment(context.Background(), params)
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{
				"error": "Assignment not found",
			})
			return
		}
		// cache it
		if err := r.Set(keys, assignment); err != nil {
			fmt.Printf("Error caching assignment: %v\n", err)
		}
		c.JSON(http.StatusOK, assignment)
		return
	}
	// get from cache
	val, err := r.Get(keys)
	if err != nil {
		fmt.Printf("Error getting key: %v\n", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Internal server error",
		})
		return
	}
	fmt.Println("retrieving from cache lfg")
	// deserialize json back to struct
	var cachedAssignment sqlite.Assignment
	if err := json.Unmarshal([]byte(val), &cachedAssignment); err != nil {
		fmt.Printf("Error unmarshalling assignment: %v\n", err)
	}
	c.JSON(http.StatusOK, cachedAssignment)
}

// Simple function to get all assignments directly from DB
func GetAllAssignments(c *gin.Context, q *sqlite.Queries) {
	// Attempt to fetch all assignments from the DB
	assignments, err := q.ListAllAssignments(context.Background())
	if err != nil {
		// If an error occurs, print it and return a server error response
		fmt.Printf("Error fetching assignments from DB: %v\n", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": fmt.Sprintf("Error fetching assignments: %v", err),
		})
		return
	}

	// Return the assignments as JSON if successful
	c.JSON(http.StatusOK, assignments)
}

func DeleteAssignment(c *gin.Context, q *sqlite.Queries) {
	courseID := c.Param("courseID")
	assignmentID := c.Param("assignmentID")
	r := redis.GetInstance()
	keys := redis.GenerateTupleKey(courseID, assignmentID)
	params := sqlite.DeleteAssignmentParams{
		CourseID: utils.ConvertStringToInt64(courseID),
		ID:       utils.ConvertStringToInt64(assignmentID),
	}

	if err := q.DeleteAssignment(context.Background(), params); err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "Assignment not found",
		})
		return
	}

	// remove from cache if its there
	if err := r.Delete(keys); err != nil {
		fmt.Printf("Error deleting key from cache: %v\n", err)
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Assignment deleted",
	})
}

func SetupRouter(cli *canvas.CanvasClient, q *sqlite.Queries, ocli *openai.Client) *gin.Engine {
	r := gin.Default()

	// Test route
	r.GET("/test", func(c *gin.Context) {
		ExampleCanvasAssignmentFetcher(cli)
	})

	r.GET("/:courseID/assignments/:assignmentID", func(c *gin.Context) {
		GetAssignment(c, q)
	})
	// r.PUT("/:courseID/assignments/:assignmentID", UpdateAssignment)
	r.DELETE("/:courseID/assignments/:assignmentID", func(c *gin.Context) {
		DeleteAssignment(c, q)
	})
	r.GET("/assignments", func(c *gin.Context) {
		// ExampleCanvasAssignmentFetcher(cli)
		HandleAssignments(cli, q)
	})

	// Route to get all assignments directly from the DB (no caching)
	r.GET("/all-assignments", func(c *gin.Context) {
		GetAllAssignments(c, q)
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
	// Load environment variables
	utils.LoadEnv()

	var (
		BASE_CANVAS_URL = utils.GetEnv("CANVAS_URL")
		CANVAS_TOKEN    = utils.GetEnv("CANVAS_TOKEN")
	)

	// Open SQLite database connection
	db, err := sql.Open("sqlite3", "./db/canvas.db")
	if err != nil {
		panic(fmt.Sprintf("Error opening database: %v", err))
	}
	defer db.Close()

	// Create the tables if they don't exist
	if _, err := db.ExecContext(context.Background(), ddl); err != nil {
		panic(fmt.Sprintf("Error creating tables: %v", err))
	}

	q := sqlite.New(db)

	// Initialize Canvas client
	client := canvas.NewCanvasClient(BASE_CANVAS_URL, CANVAS_TOKEN)

	// Retrieve the OpenAI API key from environment variables
	apiKey := utils.GetEnv("OPENAI_API_KEY")
	if apiKey == "" {
		fmt.Println("OPENAI_API_KEY is not set in the environment variables")
		return
	}

	// Initialize the OpenAI client
	openAIclient := openai.NewClient(
		option.WithAPIKey("My API Key"),
	)

	// Set up the router with dependencies
	router := SetupRouter(client, q, openAIclient)

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

	// Determine server port
	port := utils.GetEnv("PORT")
	if port == "" {
		port = "8080"
	}

	fmt.Printf("Server starting on port %s...\n", port)
	if err := router.Run(":" + port); err != nil {
		panic(fmt.Sprintf("Failed to start server: %v", err))
	}
}
